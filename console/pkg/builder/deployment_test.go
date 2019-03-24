package builder

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/console/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestDeploymentUpdate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	deployment := cube.NewDeployment(clientSet)
	is := cube.NewImageStream(clientSet)
	client := kubeFake.NewSimpleClientset()
	deploy := kube.NewDeployment(client)
	dca := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world-v1",
			Namespace: "demo",
			Labels: map[string]string{
				constant.DeploymentConfig: "hello-world",
			},
		},
	}
	_, err := deployment.Create(dca)
	db := newDeploymentService(deployment, deploy, is)
	err = db.Update("hello-world-v1", "demo", "a", "success")
	assert.Equal(t, nil, err)
	err = db.CreateApp(dca)
	assert.Equal(t, "imagestreams.cube.io \"hello-world\" not found", err.Error())

}
