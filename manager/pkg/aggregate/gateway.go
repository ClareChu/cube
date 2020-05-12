package aggregate

import (
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"os"
)

type GatewayAggregate interface {
	Create(gateway *v1alpha1.GatewayConfig) (err error)
}

type Gateway struct {
	GatewayAggregate
	ingressService service.IngressService
	routerService  service.RouterService
}

func init() {
	app.Register(NewGateway)
}

func NewGateway(ingressService service.IngressService, routerService service.RouterService) GatewayAggregate {
	return &Gateway{
		ingressService: ingressService,
		routerService:  routerService,
	}
}

func (s *Gateway) Create(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	if os.Getenv("OCP") == "openshift" {
		err = s.routerService.CreateIngress(gatewayConfig)
		return err
	}
	err = s.ingressService.CreateIngress(gatewayConfig)
	return
}
