package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// App is a specification for a App resource
type App struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec AppSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status AppStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type AppSpec struct {
	Name         string           `json:"name, omitempty" protobuf:"bytes,1,opt,name=name"`
	Namespace    string           `json:"namespace, omitempty" protobuf:"bytes,2,opt,name=namespace"`
	TemplateName string           `json:"templateName, omitempty" protobuf:"bytes,3,opt,name=templateName"`
	Version      string           `json:"version, omitempty" protobuf:"bytes,4,opt,name=version"`
	Profile      string           `json:"profile, omitempty" protobuf:"bytes,5,opt,name=profile"`
	Branch       string           `json:"branch, omitempty" protobuf:"bytes,6,opt,name=branch"`
	Context      []string         `json:"context, omitempty" protobuf:"bytes,7,opt,name=context"`
	AppRoot      string           `json:"appRoot, omitempty" protobuf:"bytes,8,opt,name=appRoot"`
	Path         string           `json:"path, omitempty" protobuf:"bytes,9,opt,name=path"`
	Project      string           `json:"project, omitempty" protobuf:"bytes,10,opt,name=project"`
	Url          string           `json:"url, omitempty" protobuf:"bytes,11,opt,name=url"`
	Container    corev1.Container `json:"container" protobuf:"bytes,12,opt,name=container"`
}

type AppStatus struct {
	StartTime string `json:"startTime, omitempty" protobuf:"bytes,1,opt,name=startTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppList is a list of App resources
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []App `json:"items"`
}
