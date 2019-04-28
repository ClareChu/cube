package aggregate

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceAggregate interface {
	Create(namespace string) error
}

type Namespace struct {
	NamespaceAggregate
	namespace *kube.Namespace
}

func init() {
	app.Register(NewNamespace)
}

func NewNamespace(namespace *kube.Namespace) NamespaceAggregate {
	return &Namespace{
		namespace: namespace,
	}
}

func (n *Namespace) Create(ns string) error {
	namespace := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{Name: ns},
		Spec: v1.NamespaceSpec{
			Finalizers: []v1.FinalizerName{
				"kubernetes",
			},
		},
	}
	options := meta_v1.GetOptions{}
	_, err := n.namespace.Get(ns, options)
	if err != nil {
		_, err = n.namespace.Create(namespace)
	}
	return err
}
