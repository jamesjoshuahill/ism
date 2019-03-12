package kube_test

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/pivotal-cf/ism/kube"
	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var _ = Describe("Instance", func() {

	var (
		kubeClient client.Client

		instance *Instance
		// registrationTimeout time.Duration
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		// registrationTimeout = time.Second

		instance = &Instance{
			KubeClient: kubeClient,
			// RegistrationTimeout: registrationTimeout,
		}
	})

	Describe("Create", func() {
		var (
			err error
			// registrationDuration time.Duration
		)

		JustBeforeEach(func() {
			b := &osbapi.Instance{
				Name:       "instance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			}

			// before := Createtime.Now()
			err = instance.Create(b)
			// registrationDuration = time.Since(before)
		})

		AfterEach(func() {
			deleteInstances(kubeClient, "instance-1")
		})

		When("the controller reacts to the instance", func() {
			// var closeChan chan bool

			BeforeEach(func() {
				// closeChan = make(chan bool)
				// go simulateRegistration(kubeClient, "instance-1", closeChan)
			})

			AfterEach(func() {
				// closeChan <- true
			})

			It("creates a new ServiceInstance resource", func() {
				Expect(err).NotTo(HaveOccurred())

				key := types.NamespacedName{
					Name:      "instance-1",
					Namespace: "default",
				}

				fetched := &v1alpha1.ServiceInstance{}
				Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())

				Expect(fetched.Spec).To(Equal(v1alpha1.ServiceInstanceSpec{
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "broker-1",
				}))
			})

			When("creating a new Instance fails", func() {
				BeforeEach(func() {
					// register the instance first, so that the second register errors
					b := &osbapi.Instance{
						Name:       "instance-1",
						PlanID:     "plan-1",
						ServiceID:  "service-1",
						BrokerName: "broker-1",
					}

					Expect(instance.Create(b)).To(Succeed())
				})

				It("propagates the error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		// When("the status of a instance is never set to registered", func() {
		// 	It("should eventually timeout", func() {
		// 		Expect(err).To(MatchError("timed out waiting for instance 'instance-1' to be registered"))
		// 	})
		//
		// 	It("times out once the timeout has been reached", func() {
		// 		estimatedExecutionTime := time.Second * 5 // flake prevention!
		//
		// 		Expect(registrationDuration).To(BeNumerically(">", registrationTimeout))
		// 		Expect(registrationDuration).To(BeNumerically("<", registrationTimeout+estimatedExecutionTime))
		// 	})
		// })
	})
})

func createdAtForInstance(kubeClient client.Client, instanceResource *v1alpha1.ServiceInstance) string {
	b := &v1alpha1.ServiceInstance{}
	namespacedName := types.NamespacedName{Name: instanceResource.Name, Namespace: instanceResource.Namespace}

	Expect(kubeClient.Get(context.TODO(), namespacedName, b)).To(Succeed())

	time := b.ObjectMeta.CreationTimestamp.String()
	return time
}

func deleteInstances(kubeClient client.Client, instanceNames ...string) {
	for _, b := range instanceNames {
		bToDelete := &v1alpha1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      b,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), bToDelete)).To(Succeed())
	}
}

// func simulateRegistration(kubeClient client.Client, instanceName string, done chan bool) {
// 	for {
// 		select {
// 		case <-done:
// 			return //exit func
// 		default:
// 			key := types.NamespacedName{
// 				Name:      instanceName,
// 				Namespace: "default",
// 			}
// 			instance := &v1alpha1.ServiceInstance{}
// 			err := kubeClient.Get(context.TODO(), key, instance)
// 			if err != nil {
// 				break //loop again
// 			}
//
// 			instance.Status.State = v1alpha1.ServiceInstanceStateCreateed
// 			Expect(kubeClient.Status().Update(context.TODO(), instance)).To(Succeed())
// 		}
// 	}
// }
