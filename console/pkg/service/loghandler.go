package service

import (
	"errors"
	"fmt"
	"hidevops.io/cube/console/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/websocket"
	"io"
	"strings"
	"time"
)

type LogsHandler struct {
	// annotation at.ContextAware
	at.ContextAware

	// websocket connection
	connection *websocket.Connection
	// kube client
	kubeClient KubeClient
}

func newLogsHandler(kubeClient KubeClient, connection *websocket.Connection) *LogsHandler {
	return &LogsHandler{connection: connection, kubeClient: kubeClient}
}

func init() {
	app.Register(newLogsHandler)
}

func (h *LogsHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)

	namespace := h.connection.Context().FormValue("namespace")
	name := h.connection.Context().FormValue("name")
	verbose := h.connection.Context().FormValue("verbose")
	if verbose == "" {
		verbose = "false"
	}

	//获取 buildConfig lastVersion
	lastVersion, err := h.kubeClient.GetBuildConfigLastVersion(namespace, name)
	if err != nil {
		log.Debugf("Error %v", err)
		//TODO WS send message
		if err := h.connection.EmitMessage([]byte("Information acquisition failed")); err != nil {
			fmt.Println("Error ", err)
		}
		return
	}
	pipelineConfig, err := h.kubeClient.GetPipelineConfig(namespace, name)
	if err != nil {
		log.Debugf("Error %v", err)
		//TODO WS send message
		if err := h.connection.EmitMessage([]byte("Information acquisition failed")); err != nil {
			fmt.Println("Error ", err)
		}
		return
	}

	label := fmt.Sprintf("%s-%s-%d", name, pipelineConfig.Spec.Version, lastVersion)

	//根据lastVersion 拼写标签信息 获取 podName
	podName, err := h.kubeClient.GetPodNameBylabel(namespace, fmt.Sprintf("app=%s", label))
	if err != nil {
		log.Debugf("Error %v", err)
		//TODO WS send message
		if err := h.connection.EmitMessage([]byte("Information acquisition failed")); err != nil {
			fmt.Println("Error ", err)
		}
		return
	}

	//使用协程实时向客户端发送状态
	podMessage := make(chan PodMessage)

	go h.kubeClient.WatchPodStatus(namespace, fmt.Sprintf("app=%s", label), 0, podMessage)

	f := true
	for f {
		select {
		case <-time.After(10 * time.Minute):
			fmt.Println(errors.New("pod query timeout 10 minutes"))
			//TODO WS send message
			if err := h.connection.EmitMessage([]byte("Pod query timeout 10 minutes")); err != nil {
				fmt.Println("Error ", err)
			}
			return

		case m := <-podMessage:
			//TODO WS send message
			if err := h.connection.EmitMessage([]byte(m.Message)); err != nil {
				fmt.Println("Error ", err)
				break
			}
			f = !m.IsEnd
		}
	}

	time.Sleep(time.Second)
	//获取 pod 日志输出流
	reader, err := h.kubeClient.GetLogs(namespace, podName, 0)
	if err != nil {
		log.Debugf("Error %v", err)
		return
	}

	//使用协程获取build状态信息
	var buildStatus = make(chan v1alpha1.BuildStatus)
	//defer close(buildStatus)
	go h.kubeClient.WatchBuildStageStatus(namespace, label, buildStatus)

	//持续向客户端发送日志信息
	var v, vs bool
	i := 1
	for err == nil {
		i++
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Debugf("Error %v", err)
			break
		}
		str = LogFormatting(str)

		if strings.Contains(str, "____________") && !vs && i <= 3 {
			vs = true
		}
		if verbose == "false" && !v && vs {
			v = strings.Contains(str, "register server shutdown on interrupt(CTRL+C/CMD+C)")
			if !v {
				continue
			}
		}
		//TODO WS send message
		if err := h.connection.EmitMessage([]byte(string(str))); err != nil {
			log.Debugf("Error %v", err)
			break
		}

		select {
		case b := <-buildStatus:
			fmt.Println(fmt.Sprintf("[INFO] %s phase %s:", b.Stages[len(b.Stages)-1].Name, b.Phase))
			if b.Phase == constant.Success && b.Stages[len(b.Stages)-1].Name == "imagePush" {
				//TODO WS send message
				if err := h.connection.EmitMessage([]byte("[INFO] end of build")); err != nil {
					log.Debugf("Error %v", err)
				}
				return
			}

			if b.Phase == constant.Fail {
				//TODO WS send message
				if err := h.connection.EmitMessage([]byte(fmt.Sprintf("[INFO] %s phase %s:", b.Stages[len(b.Stages)-1].Name, b.Phase))); err != nil {
					log.Debugf("Error %v", err)
				}
				return
			}
		default:
		}

	}

	if err == io.EOF {
		log.Debugf("Error %v", err)
		return
	}

}

func (h *LogsHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}

func LogFormatting(message string) string {

	if strings.Contains(message, "in namespace") || strings.Contains(message, "Server closed") {
		return ""
	}

	if strings.Contains(message, `{"stream"`) {
		if strings.Contains(message, "u003e") {
			return ""
		}
		message = strings.Replace(message, `{"stream":"`, "", -1)
		message = strings.Replace(message, `\n"}`, "", -1)
		message = strings.Replace(message, `\"`, "", -1)

	}

	if strings.Contains(message, `{"status"`) {
		if strings.Contains(message, "u003e") {
			return ""
		}
		message = strings.Replace(message, `"`, "", -1)
		message = strings.Replace(message, `{`, "", -1)
		message = strings.Replace(message, `}`, "", -1)
	}

	return message
}
