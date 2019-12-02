package aggregate

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"
	"strconv"
	"time"
)

type CallbackAggregate interface {
	WatchPod(name, namespace string) error
	Call(cmd *command.PipelineStart) error
	Create(params *command.PipelineReqParams) (err error)
}

type Callback struct {
	CallbackAggregate
	pipelineBuilder          builder.PipelineBuilder
	pod                      *kube.Pod
	podBuilder               builder.PodBuilder
	pipelineFactoryInterface dispatch.PipelineFactoryInterface
}

const CallBack = "callback"

func init() {
	app.Register(NewCallbackService)
}

func NewCallbackService(pipelineBuilder builder.PipelineBuilder,
	podBuilder builder.PodBuilder,
	pod *kube.Pod,
	pipelineFactoryInterface dispatch.PipelineFactoryInterface) CallbackAggregate {
	callback := &Callback{
		podBuilder:               podBuilder,
		pipelineBuilder:          pipelineBuilder,
		pod:                      pod,
		pipelineFactoryInterface: pipelineFactoryInterface,
	}
	pipelineFactoryInterface.Add(CallBack, callback)
	return callback
}

func (v *Callback) Create(params *command.PipelineReqParams) (err error) {
	ready, err := v.podBuilder.GetPod(params.Name, params.Namespace)
	if !ready {
		err = v.WatchPod(params.Name, params.Namespace)
	}
	url := params.Ingress[0].Domain + params.Ingress[0].Path
	err = v.Send(params.Callback, params.Name, params.Namespace, params.Token, err.Error(), url, params.Id)
	if err != nil {
		v.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.CallBack, constant.Fail, "")
		return err
	}
	err = v.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.CallBack, constant.Success, "")
	return err
}

func (v *Callback) WatchPod(name, namespace string) error {
	log.Debugf("build config Watch :%v", name)
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	after := time.Duration(kubeWatchTimeout) * time.Minute
	timeout := int64(constant.TimeoutSeconds)
	listOptions := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", name),
	}
	w, err := v.pod.Watch(listOptions, namespace)
	if err != nil {
		return err
	}
	for {
		select {
		case <-time.After(after):
			return errors.New("pod query timeout 10 minutes")
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Info("WatchPod resultChan: ", ok)
				return nil
			}
			switch event.Type {
			case watch.Added:
				pod := event.Object.(*corev1.Pod)
				if pod.Status.Phase == "Running" {
					for _, condition := range pod.Status.Conditions {
						if condition.Type == corev1.PodReady {
							log.Infof("yes type :%v", pod.Name)
						}
					}
				}
				log.Infof("add event type :%v", pod.Name)
			case watch.Modified:
				pod := event.Object.(*corev1.Pod)
				if pod.Status.Phase == "Running" {
					for _, condition := range pod.Status.Conditions {
						if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
							log.Infof("yes type :%v", pod.Name)
							return errors.New("ready")
						}
					}
				}
				log.Debugf("update event type :%v", pod.Status)
			case watch.Deleted:
			default:
				log.Error("Failed")
				return errors.New("failed")
			}
		}
	}
}

type Data struct {
	model.RequestBody
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Url       string `json:"url"`
}

func (v *Callback) Send(callbackUrl, name, namespace, token, status, url string, id int) error {
	log.Debugf("******************************************************************")
	log.Debugf("callback url: %s", callbackUrl)
	rep := &model.BaseResponse{
		Code:    200,
		Message: "success",
		Data: &Data{
			Id:        id,
			Name:      name,
			Namespace: namespace,
			Status:    status,
			Url:       url,
		},
	}
	time.Sleep(5 * time.Second)
	_, body, errs := gorequest.New().Get(callbackUrl).Set(constant.Authorization, token).Send(rep).End()
	log.Infof("response : %s", body)
	if errs != nil {
		return errors.New("http get callback url")
	}
	return nil
}

func (v *Callback) Call(cmd *command.PipelineStart) error {
	log.Debugf("******************************************************************")
	ready, _ := v.podBuilder.GetPod(cmd.Name, cmd.Namespace)
	if !ready {
		return errors.New("pod is not ready **")
	}
	rep := &model.BaseResponse{
		Code:    200,
		Message: "success",
		Data: &Data{
			Id:        cmd.Id,
			Name:      cmd.Name,
			Namespace: cmd.Namespace,
			Status:    "success",
			Url:       cmd.Ingress[0].Domain + cmd.Ingress[0].Path,
		},
	}
	_, body, errs := gorequest.New().Get(cmd.Callback).Set(constant.Authorization, cmd.Token).Send(rep).End()
	log.Infof("response : %s", body)
	if errs != nil {
		return errors.New("http get callback url")
	}
	return nil
}
