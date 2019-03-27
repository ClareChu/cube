package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
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
	NodeSelector     map[string]string      `json:"nodeSelector" protobuf:"bytes,2,opt,name=nodeSelector"`
	Env              []corev1.EnvVar        `json:"env" protobuf:"bytes,3,opt,name=env"`
	Port             []corev1.ContainerPort `json:"port" protobuf:"bytes,4,opt,name=port"`
	Image            string                 `json:"image"  protobuf:"bytes,5,opt,name=image"`
	LivenessProbe    *corev1.Probe          `json:"livenessProbe" protobuf:"bytes,6,opt,name=livenessProbe"`
	ReadinessProbe   *corev1.Probe          `json:"readinessProbe" protobuf:"bytes,7,opt,name=readinessProbe"`
	EnvType          []string               `json:"envType" protobuf:"bytes,8,opt,name=envType"`
	Labels           map[string]string      `json:"labels"  protobuf:"bytes,9,opt,name=labels"`
	DockerRegistry   string                 `json:"dockerRegistry" protobuf:"bytes,9,opt,name=dockerRegistry"`
	Replicas         *int32                 `json:"replicas" protobuf:"bytes,10,opt,name=replicas"`
	Profile          string                 `json:"profile"  protobuf:"bytes,11,opt,name=profile"`
	FromRegistry     string                 `json:"fromRegistry" protobuf:"bytes,12,opt,name=fromRegistry"`
	Tag              string                 `json:"tag" protobuf:"bytes,13,opt,name=tag"`
	Version          string                 `json:"version" protobuf:"bytes,14,opt,name=version"`
	DockerAuthConfig AuthConfig             `json:"dockerAuthConfig" protobuf:"bytes,15,opt,name=dockerAuthConfig"`
	Volumes          []corev1.Volume        `json:"volume" protobuf:"bytes,16,opt,name=volume"`
	VolumeMounts     []corev1.VolumeMount   `json:"volumeMounts" protobuf:"bytes,17,opt,name=volumeMounts"`
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
