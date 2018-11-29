package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/mio/console/pkg/aggregate/mocks"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestBuildConfigTemplate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	buildConfig := mio.NewBuildConfig(clientSet)
	buildAggregate := new(mocks.BuildAggregate)
	buildConfigAggregate := NewBuildConfigService(buildConfig, buildAggregate)
	bc := &command.BuildConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: "demo",
		},
	}
	build1 := &v1alpha1.BuildConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.BuildConfigKind,
			APIVersion: constant.BuildConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
			Labels: map[string]string{
				"CodeType": "java",
			},
		},
		Spec: v1alpha1.BuildSpec{
			Tags: []string{"/demo/hello-world"},
			App:  "hello-world",
		},
		Status: v1alpha1.BuildConfigStatus{
			LastVersion: 1,
		},
	}
	_, err := buildConfigAggregate.Template(bc)
	assert.Equal(t, nil, err)
	buildAggregate.On("Create", build1, "hello-world-1", "v1").Return(nil, nil)
	_, err = buildConfigAggregate.Create("hello-world", "hello-world-1", "demo", "java", "v1")
	assert.Equal(t, nil, err)
}
