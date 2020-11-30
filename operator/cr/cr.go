package cr

import (
	"hidevops.io/cube/pkg/client/clientset/versioned"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type CubeManagerFunc func(clientSet versioned.Interface) CubeManagerInterface

type CubeManagerInterface interface {
	create()
}

var cubeManagerInterface = []CubeManagerFunc{
	NewServiceConfig,
	NewGatewayConfig,
	NewDeploymentConfig,
	NewPipelineConfig,
}

type CubeManagerCustomResourceDefinition struct {
	clientSet versioned.Interface
}

func InitCube() (crd *CubeManagerCustomResourceDefinition, err error) {
	clientSet, err := GetDefaultK8sClientSet()
	if err != nil {
		return
	}
	return &CubeManagerCustomResourceDefinition{clientSet: clientSet}, nil
}

func fakeInitCube() (crd *CubeManagerCustomResourceDefinition, err error) {
	clientSet := fake.NewSimpleClientset()
	return &CubeManagerCustomResourceDefinition{clientSet: clientSet}, nil
}

// 初始化crd的所有的资源
func (c *CubeManagerCustomResourceDefinition) Run() {
	for _, d := range cubeManagerInterface {
		//d(c.clientSet).create()
		d(c.clientSet).create()
	}
}

func GetDefaultK8sClientSet() (clientset versioned.Interface, err error) {
	var config *rest.Config
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
	clientset, err = versioned.NewForConfig(config)
	return
}

//GetKubeConfig
func GetKubeConfig() (kubeConfig string) {
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeConfig = os.Getenv("KUBECONFIG")
	}
	return
}
