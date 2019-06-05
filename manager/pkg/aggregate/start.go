package aggregate

import (
	"fmt"
	//"gopkg.in/src-d/enry.v1"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
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
}

func init() {
	app.Register(NewStartService)
}

func NewStartService(pipelineConfigAggregate PipelineConfigAggregate,
	secretAggregate SecretAggregate,
	namespaceAggregate NamespaceAggregate,
	roleAggregate RoleAggregate,
	roleBindingAggregate RoleBindingAggregate) StartAggregate {
	return &Start{
		pipelineConfigAggregate: pipelineConfigAggregate,
		secretAggregate:         secretAggregate,
		roleAggregate:           roleAggregate,
		roleBindingAggregate:    roleBindingAggregate,
		namespaceAggregate:      namespaceAggregate,
	}
}

func (s *Start) Init(cmd *command.PipelineStart, propMap map[string]string) (err error) {
	//todo 通过URL部署项目
	/*	if cmd.Url != "" {
			if strings.Contains(cmd.Url, "https://") || strings.Contains(cmd.Url, "http://") {
				url := strings.Split(cmd.Url, "//")[1]
				cmd.Project = strings.Split(url, "/")[1]
				cmd.Namespace = strings.Split(url, "/")[2]
			}
		}
		langs := enry.GetLanguagesByFilename("Gemfile", []byte("<content>"), []string{})
		log.Info(langs)*/
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
		cmd.AppRoot = cmd.Name
		go func() {
			_, err = s.pipelineConfigAggregate.StartPipelineConfig(cmd)
		}()
		return
	}
	for _, ct := range cmd.Context {
		log.Info(ct)
		paths := strings.Split(ct, "/")
		name := paths[len(paths)-1]
		command := command.PipelineStart{
			Name:         name,
			Namespace:    cmd.Namespace,
			Version:      cmd.Version,
			Profile:      cmd.Profile,
			Path:         ct,
			AppRoot:      cmd.Name,
			TemplateName: cmd.TemplateName,
			Branch:       cmd.Branch,
			Project:      cmd.Project,
			Env:          cmd.Env,
		}
		go func() {
			_, err = s.pipelineConfigAggregate.StartPipelineConfig(&command)
		}()
	}
	return
}

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
		Name:      cmd.Name,
		Namespace: cmd.Namespace,
		Token:     propMap["access_token"],
	}
	err = s.secretAggregate.Create(secret)
	return err
}

func (s *Start) GetLanguagesType() error {

	return nil
}
