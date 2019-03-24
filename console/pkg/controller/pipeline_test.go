package controller

import (
	"hidevops.io/cube/console/pkg/aggregate/mocks"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	"net/http"
	"testing"
)

func TestPipeline(t *testing.T) {
	pipelineAggregate := new(mocks.PipelineAggregate)
	gateway := newPipelineController(pipelineAggregate)
	pic := &v1alpha1.PipelineConfig{}
	pi := &v1alpha1.Pipeline{}
	pipelineAggregate.On("Create", pic, "").Return(pi, nil)
	app := web.NewTestApp(t, gateway).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/pipeline").WithJSON(&command.PipelineConfigTemplate{}).Expect().Status(http.StatusOK)
	})

}
