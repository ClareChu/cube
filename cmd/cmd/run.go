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

	profile    string
	version    string
	project    string
	sourcecode string
	app        string
	branch     string
	verbose    bool
	watch      bool
}

func newRunCommand() *runCommand {
	c := &runCommand{}
	c.Use = "run"
	c.Short = "Run a command in a new pipeline"
	c.Long = "Run run command"

	pf := c.PersistentFlags()
	pf.StringVarP(&c.profile, "profile", "P", "", "--profile=dev")
	pf.StringVarP(&c.version, "version", "V", "", "--version=v1")
	pf.StringVarP(&c.project, "project", "p", "", "--project=project-name")
	pf.StringVarP(&c.app, "app", "a", "", "--app=my-app")
	pf.StringVarP(&c.branch, "branch", "b", "", "--branch=master")
	pf.StringVarP(&c.sourcecode, "sourcecode", "s", "", "--sourcecode=java")
	pf.BoolVarP(&c.verbose, "verbose", "v", false, "--verbose")
	pf.BoolVarP(&c.watch, "watch", "w", false, "--watch")
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

	pss := &api.PipelineStarts{
		Name:       c.app,
		Namespace:  c.project,
		SourceCode: c.sourcecode,
		Profile:    c.profile,
		Branch:     c.branch,
		Version:    c.version}
	if err := api.PipelineStart(user, pss); err != nil {
		return err
	}
	fmt.Println("[INFO] pipeline start succeed")
	verbose := "false"
	if c.verbose {
		verbose = "true"
	}

	if c.watch {
		time.Sleep(time.Second * 3)
		var url = fmt.Sprintf("%s?namespace=%s&name=%s&sourcecode=%s&verbose=%s", api.GetBuildLogApi(user.Server), pss.Namespace, pss.Name, pss.SourceCode, verbose)

		if err := api.ClientLoop(url, api.BuildLogOut); err != nil {
			fmt.Println("[ERROR] log acquisition failed")
			return err
		}
		fmt.Println("\n[INFO] image push success")
		fmt.Println("\n[INFO] Application logs:")
		time.Sleep(time.Second * 1)
		appUrl := fmt.Sprintf("%s?namespace=%s&name=%s&new=true&profile=%s&version=%s", api.GetAppLogApi(user.Server), pss.Namespace, pss.Name, pss.Profile, pss.Version)
		if err := api.ClientLoop(appUrl, api.LogOut); err != nil {
			fmt.Println("[ERROR] log acquisition failed")
			return err
		}
	}
	return nil
}
