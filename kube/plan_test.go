package kube_test

import (
	"context"

	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/pivotal-cf/ism/kube"
	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var _ = Describe("Plan", func() {
	var (
		kubeClient client.Client
		plan       *Plan
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		plan = &Plan{
			KubeClient: kubeClient,
		}
	})

	Describe("FindByService", func() {
		var (
			plans []*osbapi.Plan
			err   error
		)

		JustBeforeEach(func() {
			plans, err = plan.FindByService("service-uid-1")
		})

		When("plans contain owner references to services", func() {
			BeforeEach(func() {
				planResource := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-1",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "service-1",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "service-uid-1",
						}},
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name: "my-plan",
					},
				}
				Expect(kubeClient.Create(context.TODO(), planResource)).To(Succeed())

				planResource2 := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-2",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "service-2",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "service-uid-2",
						}},
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name: "my-plan-2",
					},
				}
				Expect(kubeClient.Create(context.TODO(), planResource2)).To(Succeed())
			})

			AfterEach(func() {
				deletePlans(kubeClient, "plan-1", "plan-2")
			})

			It("returns plans by service id", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(plans).To(Equal([]*osbapi.Plan{{
					Name:      "my-plan",
					ID:        "plan-1",
					ServiceID: "service-uid-1",
				}}))
			})
		})

		When("the plan owner reference is not set", func() {
			BeforeEach(func() {
				planResource := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-1",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name: "my-plan",
					},
				}
				Expect(kubeClient.Create(context.TODO(), planResource)).To(Succeed())
			})

			AfterEach(func() {
				deletePlans(kubeClient, "plan-1")
			})

			It("successfully returns no plans", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(plans).To(HaveLen(0))
			})
		})
	})
})

func deletePlans(kubeClient client.Client, planNames ...string) {
	for _, p := range planNames {
		pToDelete := &v1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), pToDelete)).To(Succeed())
	}
}
