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
	"time"
)

type runCommand struct {
	cli.SubCommand

	req api.PipelineRequest

	profile    string
	version    string
	project    string
	sourcecode string
	app        string
	branch     string
	context    []string
	verbose    bool
	watch      bool
}

func newRunCommand() *runCommand {
	c := &runCommand{}
	c.Use = "run"
	c.Short = "Run a command in a new pipeline"
	c.Long = "Run run command"

	pf := c.PersistentFlags()
	pf.StringVarP(&c.req.Profile, "profile", "P", "", "--profile=dev")
	pf.StringVarP(&c.req.Version, "version", "V", "", "--version=v1")
	pf.StringVarP(&c.req.Project, "project", "p", "", "--project=project-name")
	pf.StringVarP(&c.req.Namespace, "namespace", "n", "", "--namespace=project-name-dev")
	pf.StringVarP(&c.req.Name, "app", "a", "", "--app=my-app")
	pf.StringVarP(&c.req.Branch, "branch", "b", "", "--branch=master")
	pf.StringSliceVarP(&c.req.Context, "context", "c", nil, "--context=sub-module")
	pf.StringVarP(&c.req.SourceCode, "sourcecode", "s", "", "--sourcecode=java")
	pf.BoolVarP(&c.req.Verbose, "verbose", "v", false, "--verbose")
	pf.BoolVarP(&c.req.Watch, "watch", "w", false, "--watch")
	return c
}

func init() {
	app.Register(newRunCommand)
}

//cube run --project=demo --app=hello-world --sourcecode=java --verbose=true
func (c *runCommand) Run(args []string) error {
	log.Debug("handle run command")

	user, _, err := api.ReadConfig()
	if err != nil {
		return err
	}

	if err := api.PipelineStart(user, &c.req); err != nil {
		return err
	}
	fmt.Println("[INFO] pipeline start succeed")
	verbose := "false"
	if c.verbose {
		verbose = "true"
	}

	if c.watch {
		time.Sleep(time.Second * 3)
		var url = fmt.Sprintf("%s?namespace=%s&name=%s&sourcecode=%s&verbose=%s", api.GetBuildLogApi(user.Server), c.req.Namespace, c.req.Name, c.req.SourceCode, verbose)

		if err := api.ClientLoop(url, api.BuildLogOut); err != nil {
			fmt.Println("[ERROR] log acquisition failed")
			return err
		}
		fmt.Println("\n[INFO] image push success")
		fmt.Println("\n[INFO] Application logs:")
		time.Sleep(time.Second * 1)
		appUrl := fmt.Sprintf("%s?namespace=%s&name=%s&new=true&profile=%s&version=%s", api.GetAppLogApi(user.Server), c.req.Namespace, c.req.Name, c.req.Profile, c.req.Version)
		if err := api.ClientLoop(appUrl, api.LogOut); err != nil {
			fmt.Println("[ERROR] log acquisition failed")
			return err
		}
	}
	return nil
}
