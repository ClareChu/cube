package builder

import (
	"errors"
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodBuilder interface {
	GetPod(name, namespace string) (ready bool, err error)
}

type Pod struct {
	PodBuilder
	pod *kube.Pod
}

func init() {
	app.Register(newPodService)
}

func newPodService(pod *kube.Pod) PodBuilder {
	return &Pod{
		pod: pod,
	}
}

func (p *Pod) GetPod(name, namespace string) (ready bool, err error) {
	log.Infof("*** get pod ****")
	listOptions := v1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", name),
	}
	pods, err := p.pod.GetPodList(namespace, listOptions)
	if err != nil || len(pods.Items) == 0 {
		return false, nil
	}
	if len(pods.Items[0].Status.ContainerStatuses) == 0 {
		return false, errors.New("running")
	}
	ready = pods.Items[0].Status.ContainerStatuses[0].Ready
	return ready, errors.New("ready")
}
