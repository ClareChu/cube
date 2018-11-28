package mio

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
)

type configuration struct {
	app.Configuration `depends:"kube"`
}

func init() {
	app.Register(newConfiguration)
}

func newConfiguration() *configuration {
	return &configuration{}
}

func (c *configuration) BuildConfig(restConfig *kube.RestConfig) *BuildConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewBuildConfig(clientSet)
}

func (c *configuration) Build(restConfig *kube.RestConfig) *Build {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewBuild(clientSet)
}

func (c *configuration) Pipeline(restConfig *kube.RestConfig) *Pipeline {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewPipeline(clientSet)
}

func (c *configuration) PipelineConfig(restConfig *kube.RestConfig) *PipelineConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewPipelineConfig(clientSet)
}

func (c *configuration) SourceConfig(restConfig *kube.RestConfig) *SourceConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewSourceConfig(clientSet)
}

func (c *configuration) DeploymentConfig(restConfig *kube.RestConfig) *DeploymentConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewDeploymentConfig(clientSet)
}

func (c *configuration) GatewayConfig(restConfig *kube.RestConfig) *GatewayConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewGatewayConfig(clientSet)
}

func (c *configuration) ServiceConfig(restConfig *kube.RestConfig) *ServiceConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewServiceConfig(clientSet)
}

func (c *configuration) Deployment(restConfig *kube.RestConfig) *Deployment {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewDeployment(clientSet)
}

func (c *configuration) Tests(restConfig *kube.RestConfig) *Tests {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewTestses(clientSet)
}

func (c *configuration) TestConfig(restConfig *kube.RestConfig) *testConfig {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewTestConfig(clientSet)
}

func (c *configuration) Notify(restConfig *kube.RestConfig) *Notify {
	clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
	if err != nil {
		log.Errorf("v1alpha1.NewForConfig %v", err)
		return nil
	}
	return NewNotify(clientSet)
}
