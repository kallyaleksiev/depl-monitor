/*
Copyright 2023 Kaloyan Aleksiev.

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

package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Type which represents the configs which are going to get passed
// as labels to the `underlying`
type MonConfigs map[string]string

// MonDeplSpec defines the desired state of MonDepl
type MonDeplSpec struct {
	//+kubebuilder:validation:MinLength=0

	// Reason for the creation of this monitoring
	Reason string `json:"reason"`

	// The configs, which will be injected as labels
	// for underlying
	//+optional
	Configs MonConfigs `json:"configs,omitempty"`

	// Underlying deployment
	Underlying appsv1.DeploymentSpec `json:"underlying"`
	// Underlying batchv1.JobTemplateSpec `json:"underlying"`
}

// MonDeplStatus defines the observed state of MonDepl
type MonDeplStatus struct {
	// A pointer to currently running underlying.
	// +optional
	Active *corev1.ObjectReference `json:"active,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MonDepl is the Schema for the mondepls API
type MonDepl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonDeplSpec   `json:"spec,omitempty"`
	Status MonDeplStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MonDeplList contains a list of MonDepl
type MonDeplList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MonDepl `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MonDepl{}, &MonDeplList{})
}
