package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SourceConfig is a specification for a SourceConfig resource
type SourceConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec SourceConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status metav1.Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type SourceConfigSpec struct {
	Scm        SourceConfigScm `json:"scm"`
	SourceCode []SourceCode    `json:"sourceCode"`
}

type SourceConfigScm struct {
	Url        string `json:"url"`
	Ref        string `json:"ref"`
	Token      string `json:"token"`
	UserName   string `json:"userName"`
	ApiVersion string `json:"apiVersion"`
	ProjectId  string `json:"projectId"`
}

type SourceCode struct {
	Type    string `json:"type"`    //java cqrs
	Path    string `json:"path"`    //路径
	Content string `json:"content"` //内容
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SourceConfigList is a list of SourceConfig resources
type SourceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SourceConfig `json:"items"`
}
