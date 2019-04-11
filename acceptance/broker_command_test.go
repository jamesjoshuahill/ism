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

var _ = Describe("CLI broker command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"broker"}
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
			Eventually(session).Should(Say(`ism \[OPTIONS\] broker <list | register>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The broker command group lets you register, list, and delete service brokers"))
			Eventually(session).Should(Say("from the marketplace"))
		})
	})

	Describe("register sub command", func() {
		BeforeEach(func() {
			args = append(args, "register")
		})

		When("valid args are passed", func() {
			BeforeEach(func() {
				args = append(args, "--name", "my-broker", "--url", nodeBrokerURL, "--username", nodeBrokerUsername, "--password", nodeBrokerPassword)
			})

			AfterEach(func() {
				deleteBroker("my-broker")
			})

			It("displays a message that the registration has been successful", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Broker 'my-broker' registered\\."))
			})
		})

		When("a broker with the same name already exists", func() {
			BeforeEach(func() {
				registerBroker("register-dup-name-broker")
				args = append(args, "--name", "register-dup-name-broker", "--url", nodeBrokerURL, "--username", nodeBrokerUsername, "--password", nodeBrokerPassword)
			})

			AfterEach(func() {
				deleteBroker("register-dup-name-broker")
			})

			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("ERROR: A service broker named 'register-dup-name-broker' already exists."))
			})
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`ism \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a service broker into the marketplace"))
			})
		})

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("the required flags `--name', `--password', `--url' and `--username' were not specified"))
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
				Eventually(session).Should(Say(`ism \[OPTIONS\] broker list`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("List the service brokers in the marketplace"))
			})
		})

		When("0 brokers are registered", func() {
			It("displays 'No brokers found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No brokers found\\."))
			})
		})

		When("1 broker is registered", func() {
			BeforeEach(func() {
				registerBroker("broker-list-command-broker")
			})

			AfterEach(func() {
				deleteBroker("broker-list-command-broker")
			})

			It("displays the broker", func() {
				timeRegex := `\d{4,}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("NAME\\s+URL\\s+CREATED AT"))
				Eventually(session).Should(Say("broker-list-command-broker\\s+" + nodeBrokerURL + "\\s+" + timeRegex))
			})
		})
	})

	Describe("delete sub command", func() {
		BeforeEach(func() {
			args = append(args, "delete")
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`ism \[OPTIONS\] broker delete`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Delete a service broker from the marketplace"))
			})
		})

		When("valid args are passed", func() {
			BeforeEach(func() {
				args = append(args, "--name", "instance-deletion-broker")

				registerBroker("instance-deletion-broker")

				Expect(getBrokers()).To(HaveLen(1))
			})

			It("deletes the service broker", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Broker 'instance-deletion-broker' is being deleted\\."))

				Eventually(getBrokers).Should(HaveLen(0))
			})
		})

		When("valid args are passed and the broker has a service instance", func() {
			BeforeEach(func() {
				args = append(args, "--name", "instance-deletion-broker")

				registerBroker("instance-deletion-broker")
				createInstance("instance-deletion-instance", "instance-deletion-broker")

				Expect(getBrokers()).To(HaveLen(1))
				Expect(getBrokerInstances()).To(HaveLen(1))
			})

			AfterEach(func() {
				deleteInstance("instance-deletion-instance")
				deleteBroker("instance-deletion-broker")
				cleanBrokerData()
			})

			It("doesn't delete the instance", func() {
				Expect(getBrokerInstances()).To(HaveLen(1))
			})

			It("errors with a useful message", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Broker 'instance-deletion-broker' cannot be deleted as it has 1 instance"))
			})
		})

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("the required flag `--name' was not specified"))
			})
		})
	})
})
