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
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var ctx = context.TODO()

//go:generate counterfeiter . KubeBrokerRepo

type KubeBrokerRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.Broker, error)
	UpdateState(broker *v1alpha1.Broker, newState v1alpha1.BrokerState, message string) error
}

//go:generate counterfeiter . KubeServiceRepo

type KubeServiceRepo interface {
	Create(broker *v1alpha1.Broker, catalogService osbapi.Service) (*v1alpha1.BrokerService, error)
}

//go:generate counterfeiter . KubePlanRepo

type KubePlanRepo interface {
	Create(brokerService *v1alpha1.BrokerService, catalogPlan osbapi.Plan) error
}

//go:generate counterfeiter . BrokerClient

type BrokerClient interface {
	osbapi.Client
}

type BrokerReconciler struct {
	log                logr.Logger
	kubeBrokerRepo     KubeBrokerRepo
	kubeServiceRepo    KubeServiceRepo
	kubePlanRepo       KubePlanRepo
	createBrokerClient osbapi.CreateFunc
}

func NewBrokerReconciler(
	log logr.Logger,
	createBrokerClient osbapi.CreateFunc,
	kubeBrokerRepo KubeBrokerRepo,
	kubeServiceRepo KubeServiceRepo,
	kubePlanRepo KubePlanRepo,
) *BrokerReconciler {
	return &BrokerReconciler{
		log:                log,
		createBrokerClient: createBrokerClient,
		kubeBrokerRepo:     kubeBrokerRepo,
		kubeServiceRepo:    kubeServiceRepo,
		kubePlanRepo:       kubePlanRepo,
	}
}

func (r *BrokerReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	broker, err := r.kubeBrokerRepo.Get(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if broker.Status.State == v1alpha1.BrokerStateRegistered {
		return reconcile.Result{}, nil
	}

	catalog, err := r.catalog(broker)
	if err != nil {
		message := messageForError(err)

		if err := r.kubeBrokerRepo.UpdateState(broker, v1alpha1.BrokerStateRegistrationFailed, message); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, err
	}

	for _, catalogService := range catalog.Services {
		service, err := r.kubeServiceRepo.Create(broker, catalogService)
		if err != nil {
			return reconcile.Result{}, err
		}

		for _, catalogPlan := range catalogService.Plans {
			if err := r.kubePlanRepo.Create(service, catalogPlan); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	if err := r.kubeBrokerRepo.UpdateState(broker, v1alpha1.BrokerStateRegistered, ""); err != nil {
		return reconcile.Result{}, err
	}

	r.log.Info("Broker registered")

	return reconcile.Result{}, nil
}

func (r *BrokerReconciler) catalog(broker *v1alpha1.Broker) (*osbapi.CatalogResponse, error) {
	osbapiConfig := brokerClientConfig(broker)
	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return nil, err
	}

	return osbapiClient.GetCatalog()
}

func brokerClientConfig(broker *v1alpha1.Broker) *osbapi.ClientConfiguration {
	osbapiConfig := osbapi.DefaultClientConfiguration()
	osbapiConfig.Name = broker.Spec.Name
	osbapiConfig.URL = broker.Spec.URL
	osbapiConfig.AuthConfig = &osbapi.AuthConfig{
		BasicAuthConfig: &osbapi.BasicAuthConfig{
			Username: broker.Spec.Username,
			Password: broker.Spec.Password,
		},
	}
	return osbapiConfig
}

func messageForError(err error) string {
	httpErr, ok := osbapi.IsHTTPError(err)

	if ok {
		switch httpErr.StatusCode {
		case 200:
			return "Service broker did not return a valid catalog"
		case 401:
			return "Service broker authentication failed"
		case 404:
			return "Service broker catalog not found"
		}
	}

	return "Unknown error"
}
