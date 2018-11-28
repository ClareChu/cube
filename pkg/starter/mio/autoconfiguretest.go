// +build !test

package mio

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"testing"
)

func TestMioConfig(t *testing.T) {
	testApp := web.NewTestApplication(t).(app.ApplicationContext)

	t.Run("should get instance", func(t *testing.T) {
		mioBuildConfig := testApp.GetInstance("mioBuildConfig")
		assert.NotEqual(t, nil, mioBuildConfig)
	})
}
