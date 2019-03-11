package kube_test

import (
	"context"

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

var _ = Describe("Service", func() {
	var (
		kubeClient client.Client
		service    *Service
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		service = &Service{
			KubeClient: kubeClient,
		}
	})

	Describe("FindByBroker", func() {
		var (
			services []*osbapi.Service
			err      error
		)

		JustBeforeEach(func() {
			services, err = service.FindByBroker("broker-1")
		})

		When("services contain owner references to brokers", func() {
			BeforeEach(func() {
				serviceResource := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "broker-1",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "broker-1",
						}},
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service-1",
						Description: "service-1-desc",
					},
				}
				Expect(kubeClient.Create(context.TODO(), serviceResource)).To(Succeed())

				serviceResource2 := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-2",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "broker-2",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "broker-2",
						}},
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service-2",
						Description: "service-2-desc",
					},
				}
				Expect(kubeClient.Create(context.TODO(), serviceResource2)).To(Succeed())
			})

			AfterEach(func() {
				deleteServices(kubeClient, "service-1", "service-2")
			})

			It("returns only the services owned by the broker id", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(services).To(HaveLen(1))
				Expect(*services[0]).To(MatchFields(IgnoreExtras, Fields{
					"ID":          Equal("service-1"),
					"Name":        Equal("my-service-1"),
					"Description": Equal("service-1-desc"),
					"BrokerID":    Equal("broker-1"),
				}))
			})
		})

		When("the service owner reference is not set", func() {
			BeforeEach(func() {
				serviceResource := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service-1",
						Description: "service-1-desc",
					},
				}
				Expect(kubeClient.Create(context.TODO(), serviceResource)).To(Succeed())
			})

			AfterEach(func() {
				deleteServices(kubeClient, "service-1")
			})

			It("successfully returns no services", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(services).To(HaveLen(0))
			})
		})
	})
})

func deleteServices(kubeClient client.Client, serviceNames ...string) {
	for _, s := range serviceNames {
		sToDelete := &v1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      s,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), sToDelete)).To(Succeed())
	}
}
