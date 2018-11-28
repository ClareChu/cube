package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CODEPATH   string = "codepath"
	CLONE      string = "clone"
	COMPILE    string = "compile"
	BUILDIMAGE string = "buildImage"
	PUSHIMAGE  string = "pushImage"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type Pipeline struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec PipelineSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status PipelineStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type PipelineSpec struct {
	App       string   `json:"app"  protobuf:"bytes,1,opt,name=app"`
	Profile   string   `json:"profile"  protobuf:"bytes,2,opt,name=profile"`
	Project   string   `json:"project"  protobuf:"bytes,3,opt,name=project"`
	Namespace string   `json:"namespace"  protobuf:"bytes,5,opt,name=namespace"`
	Version   string   `json:"version"  protobuf:"bytes,7,opt,name=version"`
	Events    []Events `json:"events" protobuf:"bytes,16,opt,name=events"`
}

type Events struct {
	Name       string `json:"name" protobuf:"bytes,1,opt,name=name"`
	EventTypes string `json:"eventTypes" protobuf:"bytes,2,opt,name=eventTypes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type PipelineList struct {
	metav1.TypeMeta `json:",inline" protobuf:"bytes,1,opt,name=kind"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,2,opt,name=kind"`
	Items           []Pipeline `json:"items" protobuf:"bytes,3,opt,name=kind"`
}

type PipelineStatus struct {
	Kind           string      `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	Name           string      `json:"name" protobuf:"bytes,2,opt,name=name"`
	Namespace      string      `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	Phase          string      `json:"phase"  protobuf:"bytes,4,opt,name=phase"`
	Stages         []Stages    `json:"stages" protobuf:"bytes,5,opt,name=stages"`
	StartTimestamp metav1.Time `json:"startTimestamp" protobuf:"bytes,6,opt,name=startTimestamp"`
}

type Stages struct {
	Name                 string `json:"name" protobuf:"bytes,1,opt,name=name"`
	StartTime            int64  `json:"startTime" protobuf:"bytes,2,opt,name=startTime"`
	DurationMilliseconds int64  `json:"durationMilliseconds" protobuf:"bytes,3,opt,name=durationMilliseconds"`
}
