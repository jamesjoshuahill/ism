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
			Eventually(session).Should(Say(`ism \[OPTIONS\] instance <create | list>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The instance command group lets you create and list service instances"))
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
			})

			AfterEach(func() {
				deleteBroker("instance-creation-broker")
				deleteInstance("my-instance")
				cleanBrokerData()
			})

			It("starts creation of the service instance", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Instance 'my-instance' is being created\\."))

				Eventually(getBrokerInstances).Should(HaveLen(1))
			})
		})

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("the required flags `--broker', `--name', `--plan' and `--service' were not specified"))
			})
		})
	})

	Describe("list sub command", func() {
		BeforeEach(func() {
			args = append(args, "list")
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`ism \[OPTIONS\] instance list`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("List the service instances"))
			})
		})

		When("0 service instances are created", func() {
			It("displays 'No instances found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No instances found\\."))
			})
		})

		When("1 instance is created", func() {
			BeforeEach(func() {
				registerBroker("instance-list-command-broker")
				createInstance("instance-list-test-instance", "instance-list-command-broker")
			})

			AfterEach(func() {
				deleteInstance("instance-list-test-instance")
				deleteBroker("instance-list-command-broker")
				cleanBrokerData()
			})

			It("displays the instance", func() {
				timeRegex := `\d{4,}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("NAME\\s+SERVICE\\s+PLAN\\s+BROKER\\s+STATUS\\s+CREATED AT"))
				Eventually(session).Should(Say("instance-list-test-instance\\s+" + serviceName + "\\s+" + planName + "\\s+instance-list-command-broker\\s+" + "created" + "\\s+" + timeRegex))
			})
		})
	})
})
