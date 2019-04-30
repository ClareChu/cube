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
)

type loginCommand struct {
	cli.SubCommand

	username string
	password string
}

func newLoginCommand() *loginCommand {
	c := &loginCommand{}
	c.Use = "login"
	c.Short = "Log in to a cube/client"
	c.Long = "Run login command"
	pf := c.PersistentFlags()
	pf.StringVarP(&c.username, "username", "u", "", "--username=your-username")
	pf.StringVarP(&c.password, "password", "p", "", "--password=your-password")
	return c
}

func init() {
	app.Register(newLoginCommand)
}

func (c *loginCommand) Run(args []string) error {
	log.Debug("handle login command")
	if c.username == "" {
		c.username = api.GetInput(api.USERNAME)
	}
	if c.password == "" {
		c.password = api.GetInput(api.PASSWORD)
	}

	user, filePath, err := api.ReadConfig()
	if err != nil {
		log.Debug("Error", err)
		return err
	}

	if token, err := api.Login(api.GetLoginApi(user.Server), c.username, c.password); err != nil {
		return err
	} else {
		user.Token = token
		if err := api.WriteConfig(user, filePath); err != nil {
			log.Debug("Error", err)
			return err
		}
	}

	fmt.Println("[INFO] Login successful.")
	return nil
}
