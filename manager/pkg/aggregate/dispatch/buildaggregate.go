package dispatch

import "hidevops.io/cube/pkg/apis/cube/v1alpha1"

type BuildPackAggregate interface {
	Handle(build *v1alpha1.Build) (err error)
}
