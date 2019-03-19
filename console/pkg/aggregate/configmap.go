package aggregate

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/api/core/v1"
)

type ConfigMapsAggregate interface {
	Create(name, namespace string, data map[string]string) (configMap *v1.ConfigMap, err error)
	Get(name, namespace string) (configMap *v1.ConfigMap, err error)
}

type ConfigMaps struct {
	ConfigMapsAggregate
	configMaps *kube.ConfigMaps
}

func init() {
	app.Register(NewConfigMapsService)
}

func NewConfigMapsService(configMaps *kube.ConfigMaps) ConfigMapsAggregate {
	return &ConfigMaps{
		configMaps: configMaps,
	}
}

func (c *ConfigMaps) Create(name, namespace string, data map[string]string) (configMap *v1.ConfigMap, err error) {
	configMap, err = c.configMaps.Create(name, namespace, data)
	return configMap, err
}

func (c *ConfigMaps) Get(name, namespace string) (configMap *v1.ConfigMap, err error) {
	configMap, err = c.configMaps.Get(name, namespace)
	return configMap, err
}