package dispatch

import "hidevops.io/cube/manager/pkg/command"

type Aggregate interface {
	Create(params *command.PipelineReqParams) (err error)
}
