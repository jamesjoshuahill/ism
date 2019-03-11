package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
)

var _ = Describe("Instance create command", func() {

	var (
		fakeInstanceCreateUsecase *commandsfakes.FakeInstanceCreateUsecase
		fakeUI                    *commandsfakes.FakeUI

		createCommand InstanceCreateCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeInstanceCreateUsecase = &commandsfakes.FakeInstanceCreateUsecase{}
		fakeUI = &commandsfakes.FakeUI{}

		createCommand = InstanceCreateCommand{
			Name:                  "instance-1",
			Plan:                  "plan-1",
			Service:               "service-1",
			Broker:                "broker-1",
			InstanceCreateUsecase: fakeInstanceCreateUsecase,
			UI:                    fakeUI,
		}
	})

	JustBeforeEach(func() {
		executeErr = createCommand.Execute(nil)
	})

	When("creating an instance succeeds", func() {
		BeforeEach(func() {
			fakeInstanceCreateUsecase.CreateReturns(nil)
		})

		It("calls to create the instance", func() {
			name, planName, serviceName, brokerName := fakeInstanceCreateUsecase.CreateArgsForCall(0)

			Expect(name).To(Equal("instance-1"))
			Expect(planName).To(Equal("plan-1"))
			Expect(serviceName).To(Equal("service-1"))
			Expect(brokerName).To(Equal("broker-1"))
		})

		It("displays that the service instance was created", func() {
			text, data := fakeUI.DisplayTextArgsForCall(0)
			Expect(text).To(Equal("Instance '{{.InstanceName}}' created."))
			Expect(data[0]).To(HaveKeyWithValue("InstanceName", "instance-1"))
		})
	})

	When("creating an instace errors", func() {
		BeforeEach(func() {
			fakeInstanceCreateUsecase.CreateReturns(errors.New("error-creating-instance"))
		})

		It("returns the error", func() {
			Expect(executeErr).To(MatchError("error-creating-instance"))
		})
	})
})
