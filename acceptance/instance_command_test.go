package acceptance

import (
	"encoding/json"
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
			})

			AfterEach(func() {
				deleteBrokers("instance-creation-broker")
				deleteServiceInstances("my-instance")
			})

			It("creates the service instance", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Instance 'my-instance' created\\."))

				//TODO Replace this once we have ism instance list
				out := runKubectl("get", "serviceinstance", "-o", "json")

				type instanceList struct {
					Items []interface{}
				}

				result := instanceList{}
				err := json.Unmarshal([]byte(out), &result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
			})
		})

		When("required args are not passed", func() {
			It("displays an informative message and exits 1", func() {
				Eventually(session).Should(Exit(1))
				Eventually(session).Should(Say("the required flags `--broker', `--name', `--plan' and `--service' were not specified"))
			})
		})
	})
})
