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
	"encoding/json"
	"errors"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var errBindingNotFound = errors.New("binding not found")

type Binding struct {
	KubeClient client.Client
}

func (b *Binding) Create(binding *osbapi.Binding) error {
	bindingResource := &v1alpha1.ServiceBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      binding.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.ServiceBindingSpec{
			Name:       binding.Name,
			InstanceID: binding.InstanceID,
			PlanID:     binding.PlanID,
			ServiceID:  binding.ServiceID,
			BrokerName: binding.BrokerName,
		},
	}

	return b.KubeClient.Create(context.TODO(), bindingResource)
}

func (b *Binding) FindByName(name string) (*osbapi.Binding, error) {
	binding := &v1alpha1.ServiceBinding{}
	err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: "default"}, binding)

	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, errBindingNotFound
		}
		return nil, err
	}

	var decodedCreds map[string]interface{}
	if binding.Status.State == v1alpha1.ServiceBindingStateCreated {
		decodedCreds, err = b.getCredentials(name)
		if err != nil {
			return nil, err
		}
	}

	osbapiBinding := &osbapi.Binding{
		ID:          string(binding.ObjectMeta.UID),
		Name:        binding.Spec.Name,
		InstanceID:  binding.Spec.InstanceID,
		PlanID:      binding.Spec.PlanID,
		ServiceID:   binding.Spec.ServiceID,
		BrokerName:  binding.Spec.BrokerName,
		Status:      b.setStatus(binding.Status.State),
		CreatedAt:   binding.ObjectMeta.CreationTimestamp.String(),
		Credentials: decodedCreds,
	}

	return osbapiBinding, nil
}

func (b *Binding) FindAll() ([]*osbapi.Binding, error) {
	list := &v1alpha1.ServiceBindingList{}
	if err := b.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Binding{}, err
	}

	bindings := []*osbapi.Binding{}
	for _, binding := range list.Items {
		bindings = append(bindings, &osbapi.Binding{
			ID:         string(binding.ObjectMeta.UID),
			Name:       binding.Spec.Name,
			InstanceID: binding.Spec.InstanceID,
			PlanID:     binding.Spec.PlanID,
			ServiceID:  binding.Spec.ServiceID,
			BrokerName: binding.Spec.BrokerName,
			Status:     b.setStatus(binding.Status.State),
			CreatedAt:  binding.ObjectMeta.CreationTimestamp.String(),
			//TODO: Should we add credentials?
		})
	}

	return bindings, nil
}

func (b *Binding) setStatus(state v1alpha1.ServiceBindingState) string {
	if state == v1alpha1.ServiceBindingStateCreated {
		return created
	}
	return creating
}

func (b *Binding) getCredentials(name string) (map[string]interface{}, error) {
	secret := &corev1.Secret{}
	if err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: "default"}, secret); err != nil {
		return nil, err
	}

	encodedCreds, ok := secret.Data["credentials"]
	if !ok {
		return nil, errors.New("error fetching credentials for binding")
	}

	var decodedCreds map[string]interface{}
	if err := json.Unmarshal(encodedCreds, &decodedCreds); err != nil {
		return nil, err
	}

	return decodedCreds, nil
}
