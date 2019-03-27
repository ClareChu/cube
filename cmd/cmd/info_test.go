package cmd

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app/cli"
	"testing"
)

func TestInfoCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, NewRootCommand)

	t.Run("should info success", func(t *testing.T) {
		_, err := testApp.Run("info")
		assert.Equal(t, nil, err)
	})
}
