package command

import "hidevops.io/hiboot/pkg/model"

type Remote struct {
	model.RequestBody
	Profile        string `json:"profile"`
	DockerRegistry string `json:"dockerRegistry"`
	App            string `json:"app"`
	Namespace      string `json:"namespace"`
	Tag            string `json:"tag"`
}
