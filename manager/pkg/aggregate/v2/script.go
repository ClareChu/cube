package v2

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
)

type ScriptInterface interface {
	Handle(build *v1alpha1.Build) error
}

type Script struct {
	CodeInterface
	configMaps         *kube.ConfigMaps
	buildConfigService service.BuildConfigService
	buildPackInterface dispatch.BuildPackInterface
}

func init() {
	app.Register(NewScriptService)
}

const COMPILE = "compile"

func NewScriptService(buildConfigService service.BuildConfigService,
	buildPackInterface dispatch.BuildPackInterface) ScriptInterface {
	script := &Script{
		buildConfigService: buildConfigService,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(COMPILE, script)
	return script
}

func (s *Script) Compile(build *v1alpha1.Build) error {
	var buildCommands []*command.BuildCommand
	for _, cmd := range build.Spec.CompileCmd {
		buildCommands = append(buildCommands, &command.BuildCommand{
			CodeType:    build.Spec.CodeType,
			ExecType:    cmd.ExecType,
			Script:      cmd.Script,
			CommandName: cmd.CommandName,
			Params:      cmd.CommandParams,
		})
	}

	command := &command.CompileCommand{
		Name:       build.Name,
		Namespace:  build.Namespace,
		CompileCmd: buildCommands,
		Context:    build.Spec.Context,
	}

	err := s.buildConfigService.Compile(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), constant.DefaultPort, command)
	return err
}
