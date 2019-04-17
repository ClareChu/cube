package aggregate

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetNamespace(t *testing.T) {
	space := "demo"
	profile := "dev"
	n := GetNamespace(space, profile)
	assert.Equal(t, "demo-dev", n)
	profile = ""
	n = GetNamespace(space, profile)
	assert.Equal(t, "demo", n)
}

func TestTagTagImage(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	imageStream := cube.NewImageStream(clientSet)
	d := new(mocks.DeploymentBuilder)
	tag := NewTagService(imageStream, d)
	deploy := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				constant.DeploymentConfig: "hello-world",
			},
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	is := &v1alpha1.ImageStream{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	err := tag.TagImage(deploy)
	assert.Equal(t, "imagestreams.cube.io \"hello-world\" not found", err.Error())
	imageStream.Create(is)
	d.On("Update", "hello-world", "demo", "remoteDeploy", "success").Return(nil)
	err = tag.TagImage(deploy)
	assert.Equal(t, nil, err)
	deploy = &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				constant.DeploymentConfig: "hello-world",
			},
			Name:      "hello-world",
			Namespace: "demo",
		},
		Spec: v1alpha1.DeploymentConfigSpec{
			Profile: "dev",
		},
	}
	err = tag.TagImage(deploy)
	assert.Equal(t, nil, err)
}
