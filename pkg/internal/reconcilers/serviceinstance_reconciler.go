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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// var ctx = context.TODO()

//go:generate counterfeiter . KubeServiceInstanceRepo

type KubeServiceInstanceRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.ServiceInstance, error)
	UpdateState(serviceinstance *v1alpha1.ServiceInstance, newState v1alpha1.ServiceInstanceState) error
}

type ServiceInstanceReconciler struct {
	log                     logr.Logger
	createBrokerClient      osbapi.CreateFunc
	kubeServiceInstanceRepo KubeServiceInstanceRepo
	kubeBrokerRepo          KubeBrokerRepo
}

func NewServiceInstanceReconciler(
	log logr.Logger,
	createBrokerClient osbapi.CreateFunc,
	kubeServiceInstanceRepo KubeServiceInstanceRepo,
	kubeBrokerRepo KubeBrokerRepo,
) *ServiceInstanceReconciler {
	return &ServiceInstanceReconciler{
		log:                     log,
		createBrokerClient:      createBrokerClient,
		kubeServiceInstanceRepo: kubeServiceInstanceRepo,
		kubeBrokerRepo:          kubeBrokerRepo,
	}
}

func (r *ServiceInstanceReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	instance, err := r.kubeServiceInstanceRepo.Get(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	broker, err := r.kubeBrokerRepo.Get(types.NamespacedName{Name: instance.Spec.BrokerName, Namespace: instance.ObjectMeta.Namespace})
	if err != nil {
		//TODO what if the broker does not exist?
		return reconcile.Result{}, err
	}

	if instance.Status.State == v1alpha1.ServiceInstanceStateProvisioned {
		return reconcile.Result{}, nil
	}

	osbapiConfig := brokerClientConfig(broker)

	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	_, err = osbapiClient.ProvisionInstance(&osbapi.ProvisionRequest{
		InstanceID:        string(instance.ObjectMeta.UID),
		AcceptsIncomplete: false,
		ServiceID:         instance.Spec.ServiceID,
		PlanID:            instance.Spec.PlanID,
		OrganizationGUID:  instance.ObjectMeta.Namespace,
		SpaceGUID:         instance.ObjectMeta.Namespace,
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	r.log.Info("Instance Provisioned")

	if err := r.kubeServiceInstanceRepo.UpdateState(instance, v1alpha1.ServiceInstanceStateProvisioned); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
