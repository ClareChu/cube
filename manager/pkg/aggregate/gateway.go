package aggregate

import (
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
)

type GatewayAggregate interface {
	Create(gateway *v1alpha1.GatewayConfig) (err error)
}

type Gateway struct {
	GatewayAggregate
	ingressService service.IngressService
}

func init() {
	app.Register(NewGateway)
}

func NewGateway(ingressService service.IngressService) GatewayAggregate {
	return &Gateway{
		ingressService: ingressService,
	}
}

func (s *Gateway) Create(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	err = s.ingressService.CreateIngress(gatewayConfig)
	return
}
