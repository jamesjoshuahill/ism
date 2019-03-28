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
	osbapiv1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"github.com/pivotal-cf/ism/pkg/internal/reconcilers"
	"github.com/pivotal-cf/ism/pkg/internal/repositories"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const MaxConcurrentReconciles = 10

var log = logf.Log.WithName("serviceinstance-controller")

// Add creates a new ServiceInstance Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileServiceInstance{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("serviceinstance-controller", mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: MaxConcurrentReconciles})
	if err != nil {
		return err
	}

	// Watch for changes to ServiceInstance
	err = c.Watch(&source.Kind{Type: &osbapiv1alpha1.ServiceInstance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileServiceInstance{}

// ReconcileServiceInstance reconciles a ServiceInstance object
type ReconcileServiceInstance struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ServiceInstance object and makes changes based on the state read
// and what is in the ServiceInstance.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=osbapi.ism.io,resources=serviceinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=osbapi.ism.io,resources=serviceinstances/status,verbs=get;update;patch
func (r *ReconcileServiceInstance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	requestLogger := log.WithValues("request", request)
	requestLogger.Info("Reconcile called")

	kubeBrokerRepo := repositories.NewKubeBrokerRepo(r.Client)
	kubeServiceInstanceRepo := repositories.NewKubeServiceInstanceRepo(r.Client)

	reconciler := reconcilers.NewServiceInstanceReconciler(
		requestLogger,
		osbapi.NewClient,
		kubeServiceInstanceRepo,
		kubeBrokerRepo,
	)

	return reconciler.Reconcile(request)
}
