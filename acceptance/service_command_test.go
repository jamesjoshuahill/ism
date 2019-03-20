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

var _ = Describe("CLI service command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"service"}
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
			Eventually(session).Should(Say(`ism \[OPTIONS\] service <list>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The service command group lets you list the available services in the"))
			Eventually(session).Should(Say("marketplace"))
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
				Eventually(session).Should(Say(`ism \[OPTIONS\] service list`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("List the services that are available in the marketplace"))
			})
		})

		When("0 brokers are registered", func() {
			It("displays 'No services found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No services found\\."))
			})
		})

		When("1 broker is registered", func() {
			BeforeEach(func() {
				registerBroker("service-list-command-broker")
			})

			AfterEach(func() {
				deleteBroker("service-list-command-broker")
			})

			It("displays services and plans for the broker", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("SERVICE\\s+PLANS\\s+BROKER\\s+DESCRIPTION"))
				Eventually(session).Should(Say("overview-service\\s+simple, complex\\s+service-list-command-broker\\s+Provides an overview of any service instances and bindings that have been created by a platform"))
			})
		})
	})
})
