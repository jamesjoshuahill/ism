package repositories_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeBrokerRepo", func() {
	var (
		repo           *KubeBrokerRepo
		existingBroker *v1alpha1.Broker
		resource       types.NamespacedName
	)

	BeforeEach(func() {
		resource = types.NamespacedName{Name: "broker-1", Namespace: "default"}

		existingBroker = &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resource.Name,
				Namespace: resource.Namespace,
			},
			Spec: v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "http://example.org/broker",
				Username: "admin",
				Password: "password",
			},
		}

		repo = NewKubeBrokerRepo(kubeClient)
	})

	AfterEach(func() {
		kubeClient.Delete(context.Background(), existingBroker)
	})

	Describe("Get", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns broker", func() {
				broker, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(broker).To(Equal(existingBroker))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				_, err := repo.Get(resource)

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})

	Describe("UpdateStatus", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())
			})

			It("updates status", func() {
				newState := v1alpha1.BrokerStateRegistered
				Expect(existingBroker.Status.State).NotTo(Equal(newState))

				err := repo.UpdateState(existingBroker, newState)
				Expect(err).NotTo(HaveOccurred())

				updatedBroker, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatedBroker.Status.State).To(Equal(newState))
				Expect(existingBroker.Status.State).To(Equal(newState))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				newState := v1alpha1.BrokerStateRegistered
				err := repo.UpdateState(existingBroker, newState)

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})
})
