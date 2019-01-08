package command

import (
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildConfig struct {
	model.RequestBody
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec v1alpha1.BuildSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status BuildConfigStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type BuildConfigStatus struct {
	LastVersion int `json:"lastVersion,omitempty" protobuf:"bytes,1,opt,name=lastVersion"`
}

type SourceCodePullCommand struct {
	model.RequestBody
	Host      string `json:"host"`
	Port      string `json:"port"`
	CloneType string `protobuf:"bytes,1,opt,name=cloneType,proto3" json:"cloneType,omitempty"`
	Url       string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Branch    string `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`
	DstDir    string `protobuf:"bytes,4,opt,name=dstDir,proto3" json:"dstDir,omitempty"`
	Username  string `protobuf:"bytes,5,opt,name=username,proto3" json:"username,omitempty"`
	Password  string `protobuf:"bytes,6,opt,name=password,proto3" json:"password,omitempty"`
	Depth     int32  `protobuf:"varint,7,opt,name=depth,proto3" json:"depth,omitempty"`
	Namespace string `protobuf:"bytes,8,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name      string `protobuf:"bytes,9,opt,name=name,proto3" json:"name,omitempty"`
}

type BuildCommand struct {
	CodeType    string   `protobuf:"bytes,1,opt,name=codeType,proto3" json:"codeType,omitempty"`
	ExecType    string   `protobuf:"bytes,2,opt,name=execType,proto3" json:"execType,omitempty"`
	Script      string   `protobuf:"bytes,3,opt,name=Script,proto3" json:"Script,omitempty"`
	CommandName string   `protobuf:"bytes,4,opt,name=commandName,proto3" json:"commandName,omitempty"`
	Params      []string `protobuf:"bytes,5,rep,name=Params,proto3" json:"Params,omitempty"`
}

type CompileCommand struct {
	model.RequestBody
	CompileCmd []*BuildCommand `protobuf:"bytes,1,rep,name=CompileCmd,proto3" json:"CompileCmd,omitempty"`
	Namespace  string          `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name       string          `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

type ImageBuildCommand struct {
	model.RequestBody
	App        string   `protobuf:"bytes,1,opt,name=app,proto3" json:"app,omitempty"`
	S2IImage   string   `protobuf:"bytes,2,opt,name=s2iImage,proto3" json:"s2iImage,omitempty"`
	Tags       []string `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
	DockerFile []string `protobuf:"bytes,4,rep,name=dockerFile,proto3" json:"dockerFile,omitempty"`
	Namespace  string   `protobuf:"bytes,5,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name       string   `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Username   string   `protobuf:"bytes,7,opt,name=username,proto3" json:"username,omitempty"`
	Password   string   `protobuf:"bytes,8,opt,name=password,proto3" json:"password,omitempty"`
}

type ImagePushCommand struct {
	model.RequestBody
	Tags      []string `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty"`
	Namespace string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name      string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Username  string   `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Password  string   `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	ImageName string   `protobuf:"bytes,6,opt,name=imageName,proto3" json:"imageName,omitempty"`
}

type BuildSpec struct {
	CloneConfig SourceCodePullCommand `json:"cloneConfig"  protobuf:"bytes,1,opt,name=cloneConfig"`
	App         string                `json:"app"  protobuf:"bytes,1,opt,name=app"`
	CodeType    string                `json:"codeType"  protobuf:"bytes,1,opt,name=codeType"`
	CompileCmd  []CompileCommand      `json:"compileCmd"  protobuf:"bytes,1,opt,name=compileCmd"`
	CloneType   string                `json:"cloneType"  protobuf:"bytes,1,opt,name=cloneType"`
	S2iImage    string                `json:"s2iImage"  protobuf:"bytes,1,opt,name=s2iImage"`
	Tags        []string              `json:"tags"  protobuf:"bytes,1,opt,name=tags"`
	DockerFile  []string              `json:"dockerFile"  protobuf:"bytes,1,opt,name=dockerFile"`
}

type DeployNode struct {
	model.RequestBody
	kube.DeployData
}

type ServiceNode struct {
	model.RequestBody
	kube.DeployData
}

type BuildConfigTemplate struct {
	model.RequestBody
	Name         string   `json:"name"`
	Namespace    string   `json:"namespace"`
	SourceType   string   `json:"sourceType"`
	Events       []string `json:"events"`
	Version      string   `json:"version"`
	PipelineName string   `json:"pipelineName"`
}
