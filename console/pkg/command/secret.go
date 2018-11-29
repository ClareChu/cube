package command

import "hidevops.io/hiboot/pkg/model"

type Secret struct {
	model.RequestBody
	Username  string `json:"username"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}
