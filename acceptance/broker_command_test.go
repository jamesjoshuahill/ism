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
			Eventually(session).Should(Say("The broker command group lets you register and list service brokers from the"))
			Eventually(session).Should(Say("marketplace"))
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
				deleteBrokers("my-broker")
			})

			It("displays a message that the registration has been successful", func() {
				Eventually(session).Should(Exit(0))
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
				Eventually(session).Should(Say(`ism \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a service broker into the marketplace"))
			})
		})

		When("required arguments are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session).Should(Say("the required flags `--name', `--password', `--url' and `--username' were not specified"))
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
				Eventually(session).Should(Say("Lists the service brokers in the marketplace"))
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
				registerBroker("test-broker-2")
			})

			AfterEach(func() {
				deleteBrokers("test-broker-2")
			})

			It("displays the broker", func() {
				timeRegex := `\d{4,}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("NAME\\s+URL\\s+CREATED AT"))
				Eventually(session).Should(Say("test-broker-2\\s+" + nodeBrokerURL + "\\s+" + timeRegex))
			})
		})
	})
})
