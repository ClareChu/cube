package builder

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestDeploymentUpdate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	is := mio.NewImageStream(clientSet)
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
	assert.Equal(t, "imagestreams.mio.io \"hello-world\" not found", err.Error())




}
