package v2

import (
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
)

type DeleteNodeInterface interface {
	Handle(build *v1alpha1.Build) error
}

type DeleteNode struct {
	DeleteNodeInterface
	buildNode          builder.BuildNode
	build              BuildInterface
	serviceInterface   ServiceInterface
	buildPackInterface dispatch.BuildPackInterface
}

func init() {
	app.Register(NewDeleteNodeService)
}

const DeleteDeployment = "deleteDeployment"

func NewDeleteNodeService(buildNode builder.BuildNode,
	build BuildInterface,
	serviceInterface ServiceInterface,
	buildPackInterface dispatch.BuildPackInterface) DeleteNodeInterface {
	deleteNode := &DeleteNode{
		buildNode:          buildNode,
		build:              build,
		serviceInterface:   serviceInterface,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(DeleteDeployment, deleteNode)
	return deleteNode
}

func (d *DeleteNode) Handle(build *v1alpha1.Build) error {
	phase := constant.Success
	//TODO delete deployment config
	err := d.buildNode.DeleteDeployment(build.ObjectMeta.Labels["name"], build.Namespace)
	//TODO delete service
	err = d.serviceInterface.DeleteService(build.ObjectMeta.Labels["name"], build.Namespace)
	if err != nil {
		phase = constant.Fail
	}
	err = d.build.Update(build, constant.DeleteDeployment, phase)
	return err
}
