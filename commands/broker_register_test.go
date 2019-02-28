package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Broker Register Command", func() {

	var (
		fakeUI              *commandsfakes.FakeUI
		fakeBrokerRegistrar *commandsfakes.FakeBrokerRegistrar

		registerCommand BrokerRegisterCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeUI = &commandsfakes.FakeUI{}
		fakeBrokerRegistrar = &commandsfakes.FakeBrokerRegistrar{}

		registerCommand = BrokerRegisterCommand{
			UI:              fakeUI,
			BrokerRegistrar: fakeBrokerRegistrar,
		}
	})

	JustBeforeEach(func() {
		executeErr = registerCommand.Execute(nil)
	})

	When("given all required args", func() {
		BeforeEach(func() {
			registerCommand.Name = "broker-1"
			registerCommand.URL = "test-url"
			registerCommand.Username = "test-username"
			registerCommand.Password = "test-password"
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays that the broker was registered", func() {
			text, data := fakeUI.DisplayTextArgsForCall(0)
			Expect(text).To(Equal("Broker '{{.BrokerName}}' registered."))
			Expect(data[0]).To(HaveKeyWithValue("BrokerName", "broker-1"))
		})

		It("registers the broker", func() {
			broker := fakeBrokerRegistrar.RegisterArgsForCall(0)

			Expect(broker).To(Equal(&osbapi.Broker{
				Name:     "broker-1",
				URL:      "test-url",
				Username: "test-username",
				Password: "test-password",
			}))
		})

		When("registering the broker errors", func() {
			BeforeEach(func() {
				fakeBrokerRegistrar.RegisterReturns(errors.New("error-registering-broker"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-registering-broker"))
			})
		})
	})
})
