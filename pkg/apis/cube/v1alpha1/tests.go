package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type Tests struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Tests.
	// +optional
	Spec TestConfigSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Tests.
	// +optional
	Status TestsStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type TestsStatus struct {
	Kind           string      `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	Name           string      `json:"name" protobuf:"bytes,2,opt,name=name"`
	Namespace      string      `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	Phase          string      `json:"phase"  protobuf:"bytes,4,opt,name=phase"`
	Stages         []Stages    `json:"stages" protobuf:"bytes,5,opt,name=stages"`
	StartTimestamp metav1.Time `json:"startTimestamp" protobuf:"bytes,6,opt,name=startTimestamp"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type TestsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Tests `json:"items"`
}
