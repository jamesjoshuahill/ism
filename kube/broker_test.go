package kube_test

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/pivotal-cf/ism/kube"
	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var _ = Describe("Broker", func() {

	var (
		kubeClient client.Client

		broker *Broker
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		broker = &Broker{
			KubeClient: kubeClient,
		}
	})

	Describe("Register", func() {
		var err error

		JustBeforeEach(func() {
			b := &osbapi.Broker{
				Name:     "broker-1",
				URL:      "broker-1-url",
				Username: "broker-1-username",
				Password: "broker-1-password",
			}

			err = broker.Register(b)
		})

		AfterEach(func() {
			deleteBrokers(kubeClient, "broker-1")
		})

		It("creates a new Broker resource instance", func() {
			Expect(err).NotTo(HaveOccurred())

			key := types.NamespacedName{
				Name:      "broker-1",
				Namespace: "default",
			}

			fetched := &v1alpha1.Broker{}
			Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())

			Expect(fetched.Spec).To(Equal(v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "broker-1-url",
				Username: "broker-1-username",
				Password: "broker-1-password",
			}))
		})

		When("creating a new Broker fails", func() {
			BeforeEach(func() {
				// register the broker first, so that the second register errors
				b := &osbapi.Broker{
					Name:     "broker-1",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				}

				Expect(broker.Register(b)).To(Succeed())
			})

			It("propagates the error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FindAll", func() {
		var (
			brokers         []*osbapi.Broker
			brokerCreatedAt string
			err             error
		)

		BeforeEach(func() {
			brokerResource := &v1alpha1.Broker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "broker-1",
					Namespace: "default",
				},
				Spec: v1alpha1.BrokerSpec{
					Name:     "broker-1",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				},
			}

			Expect(kubeClient.Create(context.TODO(), brokerResource)).To(Succeed())
			brokerCreatedAt = createdAtForBroker(kubeClient, brokerResource)
		})

		JustBeforeEach(func() {
			brokers, err = broker.FindAll()
		})

		AfterEach(func() {
			deleteBrokers(kubeClient, "broker-1")
		})

		It("returns all brokers", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(*brokers[0]).To(MatchFields(IgnoreExtras, Fields{
				"CreatedAt": Equal(brokerCreatedAt),
				"Name":      Equal("broker-1"),
				"URL":       Equal("broker-1-url"),
				"Username":  Equal("broker-1-username"),
				"Password":  Equal("broker-1-password"),
			}))
		})
	})
})

func createdAtForBroker(kubeClient client.Client, brokerResource *v1alpha1.Broker) string {
	b := &v1alpha1.Broker{}
	namespacedName := types.NamespacedName{Name: brokerResource.Name, Namespace: brokerResource.Namespace}

	Expect(kubeClient.Get(context.TODO(), namespacedName, b)).To(Succeed())

	time := b.ObjectMeta.CreationTimestamp.String()
	return time
}

func deleteBrokers(kubeClient client.Client, brokerNames ...string) {
	for _, b := range brokerNames {
		bToDelete := &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      b,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), bToDelete)).To(Succeed())
	}
}
