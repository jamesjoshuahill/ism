/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms
of the under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

const (
	serviceName = "overview-service"
	planName    = "simple"
	brokerPort  = 8081
)

var (
	controllerSession *Session
	proxySession      *Session

	nodePathToCLI      string
	nodeBrokerURL      string
	nodeBrokerUsername string
	nodeBrokerPassword string
)

func TestAcceptance(t *testing.T) {
	SetDefaultEventuallyTimeout(time.Second * 5)
	SetDefaultConsistentlyDuration(time.Second * 5)

	SynchronizedBeforeSuite(func() []byte {
		printTestSetup()
		cliPath := buildCLI()
		installCRDs()

		var brokerURL, brokerUser, brokerPass string

		if testingInCluster() {
			brokerURL, brokerUser, brokerPass = deployTestBroker()
			deployController()
		} else {
			brokerURL, brokerUser, brokerPass = startTestBroker()
			startController()
		}

		data := nodeData{
			PathToCLI:      cliPath,
			BrokerURL:      brokerURL,
			BrokerUsername: brokerUser,
			BrokerPassword: brokerPass,
		}

		b, err := json.Marshal(data)
		Expect(err).NotTo(HaveOccurred())

		return []byte(b)
	}, func(rawNodeData []byte) {
		var data nodeData
		err := json.Unmarshal(rawNodeData, &data)
		Expect(err).NotTo(HaveOccurred())

		nodePathToCLI = data.PathToCLI
		nodeBrokerURL = data.BrokerURL
		nodeBrokerUsername = data.BrokerUsername
		nodeBrokerPassword = data.BrokerPassword
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		if testingInCluster() {
			deleteController()
			deleteTestBroker()
		} else {
			stopController()
			stopTestBroker()
		}

		CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func printTestSetup() {
	clusterContext := strings.TrimSpace(runKubectl("config", "current-context"))
	fmt.Printf("Running tests against the '%s' cluster\n", clusterContext)
	if testingInCluster() {
		fmt.Println("Deploying test components in kubernetes")
	} else {
		fmt.Println("Deploying test components locally")
	}
}

func testingInCluster() bool {
	return os.Getenv("TEST_ACCEPTANCE_IN_CLUSTER") != ""
}

func buildCLI() string {
	cliPath, err := Build("github.com/pivotal-cf/ism/cmd/ism")
	Expect(err).NotTo(HaveOccurred())

	return cliPath
}

func installCRDs() {
	runMake("install")
}

func startTestBroker() (string, string, string) {
	runMake("run-test-broker")
	return fmt.Sprintf("http://127.0.0.1:%d", brokerPort), "admin", "password"
}

func stopTestBroker() {
	runMake("stop-test-broker")
}

func deployTestBroker() (string, string, string) {
	runMake("deploy-test-broker")
	runKubectl("wait", "--for=condition=available", "deployment/overview-broker-deployment")
	brokerIP := runKubectl("get", "service", "overview-broker", "-o", "jsonpath={.spec.clusterIP}")

	setupProxyAccessToBroker()
	return fmt.Sprintf("http://%s:8080", brokerIP), "admin", "password"
}

func setupProxyAccessToBroker() {
	cmd := exec.Command("kubectl", "port-forward", "service/overview-broker", fmt.Sprintf("%d:8080", brokerPort))

	var err error
	proxySession, err = Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
}

func teardownProxyAccessToBroker() {
	proxySession.Terminate()
}

func deleteTestBroker() {
	runMake("delete-test-broker")
	teardownProxyAccessToBroker()
}

func startController() {
	pathToController, err := Build("github.com/pivotal-cf/ism/cmd/manager")
	Expect(err).NotTo(HaveOccurred())

	command := exec.Command(pathToController)
	controllerSession, err = Start(command, GinkgoWriter, GinkgoWriter)
	Eventually(controllerSession.Err).Should(Say("Starting the Cmd"))

	Expect(err).NotTo(HaveOccurred())
}

func stopController() {
	controllerSession.Terminate()
}

func deployController() {
	runMake("deploy")
	runKubectl("wait", "-n", "ism-system", "--for=condition=Ready", "pod/ism-controller-manager-0")
}

func deleteController() {
	runMake("delete-controller")
}

func uninstallCRDs() {
	runMake("uninstall-crds")
}

func registerBroker(brokerName string) {
	registerArgs := []string{"broker", "register",
		"--name", brokerName,
		"--url", nodeBrokerURL,
		"--username", nodeBrokerUsername,
		"--password", nodeBrokerPassword}
	command := exec.Command(nodePathToCLI, registerArgs...)
	registerSession, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(registerSession).Should(Exit(0))
}

func createInstance(instanceName, brokerName string) {
	createArgs := []string{"instance", "create",
		"--name", instanceName,
		"--service", serviceName,
		"--plan", planName,
		"--broker", brokerName}
	command := exec.Command(nodePathToCLI, createArgs...)
	createSession, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(createSession).Should(Exit(0))
}

func deleteBrokers(brokerNames ...string) {
	for _, b := range brokerNames {
		runKubectl("delete", "broker", b)
	}
}

func deleteInstances(serviceInstanceNames ...string) {
	for _, s := range serviceInstanceNames {
		runKubectl("delete", "serviceinstance", s)
	}
}

func deleteBindings(serviceBindingNames ...string) {
	for _, b := range serviceBindingNames {
		runKubectl("delete", "servicebinding", b)
	}
}

func runMake(task string) {
	command := exec.Command("make", task)
	command.Dir = filepath.Join("..")
	command.Stdout = GinkgoWriter
	command.Stderr = GinkgoWriter
	Expect(command.Run()).To(Succeed())
}

func runKubectl(args ...string) string {
	outBuf := NewBuffer()
	command := exec.Command("kubectl", args...)
	command.Dir = filepath.Join("..")
	command.Stdout = io.MultiWriter(GinkgoWriter, outBuf)
	command.Stderr = GinkgoWriter
	Expect(command.Run()).To(Succeed())

	return string(outBuf.Contents())
}

type nodeData struct {
	PathToCLI      string
	BrokerURL      string
	BrokerUsername string
	BrokerPassword string
}

type brokerData struct {
	ServiceInstances map[string]serviceInstance `json:"serviceInstances"`
}

type serviceInstance struct {
	PlanName    string                    `json:"plan_name"`
	ServiceName string                    `json:"service_name"`
	Bindings    map[string]serviceBinding `json:"bindings"`
}

type serviceBinding struct {
}

func getBrokerData() brokerData {
	brokerDataURL := fmt.Sprintf("http://127.0.0.1:%d/data", brokerPort)

	resp, err := http.Get(brokerDataURL)
	Expect(err).NotTo(HaveOccurred())
	respBytes, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	var data brokerData
	Expect(json.Unmarshal(respBytes, &data)).To(Succeed())

	return data
}

func getBrokerInstances() []serviceInstance {
	var instances []serviceInstance
	for _, instance := range getBrokerData().ServiceInstances {
		instances = append(instances, instance)
	}

	return instances
}

func getBrokerBindings() []serviceBinding {
	var bindings []serviceBinding
	instances := getBrokerData().ServiceInstances
	for _, instance := range instances {
		for _, binding := range instance.Bindings {
			bindings = append(bindings, binding)
		}
	}
	return bindings
}
