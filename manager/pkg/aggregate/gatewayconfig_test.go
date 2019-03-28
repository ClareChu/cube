package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGatewayConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	gatewayClient := cube.NewGatewayConfig(clientSet)
	pb := new(builder.PipelineBuilder)
	gate := new(mocks.GatewayAggregate)
	gatewayAggregate := NewGatewayService(gatewayClient, pb, gate)
	gc := &v1alpha1.GatewayConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.GatewayConfigKind,
			APIVersion: constant.GatewayConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: "hidevopsio",
		},
	}

	gc1 := &v1alpha1.GatewayConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.GatewayConfigKind,
			APIVersion: constant.GatewayConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: "hidevopsio-dev",
			Labels: map[string]string{
				constant.PipelineConfigName: "",
				constant.Namespace:          "hidevopsio",
			},
		},
		Spec: v1alpha1.GatewaySpec{
			Uris: []string{
				"/hidevopsio/java",
			},
			UpstreamUrl:"http://java.hidevopsio-dev.svc:8080",

		},
	}
	_, err := gatewayClient.Create(gc)
	assert.Equal(t, nil, err)
	cmd := &command.GatewayConfig{}
	_, err = gatewayAggregate.Template(cmd)
	assert.Equal(t, nil, err)
	gate.On("Create", gc1).Return(nil)
	pb.On("Update", "", "hidevopsio", "createService", "success", "").Return(nil)
	_, err = gatewayAggregate.Create("java", "", constant.TemplateDefaultNamespace, "java", "v1", "dev")
	assert.Equal(t, nil, err)
}
