package builder

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPipelineUpdate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	pipeline := cube.NewPipeline(clientSet)
	dca := &v1alpha1.Pipeline{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world-v1",
			Namespace: "demo",
		},
	}
	db := newPipelineService(pipeline)
	err := db.Update("hello-world-v1", "demo", "a", "success", "")
	assert.Equal(t, errors.New("pipelines.cube.io \"hello-world-v1\" not found").Error(), err.Error())
	_, err = pipeline.Create(dca)
	err = db.Update("hello-world-v1", "demo", "a", "success", "")
	assert.Equal(t, nil, err)
}
