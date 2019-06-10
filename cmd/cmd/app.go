package cmd

import (
	"github.com/prometheus/common/log"
	"hidevops.io/cube/cmd/api"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/cli"
	corev1 "k8s.io/api/core/v1"
)

type AppRequest struct {
	Name         string          `json:"name"`
	Namespace    string          `json:"namespace"`
	TemplateName string          `json:"templateName"`
	Version      string          `json:"version"`
	Profile      string          `json:"profile"`
	Branch       string          `json:"branch"`
	Context      []string        `json:"context"`
	AppRoot      string          `json:"appRoot"`
	Path         string          `json:"path"`
	Project      string          `json:"project"`
	Url          string          `json:"url"`
	Env          []corev1.EnvVar `json:"env"`
	EnvVar       []string        `json:"envVar"`
}

type appCommand struct {
	cli.SubCommand

	Req AppRequest
}

func newAppCommand() *appCommand {
	c := &appCommand{}
	c.Use = "create"
	c.Short = "create app"
	c.Long = "create app command"
	pf := c.PersistentFlags()
	pf.StringVarP(&c.Req.Name, "name", "n", "", "--name=your-name")
	pf.StringVarP(&c.Req.Namespace, "namespace", "s", "", "--namespace=your-k8s-namespace")
	pf.StringVarP(&c.Req.TemplateName, "template", "t", "", "--template=your-template-name")
	pf.StringVarP(&c.Req.Version, "version", "v", "", "--version=your-app-version")
	pf.StringVarP(&c.Req.Profile, "profile", "P", "", "--profile=your-app-profile")
	pf.StringVarP(&c.Req.Branch, "branch", "b", "", "--branch=your-app-branch")
	pf.StringVarP(&c.Req.Project, "project", "p", "", "--project=your-gitlab-app-name")
	pf.StringArrayVarP(&c.Req.Context, "context", "c", nil, "--context=your-context")
	pf.StringVarP(&c.Req.AppRoot, "appRoot", "r", "", "--appRoot=your-app-root")
	pf.StringArrayVarP(&c.Req.EnvVar, "envVar", "e", nil, "--env=your-env")
	return c
}

func init() {
	app.Register(newAppCommand)
}

func (c *appCommand) Run(args []string) (err error) {
	log.Infof("create app :%v", c)
	user, _, err := api.ReadConfig()
	if err != nil {
		return err
	}
	err = api.App(&c.Req, api.GetCreateAppApi(user.Server), user.Token)
	return
}
