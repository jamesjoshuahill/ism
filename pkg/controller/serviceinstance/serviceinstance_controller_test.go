/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms
of the under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package serviceinstance

import (
	"context"
	"fmt"
	"sync"
	"time"

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
		mgr         manager.Manager
		mgrClient   client.Client
		mgrStopChan chan struct{}
		mgrStopWg   *sync.WaitGroup

		reconciler        reconcile.Reconciler
		reconcileRequests chan reconcile.Request
	)

	BeforeEach(func() {
		var err error
		mgr, err = manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())

		reconciler = newReconciler(mgr)
		mgrClient = mgr.GetClient()
	})

	JustBeforeEach(func() {
		var reconcileFunc reconcile.Reconciler
		reconcileFunc, reconcileRequests = SetupTestReconcile(reconciler)
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

	Describe("controller runs concurrently", func() {
		BeforeEach(func() {
			reconciler = reconcile.Func(func(r reconcile.Request) (reconcile.Result, error) {
				if r.Namespace == "concurrent" {
					time.Sleep(time.Hour * 24)
				}
				return reconcile.Result{}, nil
			})
		})

		It("can run at least max concurrent reconciles", func() {
			for i := 0; i < MaxConcurrentReconciles; i++ {
				name := fmt.Sprintf("concurrent-serviceinstance-%d", i)
				instance := &osbapiv1alpha1.ServiceInstance{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "concurrent"}}
				Expect(mgrClient.Create(context.TODO(), instance)).To(Succeed())

				Eventually(reconcileRequests).Should(Receive(Equal(
					reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: "concurrent"}},
				)))
			}
		})
	})

	Describe("Reconcile", func() {
		//TODO: Add an integration test running against "kinda real" kube and mock http server for broker.
	})
})
