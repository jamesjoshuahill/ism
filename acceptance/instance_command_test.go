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
	"io/ioutil"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI instance command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"instance"}
	})

	JustBeforeEach(func() {
		var err error

		command := exec.Command(nodePathToCLI, args...)
		session, err = Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	When("--help is passed", func() {
		BeforeEach(func() {
			args = append(args, "--help")
		})

		It("displays help and exits 0", func() {
			Eventually(session).Should(Exit(0))
			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`ism \[OPTIONS\] instance <create>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The instance command group lets you create service instances"))
		})
	})

	Describe("create sub command", func() {
		BeforeEach(func() {
			args = append(args, "create")
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`ism \[OPTIONS\] instance create`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Create a service instance"))
			})
		})

		When("valid args are passed", func() {
			BeforeEach(func() {
				registerBroker("instance-creation-broker")
				args = append(args, "--name", "my-instance", "--service", serviceName, "--plan", planName, "--broker", "instance-creation-broker")

				Expect(getBrokerData().ServiceInstances).To(HaveLen(0))
			})

			AfterEach(func() {
				deleteBrokers("instance-creation-broker")
				deleteServiceInstances("my-instance")
			})

			It("starts creation of the service instance", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Instance 'my-instance' is being created\\."))
			})

			// TODO: When writing "ism instance list" test make sure to check the /data
			// endpoint on the broker to ensure the instance has _actually_ been created.
		})

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session).Should(Say("the required flags `--broker', `--name', `--plan' and `--service' were not specified"))
			})
		})
	})
})

type serviceInstance struct {
	PlanName    string `json:"plan_name"`
	ServiceName string `json:"service_name"`
}

type brokerData struct {
	ServiceInstances map[string]serviceInstance `json:"serviceInstances"`
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
