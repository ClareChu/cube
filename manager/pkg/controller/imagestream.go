package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type ImageStreamController struct {
	at.RestController
	imageStreamAggregate aggregate.ImageStreamAggregate
}

func init() {
	app.Register(newImageStreamControllerController)
}

func newImageStreamControllerController(imageStreamAggregate aggregate.ImageStreamAggregate) *ImageStreamController {
	return &ImageStreamController{
		imageStreamAggregate: imageStreamAggregate,
	}
}

type ImageStream struct {
	model.RequestBody
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Images    []string `json:"images"`
}

func (i *ImageStreamController) Post(is *ImageStream) (rep model.Response, err error) {
	err = i.imageStreamAggregate.Create(is.Name, is.Namespace, is.Images)
	rep = new(model.BaseResponse)
	return
}
