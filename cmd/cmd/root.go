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
	"hidevops.io/cube/cmd/api"
	"hidevops.io/hiboot/pkg/app/cli"
)

// RootCommand is the root command
type RootCommand struct {
	cli.RootCommand
}

func init() {
	user, _, err := api.ReadConfig()
	if err != nil {
		fmt.Println("Configuration information acquisition failed")
		return
	}
	api.DoUpdate(api.GetCliUpdateApi(user.Server))
}

// NewRootCommand the root command
func NewRootCommand(run *runCommand,
	login *loginCommand,
	logs *logsCommand,
	set *setCommand,
	get *getCommand,
	info *infoCommand,
	app *appCommand,
	version *versionCommand,
) *RootCommand {
	c := new(RootCommand)
	c.Use = "cube"
	c.Short = "cube command"
	c.Long = "Run cube command"
	c.ValidArgs = []string{"baz"}
	c.Add(run, login, logs, set, info, version, app, get)
	return c
}

// Run root command handler
func (c *RootCommand) Run(args []string) error {
	fmt.Printf("cube %v\n\n", appVersion)
	c.Help()
	return nil
}
