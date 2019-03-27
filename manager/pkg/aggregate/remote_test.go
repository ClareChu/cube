package aggregate

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRemoteTagImage(t *testing.T) {
	tag := GetRemoteTag("a")
	assert.Equal(t, "a", tag)
}
