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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI", func() {
	var homePath string

	BeforeEach(func() {
		var err error
		homePath, err = ioutil.TempDir("", "ism-root-command-acceptance-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(homePath)).To(Succeed())
	})

	When("no command or flag is passed", func() {
		It("displays help and exits 0", func() {
			command := exec.Command(nodePathToCLI)
			command.Env = []string{fmt.Sprintf("HOME=%s", homePath)}
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`ism \[OPTIONS\] <command>`))
			Eventually(session).Should(Say(`Available commands:`))
			Eventually(session).Should(Say(`binding`))
			Eventually(session).Should(Say(`broker`))
			Eventually(session).Should(Say(`instance`))
			Eventually(session).Should(Say(`service`))
		})
	})

	When("--help is passed", func() {
		It("displays help and exits 0", func() {
			command := exec.Command(nodePathToCLI, "--help")
			command.Env = []string{fmt.Sprintf("HOME=%s", homePath)}
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`ism \[OPTIONS\] <command>`))
			Eventually(session).Should(Say(`Available commands:`))
			Eventually(session).Should(Say(`binding`))
			Eventually(session).Should(Say(`broker`))
			Eventually(session).Should(Say(`instance`))
			Eventually(session).Should(Say(`service`))
		})
	})
})
