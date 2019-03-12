package acceptance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

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

				Expect(getBrokerData().ServiceInstances).To(HaveLen(0))
			})

			AfterEach(func() {
				deleteBrokers("instance-creation-broker")
				deleteServiceInstances("my-instance")
			})

			It("creates the service instance", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Instance 'my-instance' created\\."))

				time.Sleep(time.Second * 5)

				data := getBrokerData()
				Expect(data.ServiceInstances).To(HaveLen(1))

				for _, serviceInstance := range data.ServiceInstances {
					Expect(serviceInstance.ServiceName).To(Equal(serviceName))
					Expect(serviceInstance.PlanName).To(Equal(planName))
				}

				//TODO Add a test to check ism instance list
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

type serviceInstance struct {
	PlanName    string `json:"plan_name"`
	ServiceName string `json:"service_name"`
}

type brokerData struct {
	ServiceInstances map[string]serviceInstance `json:"serviceInstances"`
}

func getBrokerData() brokerData {
	brokerDataURL := fmt.Sprintf("%s/data", nodeBrokerURL)

	resp, err := http.Get(brokerDataURL)
	Expect(err).NotTo(HaveOccurred())
	respBytes, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	var data brokerData
	Expect(json.Unmarshal(respBytes, &data)).To(Succeed())

	return data
}
