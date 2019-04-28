package aggregate

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	rbac_v1 "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RoleBindingAggregate interface {
	Create(name, namespace string) error
}

type RoleBinding struct {
	RoleBindingAggregate
	roleBinding *kube.RoleBinding
}

func init() {
	app.Register(NewRoleBinding)
}

func NewRoleBinding(roleBinding *kube.RoleBinding) RoleBindingAggregate {
	return &RoleBinding{
		roleBinding: roleBinding,
	}
}

func (rb *RoleBinding) Create(name, namespace string) error {
	options := meta_v1.GetOptions{}
	roleBinding, err := rb.roleBinding.Get(name, namespace, options)
	if err != nil {
		roleBinding = &rbac_v1.RoleBinding{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Subjects: []rbac_v1.Subject{
				rbac_v1.Subject{
					Kind:      "ServiceAccount",
					Name:      "default",
					Namespace: namespace,
				},
			},
			RoleRef: rbac_v1.RoleRef{
				Kind:     "Role",
				Name:     name,
				APIGroup: "rbac.authorization.k8s.io",
			},
		}
		roleBinding, err = rb.roleBinding.Create(namespace, roleBinding)
	}
	return nil
}

