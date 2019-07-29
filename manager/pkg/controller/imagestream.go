package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type ImageStreamController struct {
	at.RestController
	imageStreamAggregate aggregate.ImageStreamAggregate
	volumeAggregate      aggregate.VolumeAggregate
	callbackAggregate    aggregate.CallbackAggregate
}

func init() {
	app.Register(newImageStreamControllerController)
}

func newImageStreamControllerController(imageStreamAggregate aggregate.ImageStreamAggregate,
	volumeAggregate aggregate.VolumeAggregate,
	callbackAggregate aggregate.CallbackAggregate) *ImageStreamController {
	return &ImageStreamController{
		imageStreamAggregate: imageStreamAggregate,
		volumeAggregate:      volumeAggregate,
		callbackAggregate:    callbackAggregate,
	}
}

type ImageStream struct {
	model.RequestBody
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Images    []string `json:"images"`
}

type Volumes struct {
	model.RequestBody
	Volume    v1alpha1.Volumes `json:"volume"`
	Namespace string           `json:"namespace"`
}

func (i *ImageStreamController) Post(is *ImageStream) (rep model.Response, err error) {
	err = i.imageStreamAggregate.Create(is.Name, is.Namespace, is.Images)
	rep = new(model.BaseResponse)
	return
}

func (i *ImageStreamController) PostVolume(v *Volumes) (rep model.Response, err error) {
	err = i.volumeAggregate.CreateVolume(v.Namespace, v.Volume)
	rep = new(model.BaseResponse)
	return
}

type Watch struct {
	model.RequestBody
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (i *ImageStreamController) PostWatch(w *Watch) (rep model.Response, err error) {
	err = i.callbackAggregate.WatchPod(w.Name, w.Namespace)
	rep = new(model.BaseResponse)
	return
}
