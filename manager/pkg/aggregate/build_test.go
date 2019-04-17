package aggregate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	service "hidevops.io/cube/manager/pkg/service/mocks"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"os"
	"testing"
)

func TestBuildCreate(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	configMaps := kube.NewConfigMaps(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	buildConfig := &v1alpha1.BuildConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	buildAggregate.Create(buildConfig, "hello-world", "v1")
}

func TestBuildCompile(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	configMaps := kube.NewConfigMaps(client)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
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
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	configMaps := kube.NewConfigMaps(client)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
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
		Tags:      []string{"1:", "1:latest"},
	}
	buildConfigService.On("ImageBuild", "hello-world.demo.svc", "7575", cmd).Return(nil)
	err := buildAggregate.ImageBuild(build1)
	assert.Equal(t, nil, err)
	cmd1 := &command.ImagePushCommand{
		Namespace: "demo",
		Name:      "hello-world",
		Tags:      []string{"1:", "1:latest"},
	}
	buildConfigService.On("ImagePush", "hello-world.demo.svc", "7575", cmd1).Return(nil)
	err = buildAggregate.ImagePush(build1)
	assert.Equal(t, nil, err)
}

func TestBuild_ImagePush(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	configMaps := kube.NewConfigMaps(client)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
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
	data := map[string]string{
		constant.GitlabConstant: "",
	}
	configMaps.Create(constant.GitlabConstant, constant.TemplateDefaultNamespace, data)
	err := buildAggregate.SourceCodePull(build1)
	assert.Equal(t, nil, err)
}

func TestBuildCreateService(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	configMaps := kube.NewConfigMaps(client)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
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
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	configMaps := kube.NewConfigMaps(client)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
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
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	configMaps := kube.NewConfigMaps(client)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	buildNode.On("DeleteDeployment", "hello-world", "demo").Return(nil)
	build1 := &v1alpha1.Build{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
			Labels: map[string]string{
				"name": "hello-world",
			},
		},
	}
	serviceAggregate.On("DeleteService", "hello-world", "demo").Return(nil)
	err := buildAggregate.DeleteNode(build1)
	_, err = build.Create(build1)
	assert.Equal(t, nil, err)
}

func TestBuildSelector(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	configMaps := kube.NewConfigMaps(kubeFake.NewSimpleClientset())
	data := map[string]string{
		"API_VERSION":"api/v3",
		"BASE_URL":"https://github.com",
	}
	configMaps.Create(constant.GitlabConstant, constant.TemplateDefaultNamespace, data)
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	build := cube.NewBuild(clientSet)
	buildConfigService := new(service.BuildConfigService)
	buildNode := new(builder.BuildNode)
	client := kubeFake.NewSimpleClientset()
	pod := kube.NewPod(client)
	pipelineBuilder := new(builder.PipelineBuilder)
	replicationControllerAggregate := new(mocks.ReplicationControllerAggregate)
	serviceAggregate := new(mocks.ServiceConfigAggregate)
	buildAggregate := NewBuildService(configMaps, build, buildConfigService, buildNode, pod, pipelineBuilder, replicationControllerAggregate, serviceAggregate)
	b := &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
			},
		},
	}
	cmd := &command.DeployNode{
		DeployData: kube.DeployData{
			Labels: map[string]string{
				"app":  "",
				"name": "",
			},
		},
	}
	buildNode.On("Start", cmd).Return("", nil)
	err := buildAggregate.Selector(b)
	assert.Equal(t, nil, err)

	b = &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
				v1alpha1.Task{
					Name: constant.CreateService,
				},
			},
		},
		Status: v1alpha1.BuildStatus{
			Phase: constant.Success,
			Stages: []v1alpha1.Stages{
				v1alpha1.Stages{
					Name: constant.DeployNode,
				},
			},
		},
	}
	serviceCommand := &command.ServiceNode{}
	buildNode.On("CreateServiceNode", serviceCommand).Return(nil)
	err = buildAggregate.Selector(b)
	assert.Equal(t, "builds.cube.io \"\" not found", err.Error())
	b = &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
				v1alpha1.Task{
					Name: constant.CreateService,
				},
				v1alpha1.Task{
					Name: constant.CLONE,
				},
			},
		},
		Status: v1alpha1.BuildStatus{
			Phase: constant.Success,
			Stages: []v1alpha1.Stages{
				v1alpha1.Stages{
					Name: constant.DeployNode,
				},
				v1alpha1.Stages{
					Name: constant.CreateService,
				},
			},
		},
	}
	codeCommand := &command.SourceCodePullCommand{
		Url: "//.git",
	}

	buildConfigService.On("SourceCodePull", "..svc", "7575", codeCommand).Return(errors.New("1"))
	err = buildAggregate.Selector(b)
	assert.Equal(t, "1", err.Error())

	b = &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
				v1alpha1.Task{
					Name: constant.CreateService,
				},
				v1alpha1.Task{
					Name: constant.CLONE,
				},
				v1alpha1.Task{
					Name: constant.COMPILE,
				},
			},
		},
		Status: v1alpha1.BuildStatus{
			Phase: constant.Success,
			Stages: []v1alpha1.Stages{
				v1alpha1.Stages{
					Name: constant.DeployNode,
				},
				v1alpha1.Stages{
					Name: constant.CreateService,
				},
				v1alpha1.Stages{
					Name: constant.CLONE,
				},
			},
		},
	}
	buildCommand := &command.CompileCommand{CompileCmd: []*command.BuildCommand(nil), Namespace: "", Name: ""}
	buildConfigService.On("Compile", "..svc", "7575", buildCommand).Return(errors.New("1"))
	err = buildAggregate.Selector(b)
	assert.Equal(t, "1", err.Error())

	b = &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tags: []string{
				"",
			},
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
				v1alpha1.Task{
					Name: constant.CreateService,
				},
				v1alpha1.Task{
					Name: constant.CLONE,
				},
				v1alpha1.Task{
					Name: constant.COMPILE,
				},
				v1alpha1.Task{
					Name: constant.BuildImage,
				},
			},
		},
		Status: v1alpha1.BuildStatus{
			Phase: constant.Success,
			Stages: []v1alpha1.Stages{
				v1alpha1.Stages{
					Name: constant.DeployNode,
				},
				v1alpha1.Stages{
					Name: constant.CreateService,
				},
				v1alpha1.Stages{
					Name: constant.CLONE,
				},
				v1alpha1.Stages{
					Name: constant.BuildImage,
				},
			},
		},
	}
	imageBuildCommand := &command.ImageBuildCommand{App: "", S2IImage: "", Tags: []string{":", ":latest"}, DockerFile: []string(nil)}
	buildConfigService.On("ImageBuild", "..svc", "7575", imageBuildCommand).Return(errors.New("1"))
	err = buildAggregate.Selector(b)
	assert.Equal(t, "1", err.Error())

	b = &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tags: []string{
				"",
			},
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
				v1alpha1.Task{
					Name: constant.CreateService,
				},
				v1alpha1.Task{
					Name: constant.CLONE,
				},
				v1alpha1.Task{
					Name: constant.COMPILE,
				},
				v1alpha1.Task{
					Name: constant.BuildImage,
				},
				v1alpha1.Task{
					Name: constant.PushImage,
				},
			},
		},
		Status: v1alpha1.BuildStatus{
			Phase: constant.Success,
			Stages: []v1alpha1.Stages{
				v1alpha1.Stages{
					Name: constant.DeployNode,
				},
				v1alpha1.Stages{
					Name: constant.CreateService,
				},
				v1alpha1.Stages{
					Name: constant.CLONE,
				},
				v1alpha1.Stages{
					Name: constant.BuildImage,
				},
				v1alpha1.Stages{
					Name: constant.BuildImage,
				},
			},
		},
	}
	imagePushCommand := &command.ImagePushCommand{Tags: []string{":", ":latest"}}
	buildConfigService.On("ImagePush", "..svc", "7575", imagePushCommand).Return(errors.New("1"))
	err = buildAggregate.Selector(b)
	assert.Equal(t, "1", err.Error())

	b = &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Tags: []string{
				"",
			},
			Tasks: []v1alpha1.Task{
				v1alpha1.Task{
					Name: constant.DeployNode,
				},
				v1alpha1.Task{
					Name: constant.CreateService,
				},
				v1alpha1.Task{
					Name: constant.CLONE,
				},
				v1alpha1.Task{
					Name: constant.COMPILE,
				},
				v1alpha1.Task{
					Name: constant.BuildImage,
				},
				v1alpha1.Task{
					Name: constant.PushImage,
				},
				v1alpha1.Task{
					Name: constant.DeleteDeployment,
				},
			},
		},
		Status: v1alpha1.BuildStatus{
			Phase: constant.Success,
			Stages: []v1alpha1.Stages{
				v1alpha1.Stages{
					Name: constant.DeployNode,
				},
				v1alpha1.Stages{
					Name: constant.CreateService,
				},
				v1alpha1.Stages{
					Name: constant.CLONE,
				},
				v1alpha1.Stages{
					Name: constant.BuildImage,
				},
				v1alpha1.Stages{
					Name: constant.BuildImage,
				},
				v1alpha1.Stages{
					Name: constant.PushImage,
				},
			},
		},
	}
	buildNode.On("DeleteDeployment", "", "").Return(nil)
	serviceAggregate.On("DeleteService", "", "").Return(nil)
	err = buildAggregate.Selector(b)
	assert.Equal(t, "builds.cube.io \"\" not found", err.Error())
}
