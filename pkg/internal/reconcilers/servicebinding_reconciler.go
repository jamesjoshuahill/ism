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

package reconcilers

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//go:generate counterfeiter . KubeServiceBindingRepo

type KubeServiceBindingRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.ServiceBinding, error)
	UpdateState(servicebinding *v1alpha1.ServiceBinding, newState v1alpha1.ServiceBindingState) error
}

//go:generate counterfeiter . KubeSecretRepo

type KubeSecretRepo interface {
	Create(binding *v1alpha1.ServiceBinding, creds map[string]interface{}) (*corev1.Secret, error)
}

type ServiceBindingReconciler struct {
	log                    logr.Logger
	createBrokerClient     osbapi.CreateFunc
	kubeServiceBindingRepo KubeServiceBindingRepo
	kubeBrokerRepo         KubeBrokerRepo
	kubeSecretRepo         KubeSecretRepo
}

func NewServiceBindingReconciler(
	log logr.Logger,
	createBrokerClient osbapi.CreateFunc,
	kubeServiceBindingRepo KubeServiceBindingRepo,
	kubeBrokerRepo KubeBrokerRepo,
	kubeSecretRepo KubeSecretRepo,
) *ServiceBindingReconciler {
	return &ServiceBindingReconciler{
		log:                    log,
		createBrokerClient:     createBrokerClient,
		kubeServiceBindingRepo: kubeServiceBindingRepo,
		kubeBrokerRepo:         kubeBrokerRepo,
		kubeSecretRepo:         kubeSecretRepo,
	}
}

func (r *ServiceBindingReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	binding, err := r.kubeServiceBindingRepo.Get(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	broker, err := r.kubeBrokerRepo.Get(types.NamespacedName{Name: binding.Spec.BrokerName, Namespace: binding.ObjectMeta.Namespace})
	if err != nil {
		//TODO what if the broker does not exist?
		return reconcile.Result{}, err
	}

	if binding.Status.State == v1alpha1.ServiceBindingStateCreated {
		return reconcile.Result{}, nil
	}

	credentials, err := r.bind(broker, binding)
	if err != nil {
		return reconcile.Result{}, err
	}

	secret, err := r.kubeSecretRepo.Create(binding, credentials)
	if err != nil {
		//TODO How to handle this?
		return reconcile.Result{}, err
	}

	binding.Status.SecretRef = corev1.LocalObjectReference{Name: secret.Name}

	if err := r.kubeServiceBindingRepo.UpdateState(binding, v1alpha1.ServiceBindingStateCreated); err != nil {
		return reconcile.Result{}, err
	}

	r.log.Info("Binding created")

	return reconcile.Result{}, nil
}

func (r *ServiceBindingReconciler) bind(broker *v1alpha1.Broker, binding *v1alpha1.ServiceBinding) (map[string]interface{}, error) {
	osbapiConfig := brokerClientConfig(broker)

	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return nil, err
	}

	osbapiBinding, err := osbapiClient.Bind(&osbapi.BindRequest{
		BindingID:         string(binding.ObjectMeta.UID),
		AcceptsIncomplete: false,
		InstanceID:        binding.Spec.InstanceID,
		ServiceID:         binding.Spec.ServiceID,
		PlanID:            binding.Spec.PlanID,
	})
	if err != nil {
		return nil, err
	}

	return osbapiBinding.Credentials, nil
}
