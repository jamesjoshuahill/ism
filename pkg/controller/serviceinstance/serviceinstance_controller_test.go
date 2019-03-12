package serviceinstance

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	osbapiv1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("ServiceInstance Controller", func() {
	var (
		mgrClient   client.Client
		mgrStopChan chan struct{}
		mgrStopWg   *sync.WaitGroup

		reconcileRequests chan reconcile.Request
	)

	BeforeEach(func() {
		var err error
		var reconcileFunc reconcile.Reconciler

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())

		mgrClient = mgr.GetClient()

		reconcileFunc, reconcileRequests = SetupTestReconcile(newReconciler(mgr))
		Expect(add(mgr, reconcileFunc)).To(Succeed())

		mgrStopChan, mgrStopWg = StartTestManager(mgr)
	})

	AfterEach(func() {
		close(mgrStopChan)
		mgrStopWg.Wait()
	})

	When("a serviceinstance resource is created", func() {
		It("calls the reconcile function", func() {
			instance := &osbapiv1alpha1.ServiceInstance{ObjectMeta: metav1.ObjectMeta{Name: "serviceinstance-1", Namespace: "default"}}
			Expect(mgrClient.Create(context.TODO(), instance)).To(Succeed())

			Eventually(reconcileRequests).Should(Receive(Equal(
				reconcile.Request{NamespacedName: types.NamespacedName{Name: "serviceinstance-1", Namespace: "default"}},
			)))
		})
	})

	Describe("Reconcile", func() {
		//TODO: Add an integration test running against "kinda real" kube and mock http server for broker.
	})
})
