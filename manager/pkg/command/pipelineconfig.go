package command

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineStart struct {
	model.RequestBody `json:"omitempty"`
	Id                int      `json:"id"`
	ForceUpdate       bool     `json:"forceUpdate" default:"true"`
	Name              string   `json:"name"`
	Namespace         string   `json:"namespace"`
	TemplateName      string   `json:"templateName"`
	Version           string   `json:"version" default:"v1"`
	Profile           string   `json:"profile"`
	Branch            string   `json:"branch"`
	Context           []string `json:"context"`
	AppRoot           string   `json:"appRoot" validate:"required"`
	//获取path 的目录
	Path    string `json:"path"`
	Project string `json:"project" validate:"required"`
	//gitUrl
	Url           string             `json:"url"`
	Ingress       []v1alpha1.Ingress `json:"ingress"`
	Container     corev1.Container   `json:"container"`
	InitContainer corev1.Container   `json:"initContainer"`
	Volumes       v1alpha1.Volumes   `json:"volumes"`
	Callback      string             `json:"callback"`
	IsApp         bool               `json:"isApp"`
	Token         string             `json:"token"`
	Services      []v1alpha1.Service `json:"services"`
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
	Id            int                `json:"id"`
	ForceUpdate   bool               `json:"forceUpdate"`
	Name          string             `json:"name"`
	PipelineName  string             `json:"pipeline_name"`
	Namespace     string             `json:"namespace"`
	EventType     string             `json:"event_type"`
	Version       string             `json:"version"`
	Branch        string             `json:"branch"`
	Context       string             `json:"context"`
	AppRoot       string             `json:"app_root"`
	Profile       string             `json:"profile"`
	Project       string             `json:"project"`
	BuildVersion  string             `json:"buildVersion"`
	Container     corev1.Container   `json:"container"`
	InitContainer corev1.Container   `json:"initContainer"`
	Volumes       v1alpha1.Volumes   `json:"volumes"`
	Ingress       []v1alpha1.Ingress `json:"ingress"`
	Callback      string             `json:"callback"`
	Token         string             `json:"token"`
	Services      []v1alpha1.Service `json:"services"`
}
