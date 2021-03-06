// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/cmd/api"
	"hidevops.io/hiboot/pkg/app/cli"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestRunCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, NewRootCommand)

	t.Run("should run", func(t *testing.T) {
		_, err := testApp.Run("run", "-h")
		assert.Equal(t, nil, err)
	})

	t.Run("should run", func(t *testing.T) {
		_, err := testApp.Run("run", "-s", "go", "-c", "abc/def,abc/def/hij")
		assert.Equal(t, nil, err)
	})

}

func TestGetProject(t *testing.T) {
	a, b, err := api.GetProjectInfoByCurrPath()
	assert.Equal(t, nil, err)
	log.Infof("a : %v", a)
	log.Infof("b : %v", b)
}
