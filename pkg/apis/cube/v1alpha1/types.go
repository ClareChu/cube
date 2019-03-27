package v1alpha1

type ExecCommand struct {
	ExecType      string   `protobuf:"bytes,1,opt,name=execType,proto3" json:"execType,omitempty"`
	Script        string   `protobuf:"bytes,2,opt,name=Script,proto3" json:"Script,omitempty"`
	CommandName   string   `protobuf:"bytes,3,opt,name=commandName,proto3" json:"commandName,omitempty"`
	CommandParams []string `protobuf:"bytes,4,rep,name=commandParams,proto3" json:"commandParams,omitempty"`
}
