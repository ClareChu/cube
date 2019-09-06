package cmd

import (
	"gotest.tools/assert"
	"hidevops.io/hiboot/pkg/app/cli"
	"testing"
)

func TestGetAppCommand(t *testing.T) {
	testApp := cli.NewTestApplication(t, newGetappCommand())
	t.Run("should get app", func(t *testing.T) {
		_, err := testApp.Run("get")
		assert.Equal(t, nil, err)
	})

}