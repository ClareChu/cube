package builder

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPipelineUpdate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	pipeline := mio.NewPipeline(clientSet)
	dca := &v1alpha1.Pipeline{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world-v1",
			Namespace: "demo",
		},
	}
	_, err := pipeline.Create(dca)
	db := newPipelineService(pipeline)
	err = db.Update("hello-world-v1", "demo", "a", "success", "")
	assert.Equal(t, nil, err)
}
