package aggregate

import (
	"github.com/magiconair/properties/assert"
	builder "hidevops.io/cube/console/pkg/builder/mocks"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/cube/console/pkg/constant"
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
	pb.On("Update", "", "demo", "createService", "success", "").Return(nil)
	_, err = serviceAggregate.Create("hello-world", "", "demo", "java", "v1", "dev")
	assert.Equal(t, nil, err)
	err = service.Create("java", "", constant.TemplateDefaultNamespace, "")
	assert.Equal(t, nil, err)

	err = serviceAggregate.DeleteService("java", constant.TemplateDefaultNamespace)
	assert.Equal(t, nil, err)
}
