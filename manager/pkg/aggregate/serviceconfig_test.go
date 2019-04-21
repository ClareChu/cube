package aggregate

import (
	"github.com/stretchr/testify/assert"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestServiceConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	serviceClient := cube.NewServiceConfig(clientSet)
	client := kubeFake.NewSimpleClientset()
	service := kube.NewService(client)
	pb := new(builder.PipelineBuilder)
	serviceAggregate := NewServiceConfigService(serviceClient, service, pb)
	cmd := &command.ServiceConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: constant.TemplateDefaultNamespace,
		},
	}
	_, err := serviceAggregate.Template(cmd)
	assert.Equal(t, nil, err)
	pb.On("Update", "hello-world-1", "demo", "createService", "success", "").Return(nil)
	param := &command.PipelineReqParams{
		Name: "java",
		PipelineName: "hello-world-1",
		Namespace: "demo",
		EventType: "java",
		Version: "v1",
		Profile: "dev",
	}
	_, err = serviceAggregate.Create(param)
	assert.Equal(t, nil, err)
	err = service.Create("java", "", constant.TemplateDefaultNamespace, "")
	assert.Equal(t, nil, err)

	err = serviceAggregate.DeleteService("java", constant.TemplateDefaultNamespace)
	assert.Equal(t, nil, err)
}
