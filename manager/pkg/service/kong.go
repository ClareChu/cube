package service

import (
	"github.com/kevholditch/gokong"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
)

type KongService interface {
	CreateKong(gatewayConfig *v1alpha1.GatewayConfig) (err error)
}

type Kong struct {
	KongService
}

func init() {
	app.Register(NewKong)
}

func NewKong() KongService {
	return &Kong{
	}
}

func (k *Kong) CreateKong(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	apiRequest := &gokong.ApiRequest{
		Name:                   gatewayConfig.Name,
		Hosts:                  gatewayConfig.Spec.Hosts,
		Uris:                   gatewayConfig.Spec.Uris,
		UpstreamUrl:            gatewayConfig.Spec.UpstreamUrl,
		StripUri:               gatewayConfig.Spec.StripUri,
		PreserveHost:           gatewayConfig.Spec.PreserveHost,
		Retries:                gatewayConfig.Spec.Retries,
		UpstreamConnectTimeout: gatewayConfig.Spec.UpstreamConnectTimeout,
		UpstreamSendTimeout:    gatewayConfig.Spec.UpstreamSendTimeout,
		UpstreamReadTimeout:    gatewayConfig.Spec.UpstreamReadTimeout,
		HttpsOnly:              gatewayConfig.Spec.HttpsOnly,
		HttpIfTerminated:       gatewayConfig.Spec.HttpIfTerminated,
	}
	config := &gokong.Config{
		HostAddress: gatewayConfig.Spec.KongAdminUrl,
	}
	_, err = gokong.NewClient(config).Apis().Create(apiRequest)
	return
}