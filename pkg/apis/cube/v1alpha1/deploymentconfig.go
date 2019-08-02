package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type DeploymentConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec DeploymentConfigSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status DeploymentConfigStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type DeploymentConfigSpec struct {
	NodeSelector     map[string]string         `json:"nodeSelector" protobuf:"bytes,1,opt,name=nodeSelector"`
	Image            string                    `json:"image"  protobuf:"bytes,2,opt,name=image"`
	EnvType          []string                  `json:"envType" protobuf:"bytes,3,opt,name=envType"`
	Labels           map[string]string         `json:"labels"  protobuf:"bytes,4,opt,name=labels"`
	DockerRegistry   string                    `json:"dockerRegistry" protobuf:"bytes,5,opt,name=dockerRegistry"`
	Replicas         *int32                    `json:"replicas" protobuf:"bytes,6,opt,name=replicas"`
	Profile          string                    `json:"profile"  protobuf:"bytes,7,opt,name=profile"`
	FromRegistry     string                    `json:"fromRegistry" protobuf:"bytes,8,opt,name=fromRegistry"`
	Tag              string                    `json:"tag" protobuf:"bytes,9,opt,name=tag"`
	Version          string                    `json:"version" protobuf:"bytes,10,opt,name=version"`
	DockerAuthConfig AuthConfig                `json:"dockerAuthConfig" protobuf:"bytes,11,opt,name=dockerAuthConfig"`
	Volumes          []corev1.Volume           `json:"volume" protobuf:"bytes,12,opt,name=volume"`
	Container        corev1.Container          `json:"container" protobuf:"bytes,13,opt,name=container"`
	Strategy         appsv1.DeploymentStrategy `json:"strategy,omitempty" protobuf:"bytes,14,opt,name=strategy"`
	InitContainer    corev1.Container          `json:"initContainer" protobuf:"bytes,15,opt,name=initContainer"`
	ForceUpdate      bool                      `json:"forceUpdate" protobuf:"bytes,16,opt,name=forceUpdate"`
}

type DeploymentConfigStatus struct {
	LastVersion int `json:"lastVersion,omitempty" protobuf:"bytes,1,opt,name=lastVersion"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type DeploymentConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DeploymentConfig `json:"items"`
}
