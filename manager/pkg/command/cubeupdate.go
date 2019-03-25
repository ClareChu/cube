package command

import "hidevops.io/hiboot/pkg/model"

type CubeUpdate struct {
	model.RequestBody
	Version string `json:"version"` //版本号
	Url     string `json:"url"`     // 下载链接
	Type    string `json:"type"`    // 操作系统
	Arch    string `json:"arch"`    //系统位数
	Enable  bool   `json:"enable"`  //是否更新
}
