package integration

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI", func() {
	Describe("the help command", func() {
		It("displays help for the sm cli and exits 0", func() {
			command := exec.Command(pathToSMCLI, "--help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(Say(`usage: sm \[<flags>\]`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("CLI to interact with the Services Marketplace"))
		})
	})
})
