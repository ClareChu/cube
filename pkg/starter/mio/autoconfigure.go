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
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewBuildConfig(clientSet)
	}
	return nil
}

func (c *configuration) Build(restConfig *kube.RestConfig) *Build {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewBuild(clientSet)
	}
	return nil
}

func (c *configuration) Pipeline(restConfig *kube.RestConfig) *Pipeline {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewPipeline(clientSet)
	}
	return nil
}

func (c *configuration) PipelineConfig(restConfig *kube.RestConfig) *PipelineConfig {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewPipelineConfig(clientSet)
	}
	return nil
}

func (c *configuration) SourceConfig(restConfig *kube.RestConfig) *SourceConfig {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewSourceConfig(clientSet)
	}
	return nil
}

func (c *configuration) DeploymentConfig(restConfig *kube.RestConfig) *DeploymentConfig {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewDeploymentConfig(clientSet)
	}
	return nil
}

func (c *configuration) GatewayConfig(restConfig *kube.RestConfig) *GatewayConfig {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewGatewayConfig(clientSet)
	}
	return nil
}

func (c *configuration) ServiceConfig(restConfig *kube.RestConfig) *ServiceConfig {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewServiceConfig(clientSet)
	}
	return nil
}

func (c *configuration) Deployment(restConfig *kube.RestConfig) *Deployment {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewDeployment(clientSet)
	}
	return nil
}

func (c *configuration) Tests(restConfig *kube.RestConfig) *Tests {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewTestses(clientSet)
	}
	return nil
}

func (c *configuration) TestConfig(restConfig *kube.RestConfig) *TestConfig {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewTestConfig(clientSet)
	}
	return nil
}

func (c *configuration) Notify(restConfig *kube.RestConfig) *Notify {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewNotify(clientSet)
	}
	return nil
}

func (c *configuration) ImageStream(restConfig *kube.RestConfig) *ImageStream {
	if restConfig != nil {
		clientSet, err := v1alpha1.NewForConfig(restConfig.Config)
		if err != nil {
			log.Errorf("v1alpha1.NewForConfig %v", err)
			return nil
		}
		return NewImageStream(clientSet)
	}
	return nil
}