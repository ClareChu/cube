package service

import (
	"gotest.tools/assert"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"testing"
)

func TestAppCopy(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	app := cube.NewApp(clientSet)
	as := newAppCommand(app)
	cmd := &command.PipelineStart{

	}
	_, err := as.Init(cmd)
	assert.Assert(t, nil, err)
}
