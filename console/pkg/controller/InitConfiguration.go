package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/mio/console/pkg/aggregate"
)

type InitConfiguationController struct {
	at.RestController
	buildAggregate aggregate.BuildAggregate
}

func init() {
	app.Register(newBuildController)
}

func newBuildController(buildAggregate aggregate.BuildAggregate) *BuildController {
	return &BuildController{
		buildAggregate: buildAggregate,
	}
}

