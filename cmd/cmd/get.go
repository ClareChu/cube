package cmd


import (
	"fmt"
	"hidevops.io/cube/cmd/api"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/cli"
)

type getCommand struct {
	cli.SubCommand
	Info string
	Version string
	req  api.PipelineRequest
}

func newGetCommand() *getCommand {
	c := &getCommand{}
	c.Use = "get"
	c.Short = "get app"
	c.Long = "get app command"
	pf := c.PersistentFlags()
	pf.StringVarP(&c.Info, "info", "i", "", "")
	pf.StringVarP(&c.req.Version, "version", "v", "v1", "--your-app-version")
	return c
}

func init() {
	app.Register(newGetCommand)
}

func (c *getCommand) Run(args []string) (err error) {
	user, _, err := api.ReadConfig()
	if err != nil {
		return err
	}
	if err := api.GetApp(user, &c.req); err != nil {
		fmt.Println("[ERROR] get app err:", err)
		return err
	}
	return
}
