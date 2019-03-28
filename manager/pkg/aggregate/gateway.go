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
	kongService    service.KongService
}

func init() {
	app.Register(NewGateway)
}

func NewGateway(ingressService service.IngressService, kongService service.KongService) GatewayAggregate {
	return &Gateway{
		ingressService: ingressService,
		kongService:    kongService,
	}
}

func (s *Gateway) Create(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	if true {
		err = s.ingressService.CreateIngress(gatewayConfig)
		return
	}
	err = s.kongService.CreateKong(gatewayConfig)
	return
}
