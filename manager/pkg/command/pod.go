package command

import "hidevops.io/hiboot/pkg/model"

type Pod struct {
	model.RequestBody
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
