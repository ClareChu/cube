package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestReplicationControllerAggregateImplDelete(t *testing.T) {
	client := fake.NewSimpleClientset()
	rc := kube.NewReplicationController(client)
	rc.Create("demo", "hello-world", 1)
	ra := NewReplicationControllerAggregate(rc)
	err := ra.Delete("demo", "hello-world")
	assert.Equal(t, nil, err)
}
