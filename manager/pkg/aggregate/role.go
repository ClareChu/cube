package aggregate

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	rbac_v1 "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RoleAggregate interface {
	Create(name, namespace string) error
}

type Role struct {
	RoleAggregate
	role *kube.Role
}

func init() {
	app.Register(NewRole)
}

func NewRole(role *kube.Role) RoleAggregate {
	return &Role{
		role: role,
	}
}

func (r *Role) Create(name, namespace string) error {
	options := meta_v1.GetOptions{}
	role, err := r.role.Get(name, namespace, options)
	if err != nil {
		role = &rbac_v1.Role{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Rules: []rbac_v1.PolicyRule{
				rbac_v1.PolicyRule{
					Verbs:     []string{"get", "list", "watch", "update"},
					APIGroups: []string{"cube.io"},
					Resources: []string{"builds",
						"buildconfigs",
						"deployments",
						"deploymentconfigs",
						"gatewayconfigs",
						"imagestreams",
						"pipelines",
						"pipelineconfigs",
						"serviceconfigs",
						"sourceconfigs",
						"tests",
						"testconfigs",
					},
				},
				rbac_v1.PolicyRule{
					Verbs:     []string{"get", "list", "watch", "update"},
					APIGroups: []string{""},
					Resources: []string{"secrets"},
				},
			},
		}
		role, err = r.role.Create(namespace, role)
	}
	return err
}
