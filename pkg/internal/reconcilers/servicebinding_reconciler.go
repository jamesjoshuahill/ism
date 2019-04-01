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
	"github.com/pivotal-cf/ism/pkg/finalizer"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	reconciler       = "servicebinding.osbapi.ism.io"
	bindingFinalizer = "finalizer." + reconciler
)

//go:generate counterfeiter . KubeServiceBindingRepo

type KubeServiceBindingRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.ServiceBinding, error)
	UpdateState(servicebinding *v1alpha1.ServiceBinding, newState v1alpha1.ServiceBindingState) error
	Update(servicebinding *v1alpha1.ServiceBinding) error
}

//go:generate counterfeiter . KubeSecretRepo

type KubeSecretRepo interface {
	Create(binding *v1alpha1.ServiceBinding, creds map[string]interface{}) (*corev1.Secret, error)
	Delete(corev1 *corev1.Secret) error
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
		return reconcile.Result{}, err
	}

	if !binding.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.handleDelete(broker, binding); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if binding.Status.State == "" {
		if err := r.handleCreate(broker, binding); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ServiceBindingReconciler) handleCreate(broker *v1alpha1.Broker, binding *v1alpha1.ServiceBinding) error {
	if err := r.kubeServiceBindingRepo.UpdateState(binding, v1alpha1.ServiceBindingStateCreating); err != nil {
		return err
	}

	credentials, err := r.bind(broker, binding)
	if err != nil {
		return err
	}

	secret, err := r.kubeSecretRepo.Create(binding, credentials)
	if err != nil {
		//TODO How to handle this?
		return err
	}

	finalizer.AddFinalizer(binding, bindingFinalizer)
	if err := r.kubeServiceBindingRepo.Update(binding); err != nil {
		return err
	}

	binding.Status.SecretRef = corev1.LocalObjectReference{Name: secret.Name}
	if err := r.kubeServiceBindingRepo.UpdateState(binding, v1alpha1.ServiceBindingStateCreated); err != nil {
		return err
	}

	r.log.Info("Binding created")
	return nil
}

func (r *ServiceBindingReconciler) handleDelete(broker *v1alpha1.Broker, binding *v1alpha1.ServiceBinding) error {
	if err := r.kubeServiceBindingRepo.UpdateState(binding, v1alpha1.ServiceBindingStateDeleting); err != nil {
		return err
	}

	if err := r.unbind(broker, binding); err != nil {
		return err
	}

	err := r.kubeSecretRepo.Delete(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      binding.Status.SecretRef.Name,
			Namespace: binding.Namespace,
		},
	})
	if err != nil {
		return err
	}

	finalizer.RemoveFinalizer(binding, bindingFinalizer)
	if err := r.kubeServiceBindingRepo.Update(binding); err != nil {
		return err
	}

	r.log.Info("Binding deleted")

	return nil
}

func (r *ServiceBindingReconciler) unbind(broker *v1alpha1.Broker, binding *v1alpha1.ServiceBinding) error {
	osbapiConfig := brokerClientConfig(broker)

	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return err
	}

	_, err = osbapiClient.Unbind(&osbapi.UnbindRequest{
		BindingID:         string(binding.ObjectMeta.UID),
		AcceptsIncomplete: false,
		InstanceID:        binding.Spec.InstanceID,
		ServiceID:         binding.Spec.ServiceID,
		PlanID:            binding.Spec.PlanID,
	})

	if err != nil {
		return err
	}

	return nil
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
