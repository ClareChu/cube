package service

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"
	"strconv"
	"time"
)

type DeployService interface {
	Update(replicasRequest *ReplicasRequest) (err error)
	Put(replicasRequest *ReplicasRequest) (err error)
}

type DeployServiceImpl struct {
	DeployService
	deployment *kube.Deployment
	pod        *kube.Pod
}

type ReplicasRequest struct {
	model.RequestBody
	Id        interface{}     `json:"id"`
	Name      string          `json:"name"`
	Namespace string          `json:"namespace"`
	Version   string          `json:"version"`
	App       string          `json:"app"`
	Token     string          `json:"token"`
	Replicas  *int32          `json:"replicas"`
	Status    corev1.PodPhase `json:"status"`
	Url       string          `json:"url"`
}

func init() {
	app.Register(newDeployCommand)
}

func newDeployCommand(deployment *kube.Deployment, pod *kube.Pod) DeployService {
	return &DeployServiceImpl{
		deployment: deployment,
		pod:        pod,
	}
}

func (a *DeployServiceImpl) Update(replicasRequest *ReplicasRequest) (err error) {
	option := metav1.GetOptions{}
	res, err := a.deployment.Get(replicasRequest.Name, replicasRequest.Namespace, option)
	if err != nil {
		log.Errorf("get deployment error:%v", err)
		return
	}

	res.Spec.Replicas = replicasRequest.Replicas
	return a.deployment.Update(res)
}

func (a *DeployServiceImpl) Put(replicasRequest *ReplicasRequest) (err error) {
	option := metav1.GetOptions{}
	res, err := a.deployment.Get(replicasRequest.App, replicasRequest.Namespace, option)
	if err != nil {
		log.Errorf("get deployment error:%v", err)
		return
	}
	go func() {
		res.Spec.Replicas = replicasRequest.Replicas
		err = a.deployment.Update(res)
		var condition corev1.PodPhase
		if *replicasRequest.Replicas != 0 {
			condition, _ = a.watch(replicasRequest.Name, replicasRequest.Namespace)
			replicasRequest.Status = condition
		} else {
			replicasRequest.Status = "success"
		}
		err = a.Send(replicasRequest)
	}()
	return
}

func (a *DeployServiceImpl) watch(app, namespace string) (corev1.PodPhase, error) {
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	if err != nil {
		return "", err
	}
	after := time.Duration(kubeWatchTimeout) * time.Minute
	timeout := int64(constant.TimeoutSeconds)
	listOptions := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", app),
	}
	w, err := a.pod.Watch(listOptions, namespace)
	for {
		select {
		case <-time.After(after):
			return "", errors.New("pod query timeout 10 minutes")
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Info("WatchPod resultChan: ", ok)
				return corev1.PodRunning, nil
			}
			switch event.Type {
			case watch.Added:
				pod := event.Object.(*corev1.Pod)
				if pod.Status.Phase == corev1.PodRunning {
					for _, condition := range pod.Status.Conditions {
						if condition.Type == corev1.PodReady {
							log.Debugf("yes type :%v", pod.Name)
							return corev1.PodRunning, nil
						}
					}
				}
				log.Infof("add event type :%v", pod.Name)
			case watch.Modified:
				pod := event.Object.(*corev1.Pod)
				if pod.Status.Phase == corev1.PodRunning {
					for _, condition := range pod.Status.Conditions {
						if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
							log.Debugf("yes type :%v", pod.Name)
							return corev1.PodRunning, nil
						}
					}
				}
				log.Debugf("update event type :%v", pod.Status)
			case watch.Deleted:
				log.Debugf("delete po comp ")
				return corev1.PodRunning, nil
			default:
				log.Error("Failed")
				return corev1.PodReasonUnschedulable, errors.New("failed")
			}
		}
	}
}

func (a *DeployServiceImpl) Send(replicas *ReplicasRequest) error {
	log.Debugf("******************************************************************")
	log.Debugf("callback url: %s", replicas.Url)
	rep := &model.BaseResponse{
		Code:    200,
		Message: "success",
		Data:    replicas,
	}
	time.Sleep(5 * time.Second)
	_, body, errs := gorequest.New().Get(replicas.Url).
		//Set(constant.Authorization, replicas.Token).
		Send(rep).
		Timeout(5 * time.Second).
		End()
	log.Infof("response : %s", body)
	if errs != nil {
		return errors.New("http get callback url")
	}
	return nil
}
