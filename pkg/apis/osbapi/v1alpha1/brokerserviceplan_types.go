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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BrokerServicePlanSpec defines the desired state of BrokerServicePlan
type BrokerServicePlanSpec struct {
	Name      string `json:"name"`
	ServiceID string `json:"serviceID"`
}

// BrokerServicePlanStatus defines the observed state of BrokerServicePlan
type BrokerServicePlanStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerServicePlan is the Schema for the brokerserviceplans API
// +k8s:openapi-gen=true
type BrokerServicePlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerServicePlanSpec   `json:"spec,omitempty"`
	Status BrokerServicePlanStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerServicePlanList contains a list of BrokerServicePlan
type BrokerServicePlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrokerServicePlan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrokerServicePlan{}, &BrokerServicePlanList{})
}
