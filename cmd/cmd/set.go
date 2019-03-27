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
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/cli"
	"hidevops.io/hiboot/pkg/log"
	"os"
	"regexp"
)

type setCommand struct {
	cli.SubCommand

	name   string
	token  string
	server string
}

func newSetCommand() *setCommand {
	c := &setCommand{}
	c.Use = "set"
	c.Short = "set client configuration information"
	c.Long = "Run set command"

	pf := c.PersistentFlags()
	pf.StringVarP(&c.name, "name", "n", "", "--name=test")
	pf.StringVarP(&c.token, "token", "t", "", "--token=XXXXXXXXXXXXXXX")
	pf.StringVarP(&c.server, "server", "s", "", "--server=http://100.100.100.100:8080")
	return c
}

func init() {
	app.Register(newSetCommand)
}

//cube/client set --name=test --token=XXXX --hiadmin=http://100.100.100.100:8080||default
func (c *setCommand) Run(args []string) error {
	log.Debug("handle set command")

	user, filePath, err := api.ReadConfig()
	if err != nil {
		return err
	}
	if c.name == "" && c.token == "" && c.server == "" {
		fmt.Println("No change.")
		os.Exit(0)
	}

	if c.name != "" {
		user.Name = c.name
	}

	if c.token != "" {
		user.Token = c.token
	}

	if c.server != "" {
		if c.server == "default" {
			user.Server = api.DEFAULT_SERVER
		} else {
			r, _ := regexp.Compile("^http://|https://.")
			if !r.MatchString(c.server) {
				fmt.Println("server format error,the correct one should be http://XXXXX or https://XXXXX")
				os.Exit(0)
			}
			user.Server = c.server
		}
	}

	if err := api.WriteConfig(user, filePath); err != nil {
		return err
	}
	fmt.Println("[INFO] Set successful.")
	return nil
}
