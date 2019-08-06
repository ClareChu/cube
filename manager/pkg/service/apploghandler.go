package service

import (
	"errors"
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/websocket"
	"io"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
	"time"
)

type AppLogsHandler struct {
	at.ContextAware

	connection *websocket.Connection
	kubeClient KubeClient
}

func newAppLogsHandler(kubeClient KubeClient, connection *websocket.Connection) *AppLogsHandler {
	return &AppLogsHandler{
		connection: connection,
		kubeClient: kubeClient,
	}
}

func init() {
	app.Register(newAppLogsHandler)
}

func (h *AppLogsHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)

	namespace := h.connection.Context().FormValue("namespace")
	podname := h.connection.Context().FormValue("podname")
	tail := h.connection.Context().FormValue("tail")

	profile := h.connection.Context().FormValue("profile")
	version := h.connection.Context().FormValue("version")

	name := h.connection.Context().FormValue("name")
	new := h.connection.Context().FormValue("new")

	if (name != "" && podname != "") || (name == "" && podname == "") {
		if err := h.connection.EmitMessage([]byte("name and podnames must use one")); err != nil {
			fmt.Println("Error ", err)
		}
		return
	}

	var nameList []string
	if podname != "" {
		nameList = strings.Split(podname, ",")
	}

	if name != "" {
		//使用协程实时向客户端发送状态
		podMessage := make(chan PodMessage)
		nextStageNamespace, labelSelect, err := h.kubeClient.GetPipelineApp(namespace, name, new, profile, version, podMessage)
		if err != nil {
			log.Debugf("Error %v", err)
			return
		}

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

		podList, err := h.kubeClient.GetPodList(nextStageNamespace, meta_v1.ListOptions{
			LabelSelector: labelSelect,
		})
		if err != nil || len(podList.Items) == 0 {
			log.Debugf("Error %v", err)
			if err := h.connection.EmitMessage([]byte(fmt.Sprintf("[ERROR] label %s get pod failed,in namespace %s", labelSelect, nextStageNamespace))); err != nil {
				log.Debugf("Error %v", err)
				return
			}
			return
		}

		for _, pod := range podList.Items {
			nameList = append(nameList, pod.Name)
		}

		namespace = nextStageNamespace
	}

	tailInt64 := int64(0)
	if tail != "" {
		var err error
		tailInt64, err = strconv.ParseInt(tail, 10, 64)
		if err != nil {
			if err := h.connection.EmitMessage([]byte(fmt.Sprintf("string %s to int64 failed", tail))); err != nil {
				log.Debugf("Error %v", err)
				return
			}
		}
	}

	sendLogs := func(name string) {
		//获取 pod 日志输出流
		reader, err := h.kubeClient.GetLogs(namespace, name, tailInt64)
		if err != nil {
			log.Debugf("Error %v", err)
			//TODO WS send message//
			if err := h.connection.EmitMessage([]byte(fmt.Sprintf("%s,in namespace %s", string(err.Error()), namespace))); err != nil {
				log.Debugf("Error %v", err)
			}
			return
		}

		//持续向客户端发送日志信息
		for err == nil {
			str, err := reader.ReadString('\n')
			if err != nil {
				log.Debugf("Error %v", err)
				break
			}

			//TODO WS send message
			if err := h.connection.EmitMessage([]byte(string(str))); err != nil {
				log.Debugf("Error %v", err)
				break
			}
		}
		if err == io.EOF {
			log.Debugf("Error %v", err)
			return
		}
	}

	fmt.Println(fmt.Sprintf("display pod %v log,in namespace %s", nameList, namespace))
	for _, podName := range nameList {
		go sendLogs(podName)
	}

}

func (h *AppLogsHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}



// OnDisconnect is the websocket disconnection handler
func (h *AppLogsHandler) OnPing() {
	log.Debugf("Connection with ID: %v has been pinged!", h.connection.ID())
}

// OnDisconnect is the websocket disconnection handler
func (h *AppLogsHandler) OnPong() {
	log.Debugf("Connection with ID: %v has been ponged!", h.connection.ID())
}

