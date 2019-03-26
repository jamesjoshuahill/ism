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

var _ = Describe("CLI binding command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"binding"}
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
			Eventually(session).Should(Say(`ism \[OPTIONS\] binding <create | get | list>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The binding command group lets you create, get and list service bindings"))
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
				Eventually(session).Should(Say(`ism \[OPTIONS\] binding create`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Create a service binding"))
			})
		})

		When("valid args are passed", func() {
			BeforeEach(func() {
				registerBroker("binding-creation-broker")
				createInstance("binding-creation-instance", "binding-creation-broker")
				args = append(args, "--name", "binding-creation-binding", "--instance", "binding-creation-instance")

				Expect(getBrokerBindings()).To(HaveLen(0))
			})

			AfterEach(func() {
				deleteBroker("binding-creation-broker")
				deleteInstance("binding-creation-instance")
				deleteBinding("binding-creation-binding")
				cleanBrokerData()
			})

			It("starts creation of the service binding", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Binding 'binding-creation-binding' is being created\\."))

				Eventually(getBrokerBindings).Should(HaveLen(1))
			})
		})

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("the required flags `--instance' and `--name' were not specified"))
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
				Eventually(session).Should(Say(`ism \[OPTIONS\] binding list`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("List the service bindings"))
			})
		})

		When("0 service bindings are created", func() {
			It("displays 'No bindings found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No bindings found\\."))
			})
		})

		When("1 binding is created", func() {
			BeforeEach(func() {
				registerBroker("binding-list-broker")
				createInstance("binding-list-instance", "binding-list-broker")
				createBinding("binding-list-binding", "binding-list-instance")
			})

			AfterEach(func() {
				deleteBinding("binding-list-binding")
				deleteInstance("binding-list-instance")
				deleteBroker("binding-list-broker")
				cleanBrokerData()
			})

			It("displays the binding", func() {
				timeRegex := `\d{4,}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("NAME\\s+INSTANCE\\s+STATUS\\s+CREATED AT"))
				Eventually(session).Should(Say("binding-list-binding\\s+binding-list-instance\\s+created\\s+" + timeRegex))
			})
		})
	})

	Describe("get sub command", func() {
		BeforeEach(func() {
			args = append(args, "get")
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`ism \[OPTIONS\] binding get`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Get a service binding"))
			})
		})

		When("the binding exists", func() {
			BeforeEach(func() {
				args = append(args, "--name", "binding-get-binding")
				registerBroker("binding-get-broker")
				createInstance("binding-get-instance", "binding-get-broker")
				createBinding("binding-get-binding", "binding-get-instance")
			})

			AfterEach(func() {
				deleteBinding("binding-get-binding")
				deleteInstance("binding-get-instance")
				deleteBroker("binding-get-broker")
				cleanBrokerData()
			})

			It("displays the binding and exits 0", func() {
				timeRegex := `\d{4,}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.+`

				bindings := getBrokerBindings()
				Expect(bindings).To(HaveLen(1))
				creds := bindings[0].Data["credentials"]
				username := creds["username"].(string)
				password := creds["password"].(string)

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("broker: binding-get-broker\n" +
					"createdAt:\\s+" + timeRegex + "\n" +
					"credentials:\n" +
					"\\s+password: " + password + "\n" +
					"\\s+username: " + username + "\n" +
					"instance: binding-get-instance\n" +
					"name: binding-get-binding\n" +
					"plan: simple\n" +
					"service: overview-service\n" +
					"status: created"))
			})
		})

		When("the binding does not exist", func() {
			BeforeEach(func() {
				args = append(args, "--name", "binding-get-non-existant-binding")
			})

			It("displays 'Binding not found' and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("binding not found"))
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
