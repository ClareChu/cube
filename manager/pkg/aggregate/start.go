package aggregate

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"strings"
)

type StartAggregate interface {
	Init(cmd *command.PipelineStart, propMap map[string]string) (err error)
	CreateProfile(namespace string) (err error)
	CreateSecret(cmd *command.PipelineStart, propMap map[string]string) (err error)
}

type Start struct {
	StartAggregate
	pipelineConfigAggregate PipelineConfigAggregate
	secretAggregate         SecretAggregate
	namespaceAggregate      NamespaceAggregate
	roleAggregate           RoleAggregate
	roleBindingAggregate    RoleBindingAggregate
	appService              service.AppService
	callbackAggregate       CallbackAggregate
}

func init() {
	app.Register(NewStartService)
}

func NewStartService(pipelineConfigAggregate PipelineConfigAggregate,
	secretAggregate SecretAggregate,
	namespaceAggregate NamespaceAggregate,
	roleAggregate RoleAggregate,
	roleBindingAggregate RoleBindingAggregate,
	appService service.AppService,
	callbackAggregate CallbackAggregate) StartAggregate {
	return &Start{
		pipelineConfigAggregate: pipelineConfigAggregate,
		secretAggregate:         secretAggregate,
		roleAggregate:           roleAggregate,
		roleBindingAggregate:    roleBindingAggregate,
		namespaceAggregate:      namespaceAggregate,
		appService:              appService,
		callbackAggregate:       callbackAggregate,
	}
}

func (s *Start) Init(cmd *command.PipelineStart, propMap map[string]string) (err error) {
	//TODO 获取cmd
	flag, err := s.appService.Init(cmd)
	if err != nil {
		log.Errorf("init app err :%v", err)
		return
	}
	log.Debugf("******** flag : %v, update: %v", flag, cmd.ForceUpdate)
	if flag && !cmd.ForceUpdate {
		err = s.callbackAggregate.Call(cmd)
		if err == nil {
			return
		}
	}

	//TODO 获取Namespace的值
	s.GetNamespace(cmd)
	//TODO init profile   create k8s namespace  serviceAccount default create role and roleBinding
	err = s.CreateProfile(cmd.Namespace)
	if err != nil {
		return
	}

	//TODO CREATE secret
	err = s.CreateSecret(cmd, propMap)
	if err != nil {
		return
	}

	if len(cmd.Context) == 0 || cmd.Context[0] == "" {
		go func() {
			_, err = s.pipelineConfigAggregate.StartPipelineConfig(cmd)
		}()
		return
	}
	for _, ct := range cmd.Context {
		log.Info(ct)
		paths := strings.Split(ct, "/")
		name := paths[len(paths)-1]
		cmd.Name = name
		cmd.Path = ct
		cmd.AppRoot = cmd.Name

		go func() {
			_, err = s.pipelineConfigAggregate.StartPipelineConfig(cmd)
		}()
	}
	return
}

//GetNamespace Profile
func (s *Start) GetNamespace(cmd *command.PipelineStart) {
	if cmd.Namespace == "" {
		if cmd.Profile == "" {
			cmd.Namespace = cmd.Project
		} else {
			cmd.Namespace = fmt.Sprintf("%s-%s", cmd.Project, cmd.Profile)
		}
	}
}

func (s *Start) CreateProfile(namespace string) (err error) {
	//TODO CREATE NAMESPACE
	err = s.namespaceAggregate.InitNamespace(namespace)
	if err != nil {
		return
	}
	//TODO create role
	err = s.roleAggregate.Create(constant.Default, namespace)
	if err != nil {
		return
	}
	//TODO create rolebinding
	err = s.roleBindingAggregate.Create(constant.Default, namespace)
	if err != nil {
		return
	}
	return
}

func (s *Start) CreateSecret(cmd *command.PipelineStart, propMap map[string]string) (err error) {
	secret := &command.Secret{
		Username:  propMap["username"],
		Password:  propMap["password"],
		Name:      cmd.AppRoot,
		Namespace: cmd.Namespace,
		Token:     propMap["access_token"],
	}
	err = s.secretAggregate.Create(secret)
	return err
}

func (s *Start) GetLanguagesType() error {

	return nil
}
