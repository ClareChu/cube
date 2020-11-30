package crd

import (
	"gotest.tools/assert"
	"testing"
)

func TestInitCRD(t *testing.T) {
	crd, err := fakeInitCRD()
	assert.Equal(t, nil, err)
	crd.Run()
	crd.Run()
}
