/*

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
	"os"

	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeFeatureDiscoverySpec defines the desired state of NodeFeatureDiscovery
// +k8s:openapi-gen=true
type NodeFeatureDiscoverySpec struct {
	// +optional
	Operand      OperandSpec `json:"operand"`
	WorkerConfig *ConfigMap  `json:"workerConfig,omitempty"`

	// Run NFD in multiple deployment mode
	// https://kubernetes-sigs.github.io/node-feature-discovery/v0.8/advanced/master-commandline-reference.html#-instance
	//
	// +nullable
	// +optional
	Instance string `json:"instance"`

	// ExtraLabelNs defines the list of of allowed extra label namespaces
	// By default, only allow labels in the default `feature.node.kubernetes.io` label namespace
	// +nullable
	// +kubebuilder:validation:Optional
	ExtraLabelNs []string `json:"extraLabelNs,omitempty"`

	// ResourceLabels defines the list of features
	// to be advertised as extended resources instead of labels.
	// +nullable
	// +kubebuilder:validation:Optional
	ResourceLabels []string `json:"resourceLabels,omitempty"`

	// LabelWhiteList defines a regular expression
	// for filtering feature labels based on their name.
	// Each label must match against the given reqular expression in order to be published.
	// +nullable
	// +kubebuilder:validation:Optional
	LabelWhiteList string `json:"labelWhiteList,omitempty"`

	// WorkerConfig describes configuration options for the NFD
	// worker.
	// +optional
	CustomConfig ConfigMap `json:"customConfig"`
}

// OperandSpec describes configuration options for the operand
type OperandSpec struct {
	// Namespace defines the namespace to deploy nfd-master
	// and nfd-worker pods
	// [defaults to openshift-nfd]
	//
	// +nullable
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Image defines the image to pull for the
	// NFD operand
	//
	// +kubebuilder:validation:Pattern=[a-zA-Z0-9\-]+
	// +nullable
	// +optional
	Image string `json:"image,omitempty"`

	// ImagePullPolicy defines Image pull policy for the
	// NFD operand image [defaults to Always]
	//
	// +nullable
	// +optional
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// ServicePort specifies the TCP port that nfd-master
	// listens for incoming requests.
	// +kubebuilder:validation:Optional
	ServicePort int `json:"servicePort"`
}

// ConfigMap describes configuration options for the NFD worker
type ConfigMap struct {
	// BinaryData holds the NFD configuration file
	ConfigData string `json:"configData"`
}

// NodeFeatureDiscoveryStatus defines the observed state of NodeFeatureDiscovery
// +k8s:openapi-gen=true
type NodeFeatureDiscoveryStatus struct {
	// Conditions represents the latest available observations of current state.
	//
	// +optional
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nodefeaturediscoveries,scope=Namespaced
//
// An Node Feature Discovery cluster instance
// +operator-sdk:csv:customresourcedefinitions:displayName="NodeFeatureDiscovery"
type NodeFeatureDiscovery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Node Feature Discovery
	Spec   NodeFeatureDiscoverySpec   `json:"spec,omitempty"`
	Status NodeFeatureDiscoveryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeFeatureDiscoveryList contains a list of NodeFeatureDiscovery
type NodeFeatureDiscoveryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeFeatureDiscovery `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeFeatureDiscovery{}, &NodeFeatureDiscoveryList{})
}

// ImagePath returns a compiled full valid image string
func (o *OperandSpec) ImagePath() string {
	if o.Image != "" {
		return o.Image
	}

	image := os.Getenv("NODE_FEATURE_DISCOVERY_IMAGE")
	return image
}

// ImagePolicy returns a valid corev1.PullPolicy from the string in the CR
func (o *OperandSpec) ImagePolicy(pullPolicy string) corev1.PullPolicy {
	switch corev1.PullPolicy(pullPolicy) {
	case corev1.PullAlways:
		return corev1.PullAlways
	case corev1.PullNever:
		return corev1.PullNever
	}
	return corev1.PullIfNotPresent
}

// Data returns a valid ConfigMap name
func (c *ConfigMap) Data() string {
	return c.ConfigData
}
