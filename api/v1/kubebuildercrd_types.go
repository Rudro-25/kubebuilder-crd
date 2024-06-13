/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KubebuilderCrd is the Schema for the kubebuildercrds API
type KubebuilderCrd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubebuilderCrdSpec   `json:"spec,omitempty"`
	Status KubebuilderCrdStatus `json:"status,omitempty"`
}

// KubebuilderCrdSpec defines the desired state of KubebuilderCrd
type KubebuilderCrdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	DeploymentName string        `json:"deploymentName,omitempty"`
	Replicas       *int32        `json:"replicas,omitempty"`
	Container      ContainerSpec `json:"container"`

	// +optional
	Service ServiceSpec `json:"service,omitempty"`
}

type ContainerSpec struct {
	Image string `json:"image"`
	Port  int32  `json:"port"`
}

type ServiceSpec struct {
	// +optional
	ServiceName string `json:"serviceName"`
	ServiceType string `json:"serviceType"`
	// +optional
	ServiceNodePort int32 `json:"serviceNodePort"`
}

// KubebuilderCrdStatus defines the observed state of KubebuilderCrd
type KubebuilderCrdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	AvailableReplicas *int32 `json:"availableReplicas"`
}

// +kubebuilder:object:root=true

// KubebuilderCrdList contains a list of KubebuilderCrd
type KubebuilderCrdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubebuilderCrd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubebuilderCrd{}, &KubebuilderCrdList{})
}
