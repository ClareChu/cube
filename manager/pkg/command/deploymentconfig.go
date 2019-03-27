package command

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentConfig struct {
	model.RequestBody
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec v1alpha1.DeploymentConfigSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status v1alpha1.DeploymentConfigStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type DeployConfigType struct {
	model.RequestBody
	Version      string `json:"version"`
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	SourceType   string `json:"sourceType"`
	PipelineName string `json:"pipelineName"`
	Profile      string `json:"profile"`
}
