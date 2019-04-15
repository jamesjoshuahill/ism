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

package kube

import (
	"context"
	"errors"
	"time"

	"github.com/pivotal-cf/ism/repositories"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Broker struct {
	KubeClient          client.Client
	RegistrationTimeout time.Duration
}

func (b *Broker) FindByName(name string) (*osbapi.Broker, error) {
	broker := &v1alpha1.Broker{}
	err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: "default"}, broker)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, repositories.ErrBrokerNotFound
		}
		return nil, err
	}

	osbapiBroker := &osbapi.Broker{
		Name:      broker.Spec.Name,
		URL:       broker.Spec.URL,
		Username:  broker.Spec.Username,
		Password:  broker.Spec.Password,
		CreatedAt: broker.ObjectMeta.CreationTimestamp.String(),
	}

	return osbapiBroker, nil
}

func (b *Broker) FindAll() ([]*osbapi.Broker, error) {
	list := &v1alpha1.BrokerList{}
	if err := b.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Broker{}, err
	}

	brokers := []*osbapi.Broker{}
	for _, broker := range list.Items {
		brokers = append(brokers, &osbapi.Broker{
			Name:      broker.Spec.Name,
			URL:       broker.Spec.URL,
			Username:  broker.Spec.Username,
			Password:  broker.Spec.Password,
			CreatedAt: broker.ObjectMeta.CreationTimestamp.String(),
		})
	}

	return brokers, nil
}

func (b *Broker) Delete(name string) error {
	brokerResource := &v1alpha1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}

	return b.KubeClient.Delete(context.TODO(), brokerResource)
}

func (b *Broker) Register(broker *osbapi.Broker) error {
	brokerResource := &v1alpha1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      broker.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.BrokerSpec{
			Name:     broker.Name,
			URL:      broker.URL,
			Username: broker.Username,
			Password: broker.Password,
		},
	}

	if err := b.KubeClient.Create(context.TODO(), brokerResource); err != nil {
		if kerrors.ReasonForError(err) == metav1.StatusReasonAlreadyExists {
			return repositories.ErrBrokerAlreadyExists{BrokerName: broker.Name}
		}
		return err
	}

	if err := b.waitForBrokerRegistration(brokerResource); err != nil {
		var _ = b.Delete(broker.Name) // attempt to delete the broker, but don't return an error if we can't
		return err
	}

	return nil
}

func (b *Broker) waitForBrokerRegistration(broker *v1alpha1.Broker) error {
	err := wait.Poll(time.Second/2, b.RegistrationTimeout, func() (bool, error) {
		fetchedBroker := &v1alpha1.Broker{}

		err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: broker.Name, Namespace: broker.Namespace}, fetchedBroker)
		if err == nil && fetchedBroker.Status.State == v1alpha1.BrokerStateRegistered {
			return true, nil
		}

		if err == nil && fetchedBroker.Status.State == v1alpha1.BrokerStateRegistrationFailed {
			return true, errors.New(fetchedBroker.Status.Message)
		}

		return false, nil
	})

	if err != nil {
		if err == wait.ErrWaitTimeout {
			return repositories.ErrBrokerRegisterTimeout{BrokerName: broker.Name}
		}

		return err
	}

	return nil
}
