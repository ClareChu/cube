package aggregate

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/idgen"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/builder"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/console/pkg/service"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"
	"strconv"
	"time"
)

type BuildAggregate interface {
	Create(buildConfig *v1alpha1.BuildConfig, pipelineName, version string) (build *v1alpha1.Build, err error)
	Watch(name, namespace string) (build *v1alpha1.Build, err error)
	SourceCodePull(build *v1alpha1.Build) error
	Compile(build *v1alpha1.Build) error
	ImageBuild(build *v1alpha1.Build) error
	CreateService(build *v1alpha1.Build) error
	DeployNode(build *v1alpha1.Build) error
	Selector(build *v1alpha1.Build) (err error)
	Update(build *v1alpha1.Build, event, phase string) error
	WatchPod(name, namespace string) error
	DeleteNode(build *v1alpha1.Build) error
	ImagePush(build *v1alpha1.Build) error
}

type Build struct {
	BuildAggregate
	buildClient                    *mio.Build
	buildConfigService             service.BuildConfigService
	buildNode                      builder.BuildNode
	pod                            *kube.Pod
	pipelineBuilder                builder.PipelineBuilder
	replicationControllerAggregate ReplicationControllerAggregate
	serviceConfigAggregate         ServiceConfigAggregate
}

func init() {
	app.Register(NewBuildService)
}

func NewBuildService(buildClient *mio.Build, buildConfigService service.BuildConfigService, buildNode builder.BuildNode, pod *kube.Pod, pipelineBuilder builder.PipelineBuilder, replicationControllerAggregate ReplicationControllerAggregate, serviceConfigAggregate ServiceConfigAggregate) BuildAggregate {
	return &Build{
		buildClient:                    buildClient,
		buildConfigService:             buildConfigService,
		buildNode:                      buildNode,
		pod:                            pod,
		pipelineBuilder:                pipelineBuilder,
		replicationControllerAggregate: replicationControllerAggregate,
		serviceConfigAggregate:         serviceConfigAggregate,
	}
}

func (s *Build) Create(buildConfig *v1alpha1.BuildConfig, pipelineName, version string) (build *v1alpha1.Build, err error) {
	log.Debugf("build config create :%v", buildConfig)
	number := fmt.Sprintf("%d", buildConfig.Status.LastVersion)
	nameVersion := fmt.Sprintf("%s-%s-%s", buildConfig.Name, version, number)
	build = new(v1alpha1.Build)
	copier.Copy(build, buildConfig)
	build.TypeMeta = v1.TypeMeta{
		Kind:       constant.BuildKind,
		APIVersion: constant.BuildApiVersion,
	}
	build.ObjectMeta = v1.ObjectMeta{
		Name:      nameVersion,
		Namespace: buildConfig.Namespace,
		Labels: map[string]string{
			constant.App:             nameVersion,
			constant.Number:          number,
			constant.BuildConfigName: buildConfig.Name,
			constant.PipelineName:    pipelineName,
			constant.Version:         version,
			constant.Name:            buildConfig.Name,
		},
	}
	config, err := s.buildClient.Create(build)
	log.Info("............config err:", config)
	if err != nil {
		log.Errorf("create build error :%v", err)
		return
	} else {
		_, err = s.Watch(config.Name, config.Namespace)
	}
	return config, err
}

func (s *Build) Watch(name, namespace string) (build *v1alpha1.Build, err error) {
	log.Debug("build config Watch :%v", name)
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	after := time.Duration(kubeWatchTimeout) * time.Minute

	timeout := int64(constant.TimeoutSeconds)
	if err != nil {
		err = errors.New("kube watch time out ")
		return
	}
	option := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", name),
	}
	w, err := s.buildClient.Watch(option, namespace)
	if err != nil {
		return
	}
	for {
		select {
		case <-time.After(after):
			return nil, errors.New("pod query timeout 10 minutes")
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Info(" build watch resultChan: ", ok)
				return
			}
			if event.Type == watch.Deleted {
				log.Info("Deleted: ", event.Object)
				return
			} else if event.Type == watch.Added || event.Type == watch.Modified {
				build = event.Object.(*v1alpha1.Build)
				if build.Status.Phase == constant.Fail {
					return
				}
				err = s.Selector(build)
				if err != nil {
					log.Errorf("add selector : %v", err)
					return
				}
				log.Infof("event type :%v", build.Status)
			} else {
				log.Info("build default ")
			}

		}
	}
}

func (s *Build) SourceCodePull(build *v1alpha1.Build) error {
	command := &command.SourceCodePullCommand{
		CloneType: build.Spec.CloneType,
		Url:       fmt.Sprintf("%s/%s/%s.git", build.Spec.CloneConfig.Url, build.Namespace, build.Labels[constant.BuildConfigName]),
		Branch:    build.Spec.CloneConfig.Branch,
		DstDir:    build.Spec.CloneConfig.DstDir,
		Username:  build.Spec.CloneConfig.Username,
		Password:  build.Spec.CloneConfig.Password,
		Namespace: build.Namespace,
		Name:      build.Name,
	}

	err := s.buildConfigService.SourceCodePull(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), "7575", command)
	return err
}

func (s *Build) Compile(build *v1alpha1.Build) error {
	var buildCommands []*command.BuildCommand
	for _, cmd := range build.Spec.CompileCmd {
		buildCommands = append(buildCommands, &command.BuildCommand{
			CodeType:    build.Spec.CodeType,
			ExecType:    cmd.ExecType,
			Script:      cmd.Script,
			CommandName: cmd.CommandName,
			Params:      cmd.CommandParams,
		})
	}

	command := &command.CompileCommand{
		Name:       build.Name,
		Namespace:  build.Namespace,
		CompileCmd: buildCommands,
	}

	err := s.buildConfigService.Compile(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), "7575", command)
	return err
}

func (s *Build) ImageBuild(build *v1alpha1.Build) error {
	id, err := idgen.NextString()
	if err != nil {
		log.Errorf("id err :{}", id)
	}
	command := &command.ImageBuildCommand{
		App:        build.Spec.App,
		S2IImage:   build.Spec.BaseImage,
		Tags:       []string{build.Spec.Tags[0] + ":" + build.ObjectMeta.Labels[constant.Number]},
		DockerFile: build.Spec.DockerFile,
		Name:       build.Name,
		Namespace:  build.Namespace,
		Username:   build.Spec.DockerAuthConfig.Username,
		Password:   build.Spec.DockerAuthConfig.Password,
	}
	log.Infof("build ImageBuild :%v", command)
	err = s.buildConfigService.ImageBuild(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), "7575", command)
	return err
}

func (s *Build) ImagePush(build *v1alpha1.Build) error {
	command := &command.ImagePushCommand{
		Tags:      []string{build.Spec.Tags[0] + ":" + build.ObjectMeta.Labels[constant.Number]},
		Name:      build.Name,
		Namespace: build.Namespace,
		Username:  build.Spec.DockerAuthConfig.Username,
		Password:  build.Spec.DockerAuthConfig.Password,
	}
	log.Infof("ImagePush :%v", command)
	err := s.buildConfigService.ImagePush(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), "7575", command)
	return err
}

func (s *Build) CreateService(build *v1alpha1.Build) error {
	log.Infof("aggregate create service %v", build)
	phase := constant.Success
	command := &command.ServiceNode{
		DeployData: kube.DeployData{
			Name:           build.Name,
			App:            build.ObjectMeta.Labels["name"],
			NameSpace:      build.Namespace,
			Replicas:       build.Spec.DeployData.Replicas,
			Labels:         build.Spec.DeployData.Labels,
			Image:          build.Spec.BaseImage,
			Ports:          build.Spec.DeployData.Ports,
			HostPathVolume: build.Spec.DeployData.HostPathVolume,
		},
	}
	err := s.buildNode.CreateServiceNode(command)
	if err != nil {
		log.Errorf("deploy hinode err :%v", err)
		phase = constant.Fail
		return err
	}
	err = s.Update(build, constant.CreateService, phase)
	return err
}

func (s *Build) DeployNode(build *v1alpha1.Build) error {
	phase := constant.Success
	command := &command.DeployNode{
		DeployData: kube.DeployData{
			Name:      build.Name,
			NameSpace: build.Namespace,
			Replicas:  build.Spec.DeployData.Replicas,
			Labels: map[string]string{
				constant.App:  build.Name,
				constant.Name: build.ObjectMeta.Labels[constant.Name],
			},
			Image:          build.Spec.BaseImage,
			Ports:          build.Spec.DeployData.Ports,
			HostPathVolume: build.Spec.DeployData.HostPathVolume,
			Envs:           build.Spec.DeployData.Envs,
			NodeName:       build.Spec.DeployData.Envs["NODE_NAME"],
		},
	}
	_, err := s.buildNode.Start(command)
	if err != nil {
		log.Errorf("deploy hinode err :%v", err)
		return err
	}
	err = s.WatchPod(build.Name, build.Namespace)
	if err != nil {
		phase = constant.Fail
	}

	err = s.Update(build, constant.DeployNode, phase)
	return err
}

func (s *Build) Selector(build *v1alpha1.Build) (err error) {
	tak, err := GetTask(build.Status.Stages, build.Spec.Tasks, build.Status.Phase)
	if err != nil {
		return
	}
	switch tak.Name {
	case constant.DeployNode:
		err := s.DeployNode(build)
		if err != nil {
			log.Errorf("deploy node %v", err)
		}
	case constant.CreateService:
		err = s.CreateService(build)
		if err != nil {
			log.Errorf("create service %v", err)
		}
	case constant.CLONE:
		err = s.SourceCodePull(build)
		if err != nil {
			log.Errorf("source code pull  %v", err)
		}
	case constant.COMPILE:
		err = s.Compile(build)
		if err != nil {
			log.Errorf("Compile %v", err)
		}
	case constant.BuildImage:
		err = s.ImageBuild(build)
		if err != nil {
			log.Errorf("Image Build %v", err)
		}
	case constant.PushImage:
		err = s.ImagePush(build)
		if err != nil {
			log.Errorf("Image Push %v", err)
		}
	case constant.DeleteDeployment:
		err = s.DeleteNode(build)
		if err != nil {
			log.Errorf("Delete Node %v", err)
		}
	case constant.Ending:
		err = s.Update(build, "", constant.Complete)
		log.Info("update pipeline aggregate")
		err = s.pipelineBuilder.Update(build.ObjectMeta.Labels[constant.PipelineName], build.Namespace, constant.BuildPipeline, constant.Success, build.ObjectMeta.Labels[constant.Number])
		err = fmt.Errorf("build is ending")
	default:

	}
	return
}

func GetTask(stages []v1alpha1.Stages, tasks []v1alpha1.Task, phase string) (tak v1alpha1.Task, err error) {
	if len(stages) == 0 {
		if len(tasks) == 0 {
			err = fmt.Errorf("tasks is len equ 0")
			return
		}
		tak = tasks[0]
	} else if phase == constant.Success && len(stages) != len(tasks) {
		tak = tasks[len(stages)]
	} else if len(stages) == len(tasks) {
		tak.Name = constant.Ending
	}
	return
}

func (s *Build) Update(build *v1alpha1.Build, event, phase string) error {
	stage := v1alpha1.Stages{
		Name:                 event,
		StartTime:            time.Now().Unix(),
		DurationMilliseconds: time.Now().Unix() - build.ObjectMeta.CreationTimestamp.Unix(),
	}
	build.Status.Stages = append(build.Status.Stages, stage)
	build.Status.Phase = phase
	_, err := s.buildClient.Update(build.Name, build.Namespace, build)
	return err
}

func (s *Build) WatchPod(name, namespace string) error {
	log.Debugf("build config Watch :%v", name)
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	after := time.Duration(kubeWatchTimeout) * time.Minute
	timeout := int64(constant.TimeoutSeconds)
	listOptions := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", name),
	}
	w, err := s.pod.Watch(listOptions, namespace)
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
				log.Infof("update event type :%v", pod.Status)
			case watch.Deleted:
				log.Info("Deleted: ", event.Object)
				return nil
			default:
				log.Error("Failed")
				return nil
			}
		}
	}
}

func (s *Build) DeleteNode(build *v1alpha1.Build) error {
	phase := constant.Success
	//TODO delete deployment config
	err := s.buildNode.DeleteDeployment(build.ObjectMeta.Labels["name"], build.Namespace)
	//TODO delete service
	err = s.serviceConfigAggregate.DeleteService(build.ObjectMeta.Labels["name"], build.Namespace)
	if err != nil {
		phase = constant.Fail
	}
	err = s.Update(build, constant.DeleteDeployment, phase)
	return err
}
