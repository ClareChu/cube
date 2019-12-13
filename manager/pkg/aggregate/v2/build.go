package v2

import (
	"errors"
	"fmt"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"
	"strconv"
	"time"
)

type BuildInterface interface {
	WatchPod(name, namespace string) error
	Update(build *v1alpha1.Build, event, phase string) error
}

type Build struct {
	BuildInterface
	pod                *kube.Pod
	buildClient        *cube.Build
}

func init() {
	app.Register(NewBuildService)
}

func NewBuildService(pod *kube.Pod,
	buildClient *cube.Build) BuildInterface {
	build := &Build{
		pod:                pod,
		buildClient:        buildClient,
	}
	return build
}

func (b *Build) WatchPod(name, namespace string) error {
	log.Debugf("build config Watch :%v", name)
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	after := time.Duration(kubeWatchTimeout) * time.Minute
	timeout := int64(constant.TimeoutSeconds)
	listOptions := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", name),
	}
	w, err := b.pod.Watch(listOptions, namespace)
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
							return err
						}
					}
				}
				log.Infof("add event type :%v", pod.Name)
			case watch.Modified:
				pod := event.Object.(*corev1.Pod)
				if pod.Status.Phase == "Running" {
					for _, condition := range pod.Status.Conditions {
						if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
							time.Sleep(time.Second * 30)
							log.Infof("yes type :%v", pod.Name)
							return err
						}
					}
				}
				log.Debugf("update event type :%v", pod.Status)
			case watch.Deleted:
				return nil
			default:
				log.Error("Failed")
				return nil
			}
		}
	}
}

func (b *Build) Update(build *v1alpha1.Build, event, phase string) error {
	stage := v1alpha1.Stages{
		Name:                 event,
		StartTime:            time.Now().Unix(),
		DurationMilliseconds: time.Now().Unix() - build.ObjectMeta.CreationTimestamp.Unix(),
	}
	build.Status.Stages = append(build.Status.Stages, stage)
	build.Status.Phase = phase
	_, err := b.buildClient.Update(build.Name, build.Namespace, build)
	return err
}
