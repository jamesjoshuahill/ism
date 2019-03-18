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

	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//go:generate counterfeiter . KubeServiceBindingRepo

type KubeServiceBindingRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.ServiceBinding, error)
	UpdateState(servicebinding *v1alpha1.ServiceBinding, newState v1alpha1.ServiceBindingState) error
}

type ServiceBindingReconciler struct {
	createBrokerClient     osbapi.CreateFunc
	kubeServiceBindingRepo KubeServiceBindingRepo
	kubeBrokerRepo         KubeBrokerRepo
}

func NewServiceBindingReconciler(
	createBrokerClient osbapi.CreateFunc,
	kubeServiceBindingRepo KubeServiceBindingRepo,
	kubeBrokerRepo KubeBrokerRepo,
) *ServiceBindingReconciler {
	return &ServiceBindingReconciler{
		createBrokerClient:     createBrokerClient,
		kubeServiceBindingRepo: kubeServiceBindingRepo,
		kubeBrokerRepo:         kubeBrokerRepo,
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

	osbapiConfig := brokerClientConfig(broker)

	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	_, err = osbapiClient.Bind(&osbapi.BindRequest{
		BindingID:         string(binding.ObjectMeta.UID),
		AcceptsIncomplete: false,
		InstanceID:        binding.Spec.InstanceID,
		ServiceID:         binding.Spec.ServiceID,
		PlanID:            binding.Spec.PlanID,
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.kubeServiceBindingRepo.UpdateState(binding, v1alpha1.ServiceBindingStateCreated); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
