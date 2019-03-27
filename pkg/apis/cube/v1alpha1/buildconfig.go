package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type BuildConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec BuildSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status BuildConfigStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type BuildConfigStatus struct {
	LastVersion int `json:"lastVersion,omitempty" protobuf:"bytes,1,opt,name=lastVersion"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type BuildConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []BuildConfig `json:"items"`
}

type TypeStatus struct {
	EventType string `json:"eventType,omitempty" protobuf:"bytes,2,opt,name=eventType"`
	// Status of the operation.
	// One of: "Success" or "Failure".
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Status string `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
	// A human-readable description of the status of this operation.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// A machine-readable description of why this operation is in the
	// "Failure" status. If this value is empty there
	// is no information available. A Reason clarifies an HTTP status
	// code but does not override it.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason,casttype=StatusReason"`
	// Extended data associated with the reason.  Each reason may define its
	// own extended details. This field is optional and the data returned
	// is not guaranteed to conform to any schema except that defined by
	// the reason type.
	// +optional
	Details *StatusDetails `json:"details,omitempty" protobuf:"bytes,5,opt,name=details"`
	// Suggested HTTP return code for this status, 0 if not set.
	// +optional
	Code int32 `json:"code,omitempty" protobuf:"varint,6,opt,name=code"`
}

type StatusDetails struct {
	Name                 string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	DurationMilliseconds int64  `json:"durationMilliseconds,omitempty" protobuf:"bytes,2,opt,name=durationMilliseconds"`
	StartTime            int64  `json:"startTime,omitempty" protobuf:"bytes,3,opt,name=startTime"`
}
