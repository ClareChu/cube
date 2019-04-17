package cube

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/client-go/rest"
	"testing"
)

func TestCubeConfig(t *testing.T) {
	configuration := newConfiguration()
	config := &kube.RestConfig{
		Config: &rest.Config{},
	}
	result := configuration.Deployment(nil)
	assert.Equal(t, (*Deployment)(nil), result)

	result = configuration.Deployment(config)

	buildConfig := configuration.BuildConfig(nil)
	assert.Equal(t, (*BuildConfig)(nil), buildConfig)

	buildConfig = configuration.BuildConfig(config)

	build := configuration.Build(nil)
	assert.Equal(t, (*Build)(nil), build)

	build = configuration.Build(config)

	tests := configuration.Tests(nil)
	assert.Equal(t, (*Tests)(nil), tests)

	tests = configuration.Tests(config)

	testConfigs := configuration.TestConfig(nil)
	assert.Equal(t, (*TestConfig)(nil), testConfigs)

	testConfigs = configuration.TestConfig(config)

	notify := configuration.Notify(nil)
	assert.Equal(t, (*Notify)(nil), notify)

	notify = configuration.Notify(config)

	serviceConfig := configuration.ServiceConfig(nil)
	assert.Equal(t, (*ServiceConfig)(nil), serviceConfig)

	serviceConfig = configuration.ServiceConfig(config)

	gatewayConfig := configuration.GatewayConfig(nil)
	assert.Equal(t, (*GatewayConfig)(nil), gatewayConfig)

	gatewayConfig = configuration.GatewayConfig(config)

	deploymentConfig := configuration.DeploymentConfig(nil)
	assert.Equal(t, (*DeploymentConfig)(nil), deploymentConfig)

	deploymentConfig = configuration.DeploymentConfig(config)

	sourceConfig := configuration.SourceConfig(nil)
	assert.Equal(t, (*SourceConfig)(nil), sourceConfig)

	sourceConfig = configuration.SourceConfig(config)

	pipelineConfig := configuration.PipelineConfig(nil)
	assert.Equal(t, (*PipelineConfig)(nil), pipelineConfig)

	pipelineConfig = configuration.PipelineConfig(config)

	pipeline := configuration.Pipeline(nil)
	assert.Equal(t, (*Pipeline)(nil), pipeline)

	pipeline = configuration.Pipeline(config)

	imagesStream := configuration.ImageStream(nil)
	assert.Equal(t, (*ImageStream)(nil), imagesStream)

	pipelineConfig = configuration.PipelineConfig(config)

	imagesStream = configuration.ImageStream(nil)
	assert.Equal(t, (*ImageStream)(nil), imagesStream)

	imagesStream = configuration.ImageStream(config)
}
