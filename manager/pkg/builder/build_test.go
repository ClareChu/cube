package builder

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestBuildNodeImplCreateServiceNode(t *testing.T) {
	client := fake.NewSimpleClientset()
	deployment := kube.NewDeployment(client)
	service := kube.NewService(client)
	replicaSet := kube.NewReplicaSet(client)
	build := newDeploymentConfig(deployment, service, replicaSet)
	s := &command.ServiceNode{}
	err := build.CreateServiceNode(s)
	assert.Equal(t, nil, err)

	d := &command.DeployNode{}
	_, err = build.Start(d)
	assert.Equal(t, nil, err)
	request := &kube.DeployRequest{
		Namespace: "hello-world",
		App:       "demo",
		Version:   "v1",
	}
	_, err = deployment.Deploy(request)
	assert.Equal(t, nil, err)
	set := &v1.ReplicaSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "demo-v1",
			Namespace: "hello-world",
			Labels: map[string]string{
				constant.App: "demo-v1",
			},
		},
	}
	_, err = replicaSet.Create(set)
	err = build.DeleteDeployment("demo-v1", "hello-world")
	assert.Equal(t, nil, err)

	err = build.Update("demo-v1", "hello-world")
	assert.Equal(t, nil, err)
}
