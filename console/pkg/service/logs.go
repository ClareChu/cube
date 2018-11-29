package service

import (
	"bufio"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

//go:generate mockgen -destination mock_logs.go -package aggregate hidevops.io/mio/console/pkg/aggregate KubeClient

type KubeClient interface {
	GetBuildConfigLastVersion(namespace, name string) (int, error)
	GetPodNameBylabel(namespace, label string) (string, error)
	WatchPodStatus(namespace, label string, intervals int, podMessage chan PodMessage) error
	GetLogs(namespace, name string, tail int64) (*bufio.Reader, error)
	GetPipelineConfig(namespace, name string) (config *v1alpha1.PipelineConfig, err error)
	WatchBuildStageStatus(namespace, label string, buildStatus chan v1alpha1.BuildStatus) error
	GetPodList(namespace string, opts meta_v1.ListOptions) (*corev1.PodList, error)
	GetPipelineApp(namespace, name, new, profile, version string, podMessage chan PodMessage) (string, string, error)
}

type KubeClientImpl struct {
	KubeClient
	pod                  *kube.Pod
	buildConfig          *mio.BuildConfig
	buildClient          *mio.Build
	pipelineConfigClient *mio.PipelineConfig
	pipelineClient       *mio.Pipeline
}

func init() {
	app.Register(newLogOutConfig)
}

func newLogOutConfig(buildClient *mio.Build, pod *kube.Pod, buildConfig *mio.BuildConfig, pipelineConfigClient *mio.PipelineConfig) KubeClient {
	return &KubeClientImpl{
		pod:                  pod,
		buildConfig:          buildConfig,
		pipelineConfigClient: pipelineConfigClient,
		buildClient:          buildClient,
	}
}

type PodMessage struct {
	Message string
	IsEnd   bool
}

//获取 buildconfig 的lastVersion 字段信息
func (l *KubeClientImpl) GetBuildConfigLastVersion(namespace, name string) (int, error) {
	bc, err := l.buildConfig.Get(name, namespace)
	if err != nil {
		fmt.Println("Error", err)
		return 0, err
	}

	return bc.Status.LastVersion, nil
}

//根据标签信息获取pod name
func (l *KubeClientImpl) GetPodNameBylabel(namespace, label string) (string, error) {

	podList, err := l.pod.GetPodList(namespace, meta_v1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		return "", err
	}
	for _, pod := range podList.Items {
		fmt.Println("INFO found pod name :", pod.Name)
	}

	if len(podList.Items) > 1 {
		return "", fmt.Errorf("the label %s matching pod should have only one", label)
	} else if len(podList.Items) == 0 {
		return "", fmt.Errorf("the label %s find pod failed", label)
	}
	return podList.Items[0].Name, nil
}

//获取pod Stream 信息，并且返回
func (l *KubeClientImpl) GetLogs(namespace, name string, tail int64) (*bufio.Reader, error) {
	podLogOptions := &v1.PodLogOptions{Follow: true}
	if tail != 0 {
		podLogOptions.TailLines = &tail
	}

	request, err := l.pod.GetPodLogs(namespace, name, podLogOptions)
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}

	byteReader, err := request.Stream()
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}

	reader := bufio.NewReader(byteReader)
	return reader, nil
}

//监听pod的状态，等待pod 正常或者失败 ，超时时间为10分钟
func (l *KubeClientImpl) WatchPodStatus(namespace, label string, intervals int, podMessage chan PodMessage) error {
	timeout := int64(3600)

	options := meta_v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  label,
	}

	w, err := l.pod.Watch(options, namespace)
	if err != nil {
		return err
	}

	for {
		select {
		case <-time.After(10 * time.Minute):
			fmt.Println(errors.New("Pod query timeout 10 minutes"))
			return errors.New("Pod query timeout 10 minutes")
		case event, ok := <-w.ResultChan():
			if !ok {
				fmt.Println("WatchPod resultChan: ", ok)
				return nil
			}
			pod := event.Object.(*corev1.Pod)
			if pod.Status.Phase == corev1.PodPending {
				podMessage <- PodMessage{Message: fmt.Sprintf("Event Type:%s Namespace:%s Pod:%s Status:%s \n", event.Type, pod.Namespace, pod.Name, pod.Status.Phase), IsEnd: false}
				continue
			}

			if pod.Status.Phase == corev1.PodRunning {
				//等待3秒确认pod是否再次异常
				time.Sleep(time.Second * 3)
				pod, err := l.pod.GetPods(namespace, pod.Name, meta_v1.GetOptions{})
				if err == nil {
					if pod.Status.Phase != corev1.PodRunning {
						continue
					}
				}

				fmt.Printf("pod %s has been running", pod.Name)
				if intervals == 0 {
					podMessage <- PodMessage{Message: fmt.Sprintf("Event Type:%s Namespace:%s Pod:%s Status:%s \n", event.Type, pod.Namespace, pod.Name, pod.Status.Phase), IsEnd: true}
					return nil
				}
				oldIntervals := time.Now().Unix() - pod.CreationTimestamp.Unix()
				if oldIntervals < int64(intervals) {
					podMessage <- PodMessage{Message: fmt.Sprintf("Event Type:%s Namespace:%s Pod:%s Status:%s \n", event.Type, pod.Namespace, pod.Name, pod.Status.Phase), IsEnd: true}
					return nil
				}

				continue
			} else {
				podMessage <- PodMessage{Message: fmt.Sprintf("Event Type:%s Namespace:%s Pod:%s Status:%s \n", event.Type, pod.Namespace, pod.Name, pod.Status.Phase), IsEnd: true}
				return fmt.Errorf("pod type %s status %s", string(event.Type), pod.Status.Phase)
			}
		}
	}
}

//获取 pipeline config信息
func (l *KubeClientImpl) GetPipelineConfig(namespace, name string) (config *v1alpha1.PipelineConfig, err error) {
	pc, err := l.pipelineConfigClient.Get(name, namespace)
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func (l *KubeClientImpl) WatchBuildStageStatus(namespace, label string, buildStatus chan v1alpha1.BuildStatus) error {
	timeout := int64(3600)
	listOptions := meta_v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", label),
	}
	w, err := l.buildClient.Watch(listOptions, namespace)
	if err != nil {
		return err
	}

	for {
		select {
		case <-time.After(10 * time.Minute):
			fmt.Println(errors.New("Pod query timeout 10 minutes"))
			return errors.New("Pod query timeout 10 minutes")
		case event, ok := <-w.ResultChan():
			if !ok {
				fmt.Println("WatchPipeline resultChan: ", ok)
				return nil
			}
			p := event.Object.(*v1alpha1.Build)
			buildStatus <- p.Status
			if p.Status.Phase == constant.Fail || p.Status.Phase == constant.Complete {
				return nil
			}
		}
	}
}

func (l *KubeClientImpl) GetPipelineApp(namespace, name, new, profile, version string, podMessage chan PodMessage) (string, string, error) {
	pipelineConfig, err := l.GetPipelineConfig(namespace, name)
	if err != nil {
		log.Debugf("Error %v", err)
		return "", "", err
	}
	if profile == "" {
		profile = pipelineConfig.Spec.Profile
	}

	if version == "" {
		version = pipelineConfig.Spec.Version
	}

	nextStageNamespace := fmt.Sprintf("%s-%s", namespace, profile)
	labelSelect := fmt.Sprintf("app=%s,version=%s", name, version)
	intervals := 0
	if new == "true" {
		intervals = 100
	}
	go l.WatchPodStatus(nextStageNamespace, labelSelect, intervals, podMessage)

	return nextStageNamespace, labelSelect, nil
}

func (l *KubeClientImpl) GetPodList(namespace string, opts meta_v1.ListOptions) (*corev1.PodList, error) {
	podList, err := l.pod.GetPodList(namespace, opts)
	if err != nil {
		return nil, err
	}
	return podList, nil
}
