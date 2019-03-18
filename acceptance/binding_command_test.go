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
	"time"

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
			Eventually(session).Should(Say(`ism \[OPTIONS\] binding <create>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The binding command group lets you create service bindings"))
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

		PWhen("valid args are passed", func() {
			BeforeEach(func() {
				registerBroker("binding-creation-broker")
				createInstance("binding-creation-instance", "binding-creation-broker")
				args = append(args, "--name", "binding-creation-binding", "--instance-name", "binding-creation-instance")

				Expect(getBrokerBindings()).To(HaveLen(0))

				//TODO remove this
				time.Sleep(time.Millisecond * 500)
			})

			AfterEach(func() {
				deleteBrokers("binding-creation-broker")
				deleteInstances("binding-creation-instance")
				// deleteServiceBindings("binding-creation-binding")
			})

			It("starts creation of the service binding", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Binding 'my-binding' is being created\\."))

				// TODO: Revisit this when it comes to implementing asynchronous provisioning
				// Allow time for controller to set instance status to "created"
				Expect(getBrokerBindings()).To(HaveLen(1))
			})
		})

		// 	// TODO: When writing "ism binding list" test make sure to check the /data
		// 	// endpoint on the broker to ensure the binding has _actually_ been created.
		// })

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session).Should(Say("the required flags `--instance-name' and `--name' were not specified"))
			})
		})
	})
	//
	// 	Describe("list sub command", func() {
	// 		BeforeEach(func() {
	// 			args = append(args, "list")
	// 		})
	//
	// 		When("--help is passed", func() {
	// 			BeforeEach(func() {
	// 				args = append(args, "--help")
	// 			})
	//
	// 			It("displays help and exits 0", func() {
	// 				Eventually(session).Should(Exit(0))
	// 				Eventually(session).Should(Say("Usage:"))
	// 				Eventually(session).Should(Say(`ism \[OPTIONS\] binding list`))
	// 				Eventually(session).Should(Say("\n"))
	// 				Eventually(session).Should(Say("List the service bindings"))
	// 			})
	// 		})
	//
	// 		When("0 service bindings are created", func() {
	// 			It("displays 'No bindings found.' and exits 0", func() {
	// 				Eventually(session).Should(Exit(0))
	// 				Eventually(session).Should(Say("No bindings found\\."))
	// 			})
	// 		})
	//
	// 		When("1 binding is created", func() {
	// 			BeforeEach(func() {
	// 				registerBroker("binding-list-command-broker")
	// 				createbinding("binding-list-test-binding", "binding-list-command-broker")
	//
	// 				// TODO: Revisit this when it comes to implementing asynchronous provisioning
	// 				// Allow time for controller to set binding status to "created"
	// 				time.Sleep(time.Millisecond * 500)
	// 			})
	//
	// 			AfterEach(func() {
	// 				deleteServicebindings("binding-list-test-binding")
	// 				deleteBrokers("binding-list-command-broker")
	// 			})
	//
	// 			It("displays the binding", func() {
	// 				timeRegex := `\d{4,}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`
	//
	// 				Eventually(session).Should(Exit(0))
	// 				Eventually(session).Should(Say("NAME\\s+SERVICE\\s+PLAN\\s+BROKER\\s+STATUS\\s+CREATED AT"))
	// 				Eventually(session).Should(Say("binding-list-test-binding\\s+" + serviceName + "\\s+" + planName + "\\s+binding-list-command-broker\\s+" + "created" + "\\s+" + timeRegex))
	// 			})
	// 		})
	// 	})
	// })
	//
	// type servicebinding struct {
	// 	PlanName    string `json:"plan_name"`
	// 	ServiceName string `json:"service_name"`
	// }
	//
	// type brokerData struct {
	// 	Servicebindings map[string]servicebinding `json:"servicebindings"`
	// }
	//
	// func getBrokerData() brokerData {
	// 	brokerDataURL := fmt.Sprintf("http://127.0.0.1:%d/data", brokerPort)
	//
	// 	resp, err := http.Get(brokerDataURL)
	// 	Expect(err).NotTo(HaveOccurred())
	// 	respBytes, err := ioutil.ReadAll(resp.Body)
	// 	Expect(err).NotTo(HaveOccurred())
	//
	// 	var data brokerData
	// 	Expect(json.Unmarshal(respBytes, &data)).To(Succeed())
	//
	// 	return data

})
