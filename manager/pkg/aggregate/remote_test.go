package aggregate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoteTagImage(t *testing.T) {
	tag := GetRemoteTag("a")
	assert.Equal(t, "a", tag)
}
