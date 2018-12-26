package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ImageStream struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec ImageStreamSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status metav1.Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ImageStreamSpec struct {
	DockerImageRepository string `json:"dockerImageRepository,omitempty" protobuf:"bytes,1,opt,name=dockerImageRepository"`
	Tags                  []Tag  `json:"tags,omitempty" protobuf:"bytes,2,opt,name=tags`
}

type Tag struct {
	Created              string `json:"created,omitempty" protobuf:"bytes,1,opt,name=created"`
	DockerImageReference string `json:"dockerImageReference,omitempty" protobuf:"bytes,1,opt,name=dockerImageReference"`
	Generation           string `json:"generation" protobuf:"bytes,2,opt,name=generation"`
	Version              string `json:"version" protobuf:"bytes,3,opt,name=version"`
}

type Status struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ImageStreamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ImageStream `json:"items"`
}
