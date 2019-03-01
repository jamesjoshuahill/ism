package reconcilers_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/reconcilers"
	"github.com/pivotal-cf/ism/pkg/internal/reconcilers/reconcilersfakes"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("BrokerReconciler", func() {
	var (
		reconciler *BrokerReconciler
		err        error

		createBrokerClient         osbapi.CreateFunc
		brokerClientConfiguredWith *osbapi.ClientConfiguration

		returnedBroker v1alpha1.Broker

		fakeBrokerClient    *reconcilersfakes.FakeBrokerClient
		fakeKubeBrokerRepo  *reconcilersfakes.FakeKubeBrokerRepo
		fakeKubeServiceRepo *reconcilersfakes.FakeKubeServiceRepo
		fakeKubePlanRepo    *reconcilersfakes.FakeKubePlanRepo

		catalogPlanOne   = osbapi.Plan{ID: "id-plan-1", Name: "plan-1"}
		catalogPlanTwo   = osbapi.Plan{ID: "id-plan-2", Name: "plan-2"}
		catalogPlanThree = osbapi.Plan{ID: "id-plan-3", Name: "plan-3"}

		catalogServiceOne = osbapi.Service{
			ID:          "id-service-1",
			Name:        "service-1",
			Description: "some fancy description",
			Plans:       []osbapi.Plan{catalogPlanOne},
		}
		catalogServiceTwo = osbapi.Service{
			ID:          "id-service-2",
			Name:        "service-2",
			Description: "poorly written description",
			Plans:       []osbapi.Plan{catalogPlanTwo, catalogPlanThree},
		}
	)

	BeforeEach(func() {
		fakeBrokerClient = &reconcilersfakes.FakeBrokerClient{}
		fakeKubeBrokerRepo = &reconcilersfakes.FakeKubeBrokerRepo{}
		fakeKubeServiceRepo = &reconcilersfakes.FakeKubeServiceRepo{}
		fakeKubePlanRepo = &reconcilersfakes.FakeKubePlanRepo{}

		createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
			brokerClientConfiguredWith = config
			return fakeBrokerClient, nil
		}

		returnedBroker = v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "broker-1",
				Namespace: "default",
			},
			Spec: v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "broker-url",
				Username: "broker-username",
				Password: "broker-password",
			},
		}
		fakeKubeBrokerRepo.GetReturns(&returnedBroker, nil)
		fakeKubeServiceRepo.CreateReturnsOnCall(0, catalogServiceToBrokerService(&catalogServiceOne), nil)
		fakeKubeServiceRepo.CreateReturnsOnCall(1, catalogServiceToBrokerService(&catalogServiceTwo), nil)

		fakeBrokerClient.GetCatalogReturns(&osbapi.CatalogResponse{
			Services: []osbapi.Service{
				catalogServiceOne,
				catalogServiceTwo,
			},
		}, nil)
	})

	JustBeforeEach(func() {
		reconciler = NewBrokerReconciler(
			createBrokerClient,
			fakeKubeBrokerRepo,
			fakeKubeServiceRepo,
			fakeKubePlanRepo,
		)

		_, err = reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{Name: "broker-1", Namespace: "default"},
		})
	})

	It("fetches the broker resource using the kube broker repo", func() {
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeKubeBrokerRepo.GetCallCount()).To(Equal(1))
		namespacedName := fakeKubeBrokerRepo.GetArgsForCall(0)
		Expect(namespacedName).To(Equal(types.NamespacedName{Name: "broker-1", Namespace: "default"}))
	})

	It("configures the broker client with correct options", func() {
		Expect(*brokerClientConfiguredWith).To(Equal(osbapi.ClientConfiguration{
			Name:                "broker-1",
			URL:                 "broker-url",
			APIVersion:          osbapi.LatestAPIVersion(),
			TimeoutSeconds:      60,
			EnableAlphaFeatures: false,
			AuthConfig: &osbapi.AuthConfig{
				BasicAuthConfig: &osbapi.BasicAuthConfig{
					Username: "broker-username",
					Password: "broker-password",
				},
			},
		}))
	})

	It("fetches the catalog using the broker client", func() {
		Expect(fakeBrokerClient.GetCatalogCallCount()).To(Equal(1))
	})

	It("updates the broker status to registered", func() {
		Expect(fakeKubeBrokerRepo.UpdateStateCallCount()).To(Equal(1))
		broker, newState := fakeKubeBrokerRepo.UpdateStateArgsForCall(0)
		Expect(newState).To(Equal(v1alpha1.BrokerStateRegistered))
		Expect(*broker).To(Equal(returnedBroker))
	})

	It("creates service resources using the kube service repo", func() {
		broker, catalogService := fakeKubeServiceRepo.CreateArgsForCall(0)
		Expect(*broker).To(Equal(returnedBroker))
		Expect(catalogService).To(Equal(catalogServiceOne))

		broker, catalogService = fakeKubeServiceRepo.CreateArgsForCall(1)
		Expect(*broker).To(Equal(returnedBroker))
		Expect(catalogService).To(Equal(catalogServiceTwo))
	})

	It("creates plan resources using the kube plan repo", func() {
		serviceArg, catalogPlanArg := fakeKubePlanRepo.CreateArgsForCall(0)
		Expect(serviceArg.ObjectMeta.Name).To(Equal("service-1"))
		Expect(catalogPlanArg).To(Equal(catalogPlanOne))

		serviceArg, catalogPlanArg = fakeKubePlanRepo.CreateArgsForCall(1)
		Expect(serviceArg.ObjectMeta.Name).To(Equal("service-2"))
		Expect(catalogPlanArg).To(Equal(catalogPlanTwo))

		serviceArg, catalogPlanArg = fakeKubePlanRepo.CreateArgsForCall(2)
		Expect(serviceArg.ObjectMeta.Name).To(Equal("service-2"))
		Expect(catalogPlanArg).To(Equal(catalogPlanThree))
	})

	When("the broker state reports it is already registered", func() {
		BeforeEach(func() {
			returnedBroker.Status.State = v1alpha1.BrokerStateRegistered
		})

		It("doesn't call the broker", func() {
			Expect(fakeBrokerClient.GetCatalogCallCount()).To(Equal(0))
		})

		It("doesn't update the status", func() {
			Expect(fakeKubeBrokerRepo.UpdateStateCallCount()).To(Equal(0))
		})

		It("still reconciles successfully ", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("updating the broker status errors", func() {
		BeforeEach(func() {
			fakeKubeBrokerRepo.UpdateStateReturns(errors.New("error-updating-status"))
		})

		//TODO: test the state of service / plan creation here.
		It("returns the error", func() {
			Expect(err).To(MatchError("error-updating-status"))
		})
	})

	When("fetching the broker resource using the kube broker repo fails", func() {
		BeforeEach(func() {
			fakeKubeBrokerRepo.GetReturns(nil, errors.New("error-getting-broker"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-broker"))
		})
	})

	When("the broker has been deleted", func() {
		BeforeEach(func() {
			notFoundError := kerrors.NewNotFound(schema.GroupResource{}, "broker")
			fakeKubeBrokerRepo.GetReturns(nil, notFoundError)
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("configuring the broker client fails", func() {
		BeforeEach(func() {
			createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
				return nil, errors.New("error-configuring-broker-client")
			}
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-configuring-broker-client"))
		})
	})

	When("fetching the catalog using the broker client fails", func() {
		BeforeEach(func() {
			fakeBrokerClient.GetCatalogReturns(nil, errors.New("error-getting-catalog"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-catalog"))
		})
	})

	When("creating service resource fails", func() {
		BeforeEach(func() {
			fakeKubeServiceRepo.CreateReturnsOnCall(0, nil, errors.New("error-creating-service"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-creating-service"))
		})
	})

	When("creating service plan resource fails", func() {
		BeforeEach(func() {
			fakeKubePlanRepo.CreateReturnsOnCall(0, errors.New("error-creating-plan"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-creating-plan"))
		})
	})
})

func catalogServiceToBrokerService(osbapiService *osbapi.Service) *v1alpha1.BrokerService {
	return &v1alpha1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      osbapiService.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.BrokerServiceSpec{
			Name:        osbapiService.Name,
			Description: osbapiService.Description,
		},
	}
}
