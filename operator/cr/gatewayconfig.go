package cr

import (
	"hidevops.io/cube/operator/client"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GatewayConfig struct {
	clientSet versioned.Interface
	Resource  string
}

const (
	GatewayConfigResource = "gatewayconfigs"
)

func NewGatewayConfig(clientSet versioned.Interface) CubeManagerInterface {
	return &GatewayConfig{
		clientSet: clientSet,
		Resource:  GatewayConfigResource,
	}
}

func (g *GatewayConfig) create() {
	ide := &v1alpha1.GatewayConfig{
		ObjectMeta: v1.ObjectMeta{
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                  "traefik",
				"traefik.frontend.rule.type":                   "PathPrefix",
				"traefik.ingress.kubernetes.io/rewrite-target": "/${namespace}/${name}",
			},
			Name:      IDEName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.GatewaySpec{
			Hosts: []string{
				"dev.apps.cloud2go.cn",
			},
			HttpIfTerminated:       false,
			Port:                   8443,
			PreserveHost:           true,
			Retries:                "5",
			StripUri:               true,
			UpstreamConnectTimeout: 60000,
			UpstreamReadTimeout:    60000,
			UpstreamSendTimeout:    60000,
		},
	}

	// create ide pipeline config 流水线
	//client.CreateRestClient(g.clientSet, ide, g.Resource)
	client.CreatGatewayConfig(g.clientSet, ide)

	idePath := &v1alpha1.GatewayConfig{
		ObjectMeta: v1.ObjectMeta{
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                  "traefik",
				"traefik.frontend.rule.type":                   "PathPrefix",
				"traefik.ingress.kubernetes.io/rewrite-target": "/",
			},
			Name:      IDEPATHName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.GatewaySpec{
			Hosts: []string{
				"dev.apps.cloud2go.cn",
			},
			HttpIfTerminated:       false,
			Port:                   8443,
			PreserveHost:           true,
			Retries:                "5",
			StripUri:               true,
			UpstreamConnectTimeout: 60000,
			UpstreamReadTimeout:    60000,
			UpstreamSendTimeout:    60000,
		},
	}

	// create ide pipeline config 流水线
	//client.CreateRestClient(g.clientSet, idePath, g.Resource)
	client.CreatGatewayConfig(g.clientSet, idePath)
}
