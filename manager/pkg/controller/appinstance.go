package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type AppInstanceController struct {
	at.JwtRestController
	buildNode              builder.BuildNode
	serviceConfigAggregate aggregate.ServiceConfigAggregate
	volumeAggregate        aggregate.VolumeAggregate
	ingressService         service.IngressService
	podBuilder             builder.PodBuilder
}

func init() {
	app.Register(newAppInstanceController)
}

func newAppInstanceController(buildNode builder.BuildNode, serviceConfigAggregate aggregate.ServiceConfigAggregate,
	ingressService service.IngressService, podBuilder builder.PodBuilder,
	volumeAggregate aggregate.VolumeAggregate) *AppInstanceController {
	return &AppInstanceController{
		buildNode:              buildNode,
		serviceConfigAggregate: serviceConfigAggregate,
		ingressService:         ingressService,
		podBuilder:             podBuilder,
		volumeAggregate:        volumeAggregate,
	}
}

func (p *AppInstanceController) Delete(pod *command.Pod) (rep model.Response, err error) {
	rep = new(model.BaseResponse)
	//TODO delete deployment config
	err = p.buildNode.Delete(pod.Name, pod.Namespace)
	if err != nil {
		return
	}
	//TODO delete service
	err = p.serviceConfigAggregate.Delete(pod.Name, pod.Namespace)

	//TODO delete ingress
	err = p.ingressService.Delete(pod.Name, pod.Namespace)

	//TODO delete pvc and pv
	err = p.volumeAggregate.Delete(pod.Name, pod.Namespace)
	return
}

func (p *AppInstanceController) Get(pod *command.Pod) (rep model.Response, err error) {
	read, err := p.podBuilder.GetPod(pod.Name, pod.Namespace)
	rep = new(model.BaseResponse)
	rep.SetData(read)
	return
}
