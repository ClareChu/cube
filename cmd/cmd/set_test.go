package cmd

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app/cli"
	"testing"
)

func TestSetCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, NewRootCommand)

	t.Run("should set success", func(t *testing.T) {
		_, err := testApp.Run("set", "--server=xxx")
		assert.NotEqual(t, nil, err)
	})

}
