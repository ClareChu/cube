package cr

import (
	"gotest.tools/assert"
	"testing"
)

func TestCubeManagerCustomResourceDefinition_Run(t *testing.T) {
	crd, err := fakeInitCube()
	assert.Equal(t, nil, err)
	crd.Run()
	crd.Run()
}
