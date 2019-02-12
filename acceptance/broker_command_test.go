package acceptance

import (
	"io"
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

		command := exec.Command(pathToSMCLI, args...)
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
			Eventually(session).Should(Say(`sm \[OPTIONS\] broker <register>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The broker command group lets you register, update and deregister Service"))
			Eventually(session).Should(Say("Brokers from the marketplace"))
		})
	})

	Describe("register sub command", func() {
		BeforeEach(func() {
			args = append(args, "register")
		})

		When("valid args are passed", func() {
			BeforeEach(func() {
				args = append(args, "--name", "my-broker", "--url", "url", "--username", "username", "--password", "password")
			})

			AfterEach(func() {
				deleteAllBrokers()
			})

			It("successfully registers the broker, and displays a message", func() {
				Eventually(session).Should(Exit(0))

				ensureBrokerExists("my-broker")

				Eventually(session).Should(Say("Broker 'my-broker' registered\\."))
			})
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a Service Broker into the marketplace"))
			})
		})

		When("required arguments are not passed", func() {
			It("displays an informative message and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("the required flags `--name', `--password', `--url' and `--username' were not specified"))
			})
		})
	})
})

func ensureBrokerExists(brokerName string) {
	outBuffer := NewBuffer()
	getBrokersCmd := exec.Command("kubectl", "get", "brokers")
	getBrokersCmd.Stdout = io.MultiWriter(outBuffer, GinkgoWriter)
	getBrokersCmd.Stderr = GinkgoWriter

	Expect(getBrokersCmd.Run()).To(Succeed())
	Expect(outBuffer).To(Say(brokerName))
}

func deleteAllBrokers() {
	deleteBrokersCmd := exec.Command("kubectl", "delete", "brokers", "--all")
	deleteBrokersCmd.Stdout = GinkgoWriter
	deleteBrokersCmd.Stderr = GinkgoWriter
	Expect(deleteBrokersCmd.Run()).To(Succeed())
}
