package entity

type Clone struct {
	CloneType string `protobuf:"bytes,1,opt,name=cloneType,proto3" json:"cloneType,omitempty"`
	Url       string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Branch    string `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`
	DstDir    string `protobuf:"bytes,4,opt,name=dstDir,proto3" json:"dstDir,omitempty"`
	Username  string `protobuf:"bytes,5,opt,name=username,proto3" json:"username,omitempty"`
	Password  string `protobuf:"bytes,6,opt,name=password,proto3" json:"password,omitempty"`
	Depth     int32  `protobuf:"varint,7,opt,name=depth,proto3" json:"depth,omitempty"`
	Namespace string `protobuf:"bytes,8,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name      string `protobuf:"bytes,9,opt,name=name,proto3" json:"name,omitempty"`
	Token     string `protobuf:"bytes,10,opt,name=token,proto3" json:"token,omitempty"`
}

type Command struct {
	CodeType    string   `protobuf:"bytes,1,opt,name=codeType,proto3" json:"codeType,omitempty"`
	ExecType    string   `protobuf:"bytes,2,opt,name=execType,proto3" json:"execType,omitempty"`
	Script      string   `protobuf:"bytes,3,opt,name=script,proto3" json:"script,omitempty"`
	CommandName string   `protobuf:"bytes,4,opt,name=commandName,proto3" json:"commandName,omitempty"`
	Params      []string `protobuf:"bytes,5,rep,name=params,proto3" json:"params,omitempty"`
}
