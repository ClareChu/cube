package cmd

import (
	"gotest.tools/assert"
	"hidevops.io/hiboot/pkg/app/cli"
	"testing"
)

func TestNewAppCommand(t *testing.T) {
	testApp := cli.NewTestApplication(t, newAppCommand())
	t.Run("should create app", func(t *testing.T) {
		_, err := testApp.Run("create")
		assert.Equal(t, nil, err)
	})

}