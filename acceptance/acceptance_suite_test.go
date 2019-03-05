package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
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

var (
	controllerSession *Session

	nodePathToCLI      string
	nodeBrokerURL      string
	nodeBrokerUsername string
	nodeBrokerPassword string
)

func TestAcceptance(t *testing.T) {
	SetDefaultEventuallyTimeout(time.Second * 5)
	SetDefaultConsistentlyDuration(time.Second * 5)

	SynchronizedBeforeSuite(func() []byte {
		printClusterName()
		cliPath := buildCLI()
		installCRDs()

		var brokerURL, brokerUser, brokerPass string
		if os.Getenv("TEST_ACCEPTANCE_IN_CLUSTER") != "" {
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
		if os.Getenv("TEST_ACCEPTANCE_IN_CLUSTER") != "" {
			deleteController()
			deleteTestBroker()
		} else {
			stopController()
			stopTestBroker()
		}

		uninstallCRDs()
		CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func printClusterName() {
	clusterContext := strings.TrimSpace(runKubectl("config", "current-context"))
	fmt.Printf("Running tests against the '%s' cluster\n", clusterContext)
}

func buildCLI() string {
	cliPath, err := Build("github.com/pivotal-cf/ism/cmd/ism")
	Expect(err).NotTo(HaveOccurred())

	return cliPath
}

func installCRDs() {
	runMake("install-crds")
}

func startTestBroker() (string, string, string) {
	runMake("run-test-broker")
	return "http://127.0.0.1:1122", "admin", "password"
}

func stopTestBroker() {
	runMake("terminate-test-broker")
}

func deployTestBroker() (string, string, string) {
	runMake("deploy-test-broker")
	runKubectl("wait", "--for=condition=available", "deployment/overview-broker-deployment")
	brokerIP := runKubectl("get", "service", "overview-broker", "-o", "jsonpath={.spec.clusterIP}")
	return fmt.Sprintf("http://%s:8080", brokerIP), "admin", "password"
}

func deleteTestBroker() {
	runMake("delete-test-broker")
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
	runMake("deploy-controller")
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

	//TODO: Temporarily sleep until #164240938 is done.
	time.Sleep(10 * time.Second)
}

func deleteBrokers(brokerNames ...string) {
	for _, b := range brokerNames {
		runKubectl("delete", "broker", b)
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
