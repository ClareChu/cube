package aggregate

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetNamespace(t *testing.T) {
	space := "demo"
	profile := "dev"
	n := GetNamespace(space, profile)
	assert.Equal(t, "demo-dev", n)
	profile = ""
	n = GetNamespace(space, profile)
	assert.Equal(t, "demo", n)
}

func TestTagTagImage(t *testing.T) {
}