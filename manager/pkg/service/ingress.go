package service

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type IngressService interface {
	CreateIngress(gatewayConfig *v1alpha1.GatewayConfig) (err error)
	Delete(name, namespace string) (err error)
}

type Ingress struct {
	IngressService
	ingress *kube.Ingress
}

func init() {
	app.Register(NewIngress)
}

func NewIngress(ingress *kube.Ingress) IngressService {
	return &Ingress{
		ingress: ingress,
	}
}

func (i *Ingress) CreateIngress(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	log.Debugf("create ingress name: %v  namespace: %v", gatewayConfig.Name, gatewayConfig.Namespace)
	ing, err := i.ingress.Get(gatewayConfig.Namespace, gatewayConfig.Name, v1.GetOptions{})
	ingress := &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Annotations: gatewayConfig.ObjectMeta.Annotations,
			Name:        gatewayConfig.Name,
			Namespace:   gatewayConfig.Namespace,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: gatewayConfig.Spec.Hosts[0],
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: gatewayConfig.Spec.Uris[0],
									Backend: v1beta1.IngressBackend{
										ServiceName: gatewayConfig.Name,
										ServicePort: intstr.IntOrString{
											IntVal: gatewayConfig.Spec.Port,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if err != nil {
		log.Errorf("get ingress err :%v", err)
		_, err = i.ingress.Create(ingress)
		log.Infof("create ingress err: %v", err)
		return
	}
	ing.Spec = ingress.Spec
	_, err = i.ingress.Update(ingress)
	return
}

func (i *Ingress) Delete(name, namespace string) (err error) {
	log.Debugf("delete ingress")
	options := &v1.DeleteOptions{

	}
	return i.ingress.Delete(namespace, name, options)
}
