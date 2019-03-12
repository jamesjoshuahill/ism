package usecases_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/ism/osbapi"
	. "github.com/pivotal-cf/ism/usecases"
	"github.com/pivotal-cf/ism/usecases/usecasesfakes"
)

var _ = Describe("Instance Create Usecase", func() {

	var (
		fakeBrokerFetcher   *usecasesfakes.FakeBrokerFetcher
		fakeServiceFetcher  *usecasesfakes.FakeServiceFetcher
		fakePlanFetcher     *usecasesfakes.FakePlanFetcher
		fakeInstanceCreator *usecasesfakes.FakeInstanceCreator

		instanceCreateUsecase InstanceCreateUsecase

		brokers  []*osbapi.Broker
		services []*osbapi.Service
		plans    []*osbapi.Plan

		executeErr error
	)

	BeforeEach(func() {
		fakeBrokerFetcher = &usecasesfakes.FakeBrokerFetcher{}
		fakeServiceFetcher = &usecasesfakes.FakeServiceFetcher{}
		fakePlanFetcher = &usecasesfakes.FakePlanFetcher{}
		fakeInstanceCreator = &usecasesfakes.FakeInstanceCreator{}

		instanceCreateUsecase = InstanceCreateUsecase{
			BrokerFetcher:   fakeBrokerFetcher,
			ServiceFetcher:  fakeServiceFetcher,
			PlanFetcher:     fakePlanFetcher,
			InstanceCreator: fakeInstanceCreator,
		}

		brokers = []*osbapi.Broker{{
			Name: "another-broker",
		}, {
			Name: "my-broker",
		}}

		services = []*osbapi.Service{{
			Name:       "my-service",
			ID:         "service-1",
			BrokerName: "my-broker",
		}, {
			Name:       "another-service",
			ID:         "service-2",
			BrokerName: "my-broker",
		}}

		plans = []*osbapi.Plan{{
			Name:      "my-plan",
			ID:        "plan-1",
			ServiceID: "service-1",
		}, {
			Name:      "another-plan",
			ID:        "plan-2",
			ServiceID: "service-1",
		}}
	})

	JustBeforeEach(func() {
		executeErr = instanceCreateUsecase.Create("my-instance", "my-plan", "my-service", "my-broker")
	})

	When("passed valid args", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns(brokers, nil)
			fakeServiceFetcher.GetServicesReturns(services, nil)
			fakePlanFetcher.GetPlansReturns(plans, nil)
		})

		It("creates an instance", func() {
			Expect(fakeBrokerFetcher.GetBrokersCallCount()).To(Equal(1))
			Expect(fakeServiceFetcher.GetServicesCallCount()).To(Equal(1))
			Expect(fakePlanFetcher.GetPlansCallCount()).To(Equal(1))

			passedBrokerName := fakeServiceFetcher.GetServicesArgsForCall(0)
			Expect(passedBrokerName).To(Equal("my-broker"))

			passedServiceID := fakePlanFetcher.GetPlansArgsForCall(0)
			Expect(passedServiceID).To(Equal("service-1"))

			Expect(fakeInstanceCreator.CreateCallCount()).To(Equal(1))
			passedInstance := fakeInstanceCreator.CreateArgsForCall(0)

			Expect(*passedInstance).To(Equal(osbapi.Instance{
				Name:       "my-instance",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "my-broker",
			}))

			Expect(executeErr).NotTo(HaveOccurred())
		})
	})

	When("there are no brokers", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{}, nil)
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Broker 'my-broker' does not exist"))
		})
	})

	When("there are no services", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns(brokers, nil)
			fakeServiceFetcher.GetServicesReturns([]*osbapi.Service{}, nil)
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Service 'my-service' does not exist"))
		})
	})

	When("there are no plans", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns(brokers, nil)
			fakeServiceFetcher.GetServicesReturns(services, nil)
			fakePlanFetcher.GetPlansReturns([]*osbapi.Plan{}, nil)
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Plan 'my-plan' does not exist"))
		})
	})

	When("fetching brokers errors", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{}, errors.New("error-getting-brokers"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-brokers"))
		})
	})

	When("fetching services errors", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns(brokers, nil)
			fakeServiceFetcher.GetServicesReturns([]*osbapi.Service{}, errors.New("error-getting-services"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-services"))
		})
	})

	When("fetching plans errors", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns(brokers, nil)
			fakeServiceFetcher.GetServicesReturns(services, nil)
			fakePlanFetcher.GetPlansReturns([]*osbapi.Plan{}, errors.New("error-getting-plans"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-plans"))
		})
	})

	When("creating the instance errors", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns(brokers, nil)
			fakeServiceFetcher.GetServicesReturns(services, nil)
			fakePlanFetcher.GetPlansReturns(plans, nil)
			fakeInstanceCreator.CreateReturns(errors.New("error-creating-instance"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-creating-instance"))
		})
	})
})
