package api

import (
	"fmt"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"runtime"
	"strings"
)

const (
	DEFAULT_SERVER string = "http://cube.app.hidevops.io/websocket"

	CLI_DIR     string = ".cube"
	CONFIG_FILE string = "config"
	USERNAME    string = "Username "
	PASSWORD    string = "Password "
	CLI_VERSION string = "v1.0.0"

	END_STR string = "end of build"
)

func GetPipelineStartApi(server string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	return fmt.Sprintf("%s/pipelineConfig", server)
}

func GetCreateAppApi(server string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	return fmt.Sprintf("%s/app", server)
}

func GetLoginApi(server string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	return fmt.Sprintf("%s/login", server)
}

func GetBuildLogApi(server string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	url := strings.Replace(strings.Replace(fmt.Sprintf("%s/webSocket/buildLogs", server), "http://", "ws://", -1), "https://", "ws://", -1)
	return url
}

func GetAppLogApi(server string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	url := strings.Replace(strings.Replace(fmt.Sprintf("%s/webSocket/appLogs", server), "http://", "ws://", -1), "https://", "ws://", -1)
	return url
}

func GetAppApi(server, name string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	return fmt.Sprintf("%s/app/name/%s", server, name)
}

func GetSourceCodeTypeApi(server, name, namespace string) string {
	if server == "" {
		server = DEFAULT_SERVER
	}
	return fmt.Sprintf("%s/sourceConfig/namespace/%s/name/%s", server, namespace, name)
}

// GetCliUpdateApi
func GetCliUpdateApi(server string) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if server == "" {
		server = DEFAULT_SERVER
	}
	// TODO: cliUpdate to be changed on backend side
	return fmt.Sprintf("%s/cliUpdate/type/%s/arch/%s/version/%s", server, goos, goarch, CLI_VERSION)
}

var Message = make(chan string)

type PipelineRequest struct {
	Name         string           `json:"name"`
	Project      string           `json:"project"`
	Namespace    string           `json:"namespace"`
	TemplateName string           `json:"templateName"`
	Profile      string           `json:"profile"`
	Branch       string           `json:"branch"`
	AppRoot      string           `json:"appRoot"`
	Context      []string         `json:"context"`
	Version      string           `json:"version"`
	Verbose      bool             `json:"verbose"`
	Watch        bool             `json:"watch"`
	Container    corev1.Container `json:"container"`
	EnvVar       []string         `json:"envVar"`
	Ports        []string         `json:"ports"`

	Id                int      `json:"id"`
	ForceUpdate       bool     `json:"forceUpdate" default:"true"`
	//获取path 的目录
	Path    string `json:"path"`
	//gitUrl
	Url           string             `json:"url"`
	Ingress       []v1alpha1.Ingress `json:"ingress"`
	InitContainer corev1.Container   `json:"initContainer"`
	Volumes       v1alpha1.Volumes   `json:"volumes"`
	Callback      string             `json:"callback"`
	IsApp         bool               `json:"isApp"`
	Token         string             `json:"token"`
	Services      []v1alpha1.Service `json:"services"`
}

type User struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Server string `json:"server"`
}

type AppRequest struct {
	Name         string           `json:"name"`
	Namespace    string           `json:"namespace"`
	TemplateName string           `json:"templateName"`
	Version      string           `json:"version"`
	Profile      string           `json:"profile"`
	Branch       string           `json:"branch"`
	Context      []string         `json:"context"`
	AppRoot      string           `json:"appRoot"`
	Path         string           `json:"path"`
	Project      string           `json:"project"`
	Url          string           `json:"url"`
	EnvVar       []string         `json:"envVar"`
	Container    corev1.Container `json:"container"`
	Ports        []string         `json:"ports"`
}
