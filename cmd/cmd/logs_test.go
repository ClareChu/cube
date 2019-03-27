package cmd

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app/cli"
	"testing"
)

func TestLogCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, NewRootCommand)

	t.Run("should logs success", func(t *testing.T) {
		_, err := testApp.Run("logs")
		assert.NotEqual(t, nil, err)
	})
}
