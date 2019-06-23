package cmd

import (
	"hidevops.io/cube/cmd/api"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/cli"
	"hidevops.io/hiboot/pkg/log"
)

type appCommand struct {
	cli.SubCommand

	Req api.AppRequest
}

func newAppCommand() *appCommand {
	c := &appCommand{}
	c.Use = "create"
	c.Short = "create app"
	c.Long = "create app command"
	pf := c.PersistentFlags()
	pf.StringVarP(&c.Req.Name, "app", "a", "", "--name=your-app-name")
	pf.StringVarP(&c.Req.Namespace, "namespace", "n", "", "--namespace=your-k8s-namespace")
	pf.StringVarP(&c.Req.TemplateName, "template", "t", "", "--template=your-template-name")
	pf.StringVarP(&c.Req.Version, "version", "v", "v1", "--version=your-app-version")
	pf.StringVarP(&c.Req.Profile, "profile", "P", "", "--profile=your-app-profile")
	pf.StringVarP(&c.Req.Branch, "branch", "b", "", "--branch=your-app-branch")
	pf.StringVarP(&c.Req.Project, "project", "p", "", "--project=your-gitlab-app-name")
	pf.StringArrayVarP(&c.Req.Context, "context", "c", nil, "--context=your-context")
	pf.StringVarP(&c.Req.AppRoot, "appRoot", "r", "", "--appRoot=your-app-root")
	pf.StringArrayVarP(&c.Req.EnvVar, "envVar", "e", nil, "--env=your-env")
	pf.StringArrayVarP(&c.Req.Container.Command, "cmd", "c", nil, "--env=your-cmd")
	return c
}

func init() {
	app.Register(newAppCommand)
}

func (c *appCommand) Run(args []string) (err error) {
	log.Infof("start create app")
	user, _, err := api.ReadConfig()
	if err != nil {
		return err
	}
	var value1 complex64 = 3.1 + 5i
	value1 = value1 + 1
	err = api.App(&c.Req, user.Server, user.Token)
	return
}
