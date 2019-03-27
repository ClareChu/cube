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
	"encoding/json"
	"fmt"
	"hidevops.io/cube/cmd/api"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/cli"
	"hidevops.io/hiboot/pkg/log"
)

type infoCommand struct {
	cli.SubCommand

	info string
}

func newInfoCommand() *infoCommand {
	c := &infoCommand{}
	c.Use = "info"
	c.Short = "Display client configuration information"
	c.Long = "Run info command"

	pf := c.PersistentFlags()
	pf.StringVarP(&c.info, "info", "i", "", "")
	return c
}

func init() {
	app.Register(newInfoCommand)
}

//cube info
func (c *infoCommand) Run(args []string) error {
	log.Debug("handle info command")

	user, filePath, err := api.ReadConfig()
	if err != nil {
		return err
	}

	userByte, err := json.Marshal(user)
	if err != nil {
		return err
	}
	fmt.Println(filePath, ": ", string(userByte))
	return nil
}
