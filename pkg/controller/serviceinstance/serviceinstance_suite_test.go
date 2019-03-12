package serviceinstance

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pivotal-cf/ism/pkg/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var cfg *rest.Config

func TestServiceInstance(t *testing.T) {
	var testEnv *envtest.Environment

	SetDefaultEventuallyTimeout(time.Second * 5)
	SetDefaultConsistentlyDuration(time.Second * 5)

	BeforeSuite(func() {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crds")},
		}
		apis.AddToScheme(scheme.Scheme)

		var err error
		cfg, err = testEnv.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		testEnv.Stop()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceInstance Suite")
}

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)

	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		requests <- req
		return result, err
	})

	return fn, requests
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager) (chan struct{}, *sync.WaitGroup) {
	stopChan := make(chan struct{})
	wg := &sync.WaitGroup{}

	go func() {
		wg.Add(1)
		Expect(mgr.Start(stopChan)).To(Succeed())
		wg.Done()
	}()

	return stopChan, wg
}
