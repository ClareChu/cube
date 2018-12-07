package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/aggregate/mocks"
	builder "hidevops.io/mio/console/pkg/builder/mocks"
	"hidevops.io/mio/console/pkg/command"
	service "hidevops.io/mio/console/pkg/service/mocks"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestBuildCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)

	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	buildConfig := &v1alpha1.BuildConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	buildAggregate.Create(buildConfig, "hello-world", "v1")
}

func TestBuildCompile(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	cmd := &command.CompileCommand{
		Namespace: "demo",
		Name:      "hello-world",
	}
	buildConfigService.On("Compile", "hello-world.demo.svc", "7575", cmd).Return(nil)
	err := buildAggregate.Compile(build1)
	assert.Equal(t, nil, err)

}

func TestBuild_ImageBuild(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
		Spec: v1alpha1.BuildSpec{
			Tags: []string{"1"},
		},
	}
	cmd := &command.ImageBuildCommand{
		Namespace: "demo",
		Name:      "hello-world",
		Tags:      []string{"1:"},
	}
	buildConfigService.On("ImageBuild", "hello-world.demo.svc", "7575", cmd).Return(nil)
	err := buildAggregate.ImageBuild(build1)
	assert.Equal(t, nil, err)
	cmd1 := &command.ImagePushCommand{
		Namespace: "demo",
		Name:      "hello-world",
		Tags:      []string{"1:"},
	}
	buildConfigService.On("ImagePush", "hello-world.demo.svc", "7575", cmd1).Return(nil)
	err = buildAggregate.ImagePush(build1)
	assert.Equal(t, nil, err)
}

func TestBuild_ImagePush(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	cmd := &command.SourceCodePullCommand{
		Namespace: "demo",
		Name:      "hello-world",
		Url:       "/demo/.git",
	}
	buildConfigService.On("SourceCodePull", "hello-world.demo.svc", "7575", cmd).Return(nil)
	err := buildAggregate.SourceCodePull(build1)
	assert.Equal(t, nil, err)
}

func TestBuildCreateService(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	cmd := &command.ServiceNode{
		DeployData: kube.DeployData{
			Name:      "hello-world",
			NameSpace: "demo",
		},
	}
	_, err := build.Create(build1)
	assert.Equal(t, nil, err)
	buildNode.On("CreateServiceNode", cmd).Return(nil)
	err = buildAggregate.CreateService(build1)
	assert.Equal(t, nil, err)
}

func TestBuildDeployNode(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
			Labels: map[string]string{
				"name": "1",
			},
		},
	}
	cmd := &command.DeployNode{
		DeployData: kube.DeployData{
			Name:      "hello-world",
			NameSpace: "demo",
			Labels: map[string]string{
				"name": "1",
				"app":  "hello-world",
			},
		},
	}
	_, err := build.Create(build1)
	buildNode.On("Start", cmd).Return("", nil)
	err = buildAggregate.DeployNode(build1)
	assert.Equal(t, nil, err)
}

func TestBuildDeleteNode(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	build := mio.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	buildNode.On("DeleteDeployment", "hello-world", "demo").Return(nil)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
			Labels: map[string]string{
				"name": "1",
			},
		},
	}
	serviceAggregate.On("DeleteService", "hello-world", "demo").Return(nil)
	err := buildAggregate.DeleteNode(build1)
	_, err = build.Create(build1)
	assert.Equal(t, nil, err)
}
