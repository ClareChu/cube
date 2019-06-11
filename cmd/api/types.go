package api

import (
	"fmt"
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
	return  fmt.Sprintf("%s/app/name/%s", server, name)
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
	Name         string          `json:"name"`
	Project      string          `json:"project"`
	Namespace    string          `json:"namespace"`
	TemplateName string          `json:"templateName"`
	Profile      string          `json:"profile"`
	Branch       string          `json:"branch"`
	AppRoot      string          `json:"appRoot"`
	Context      []string        `json:"context"`
	Version      string          `json:"version"`
	Verbose      bool            `json:"verbose"`
	Watch        bool            `json:"watch"`
	Env          []corev1.EnvVar `json:"env"`
	EnvVar       []string        `json:"envVar"`
}

type User struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Server string `json:"server"`
}


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