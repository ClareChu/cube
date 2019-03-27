package service

import (
	"fmt"
	cubev1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

type BuildConfigClient interface {
	Create(build *cubev1alpha1.Build) (config *cubev1alpha1.Build, err error)
	Get(name, namespace string) (config *cubev1alpha1.Build, err error)
	Watch(listOptions v1.ListOptions, namespace, name string) (watch.Interface, error)
	Delete(name, namespace string) error
	Update(name, namespace string, config *cubev1alpha1.Build) (*cubev1alpha1.Build, error)
	UpdateBuildStatus(namespace, name, eventType, status string) (*cubev1alpha1.Build, error)
}

const (
	SourceCodePull string = "souceCodePull"
	Compile        string = "compile"
	ImageBuild     string = "imageBuild"
	ImagePush      string = "imagePush"
	SourceCodeTest string = "souceCodeTest"
	Command        string = "command"

	Pulled    string = "pulled"
	Created   string = "created"
	Completed string = "completed"
	Failed    string = "failed"

	Success       = "success"
	Fail          = "fail"
	Running       = "running"
	Pending       = "pending"
	BuildPipeline = "buildPipeline"
)

type buildConfigClientImpl struct {
	BuildConfigClient
	build *cube.Build
}

func init() {
	app.Register(newBuildConfigClient)
}

func newBuildConfigClient(build *cube.Build) BuildConfigClient {
	return &buildConfigClientImpl{
		build: build,
	}
}

func (b *buildConfigClientImpl) Create(build *cubev1alpha1.Build) (config *cubev1alpha1.Build, err error) {

	config, err = b.build.Create(build)
	return
}

func (b *buildConfigClientImpl) Get(name, namespace string) (config *cubev1alpha1.Build, err error) {

	config, err = b.build.Get(name, namespace)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (b *buildConfigClientImpl) Watch(listOptions v1.ListOptions, namespace, name string) (watch.Interface, error) {

	listOptions.LabelSelector = fmt.Sprintf("app=%s", name)
	w, err := b.build.Watch(listOptions, namespace)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (b *buildConfigClientImpl) Delete(name, namespace string) error {
	return b.build.Delete(name, namespace)
}

func (b *buildConfigClientImpl) Update(name, namespace string, config *cubev1alpha1.Build) (*cubev1alpha1.Build, error) {
	result, err := b.build.Update(name, namespace, config)
	return result, err
}

func (b *buildConfigClientImpl) UpdateBuildStatus(namespace, name, eventType, status string) (*cubev1alpha1.Build, error) {
	bc, err := b.Get(name, namespace)
	if err != nil {
		return nil, err
	}

	bc.Status.Phase = status

	if len(bc.Status.Stages) == 0 {
		bc.Status.Stages = append(bc.Status.Stages, cubev1alpha1.Stages{
			Name:      eventType,
			StartTime: time.Now().Unix()},
		)
	}
	for i, stage := range bc.Status.Stages {

		if eventType == stage.Name {
			bc.Status.Stages[i].DurationMilliseconds = (time.Now().UnixNano() - stage.StartTime) / 1e6
			break
		}

		if len(bc.Status.Stages) == (i + 1) {
			bc.Status.Stages = append(bc.Status.Stages, cubev1alpha1.Stages{
				Name:      eventType,
				StartTime: time.Now().Unix()},
			)
		}
	}

	bc, err = b.Update(name, namespace, bc)
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}

	return bc, nil
}
