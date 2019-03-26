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
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/cli"
)

const appVersion = "v1.0.1"

// TODO: should get api version as well
type versionCommand struct {
	cli.SubCommand
}

func newVersionCommand() *versionCommand {
	c := &versionCommand{}
	c.Use = "version"
	c.Short = "Fetch the version of a cli and api"
	c.Long = "Run version command"
	return c
}

func init() {
	app.Register(newVersionCommand)
}

//cube/client log --profile=test --project=demo --app=hello-world --sourcecode=java --verbose=true
func (c *versionCommand) Run(args []string) error {

	fmt.Printf("%v\n", appVersion)

	return nil
}
