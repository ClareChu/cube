package aggregate

import (
	"github.com/magiconair/properties/assert"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestGatewayConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	gatewayClient := cube.NewGatewayConfig(clientSet)
	pb := new(builder.PipelineBuilder)

	client := kubeFake.NewSimpleClientset()
	ingress := kube.NewIngress(client)
	gatewayAggregate := NewGatewayService(gatewayClient, pb, ingress)
	gc := &v1alpha1.GatewayConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.GatewayConfigKind,
			APIVersion: constant.GatewayConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: constant.TemplateDefaultNamespace,
		},
	}
	_, err := gatewayClient.Create(gc)
	assert.Equal(t, nil, err)
	cmd := &command.GatewayConfig{}
	_, err = gatewayAggregate.Template(cmd)
	assert.Equal(t, nil, err)
	pb.On("Update", "", "demo", "createService", "fail", "").Return(nil)
	_, err = gatewayAggregate.Create("hello-world", "", "demo", "java", "v1", "dev")
	assert.Equal(t, nil, err)
}
