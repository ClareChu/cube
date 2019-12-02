package dispatch

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/manager/pkg/aggregate"
	"testing"
)

func TestNewPipelineFactory(t *testing.T) {
	pf := NewPipelineFactory()
	cb := &aggregate.Callback{}
	pf.Add("v", cb)
	a := pf.Get("v")
	err := a.Create(nil)
	assert.Equal(t, nil, err)
}