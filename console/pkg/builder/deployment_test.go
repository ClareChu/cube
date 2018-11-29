package builder

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestDeploymentUpdate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	dca := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world-v1",
			Namespace: "demo",
		},
	}
	_, err := deployment.Create(dca)
	db := newDeploymentService(deployment)
	err = db.Update("hello-world-v1", "demo", "a", "success")
	assert.Equal(t, nil, err)
}
