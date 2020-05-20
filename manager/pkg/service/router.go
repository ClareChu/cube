package service

import (
	routeV1 "github.com/openshift/api/route/v1"
	v1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type RouterService interface {
	//CreateIngress(gatewayConfig *v1alpha1.GatewayConfig) (err error)
	CreateRouter(gatewayConfig *v1alpha1.GatewayConfig) (err error)
}

type Router struct {
	RouterService
}

func init() {
	app.Register(NewRouter)
}

func NewRouter() RouterService {
	return &Router{
	}
}

func (i *Router) CreateRouter(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	log.Debugf("create ingress name: %v  namespace: %v", gatewayConfig.Name, gatewayConfig.Namespace)
	restConfig, err := GetDefaultK8sClientSet()
	if err != nil {
		log.Errorf("get k8s rest client error :%v", err)
		return err
	}
	client, err := v1.NewForConfig(restConfig)
	if err != nil {
		log.Errorf("get k8s client set error :%v", err)
		return
	}
	routeAPI := client.Routes(gatewayConfig.Namespace)
	ops := metav1.GetOptions{}
	_, err = routeAPI.Get(gatewayConfig.Name, ops)
	if err == nil {
		log.Debug("get openshift router exist !!!")
		return err
	}
	log.Debugf("get openshift router err :%v", err)
	route := &routeV1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gatewayConfig.Name,
			Namespace: gatewayConfig.Namespace,
		},
		Spec: routeV1.RouteSpec{
			To: routeV1.RouteTargetReference{
				Kind: "Service",
				Name: gatewayConfig.Name,
			},
			Host: gatewayConfig.Spec.Hosts[0],
			//Path: gatewayConfig.Spec.Uris[0],
			Port: &routeV1.RoutePort{
				TargetPort: intstr.IntOrString{
					IntVal: gatewayConfig.Spec.Port,
				},
			},
		},
	}
	_, err = routeAPI.Create(route)
	log.Debugf("create openshift router error :%v", err)
	return
}

func GetDefaultK8sClientSet() (config *rest.Config, err error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		kubeConfig := GetKubeConfig()
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return
		}
	}

	return
}

//
func GetKubeConfig() (kubeConfig string) {
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeConfig = os.Getenv("KUBECONFIG")
	}
	return
}
