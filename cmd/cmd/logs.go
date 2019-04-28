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
)

type logsCommand struct {
	cli.SubCommand

	profile string
	project string
	app     string
	verbose bool
}

func newLogsCommand() *logsCommand {
	c := &logsCommand{}
	c.Use = "logs"
	c.Short = "Fetch the logs of a container"
	c.Long = "Run logs command"

	pf := c.PersistentFlags()
	pf.StringVarP(&c.profile, "profile", "P", "dev", "--profile=test")
	pf.StringVarP(&c.project, "project", "p", "", "--project=project-name")
	pf.StringVarP(&c.app, "app", "a", "", "--app=my-app")
	pf.BoolVarP(&c.verbose, "verbose", "v", false, "--verbose")
	return c
}

func init() {
	app.Register(newLogsCommand)
}

//cube/client log --profile=test --project=demo --app=hello-world --sourcecode=java --verbose=true
func (c *logsCommand) Run(args []string) error {
	log.Debug("handle logs command")

	user, _, err := api.ReadConfig()
	if err != nil {
		return err
	}

	pss := &api.PipelineRequest{Name: c.app, Namespace: c.project, SourceCode: "_"}
	if _, err := api.StartInit(user, pss); err != nil {
		return err
	}

	verbose := "false"
	if c.verbose {
		verbose = "true"
	}
	var url = fmt.Sprintf("%s?namespace=%s&name=%s&sourcecode=%s&verbose=%s", api.GetBuildLogApi(user.Server), pss.Namespace, pss.Name, pss.SourceCode, verbose)
	if err := api.ClientLoop(url, api.BuildLogOut); err != nil {
		fmt.Println("[ERROR] log acquisition failed")
		os.Exit(0)
	}
	return nil
}
