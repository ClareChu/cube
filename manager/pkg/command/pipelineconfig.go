package command

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineStart struct {
	model.RequestBody	  `json:"omitempty"`
	Name         string   `json:"name"`
	Namespace    string   `json:"namespace"`
	SourceCode   string   `json:"sourceCode"`
	Version      string   `json:"version"`
	Profile      string   `json:"profile"`
	Branch       string   `json:"branch"`
	Context      []string `json:"context"`
	ParentModule string   `json:"parentModule"`
	Path         string   `json:"path"`
	Project      string   `json:"project"`
}

type PipelineConfigTemplate struct {
	model.RequestBody
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec v1alpha1.PipelineSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status v1alpha1.PipelineConfigStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type StartPipeline struct {
	model.RequestBody
	SourceCode string
	App        string
	Namespace  string
}

type PipelineReqParams struct {
	Name         string `json:"name"`
	PipelineName string `json:"pipeline_name"`
	Namespace    string `json:"namespace"`
	EventType    string `json:"event_type"`
	Version      string `json:"version"`
	Branch       string `json:"branch"`
	Context      string `json:"context"`
	ParentModule string `json:"parent_module"`
	Profile      string `json:"profile"`
	Project      string `json:"project"`
}
