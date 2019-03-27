package aggregate

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicationControllerAggregate interface {
	Delete(name, namespace string) error
}

type ReplicationControllerAggregateImpl struct {
	replicationController *kube.ReplicationController
}

func init() {
	app.Register(NewReplicationControllerAggregate)
}

func NewReplicationControllerAggregate(replicationController *kube.ReplicationController) ReplicationControllerAggregate {
	return &ReplicationControllerAggregateImpl{
		replicationController: replicationController,
	}
}

func (r *ReplicationControllerAggregateImpl) Delete(name, namespace string) error {
	option := &metav1.DeleteOptions{}
	err := r.replicationController.Delete(name, namespace, option)
	return err
}
