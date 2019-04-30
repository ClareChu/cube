package aggregate

import (
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceAggregate interface {
	InitNamespace(namespace string) error
	Create(namespace string) error
	HarborCreate(namespace string) error
}

type HarborNamespace struct {
	ProjectName string   `json:"project_name"`
	Metadata    Metadata `json:"metadata"`
}

type Metadata struct {
	Public string `json:"public"`
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

func (n *Namespace) InitNamespace(namespace string) error {
	err := n.Create(namespace)
	if err != nil {
		return nil
	}
	err = n.HarborCreate(namespace)
	return err
}

func (n *Namespace) Create(ns string) error {
	namespace := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{Name: ns},
		Spec: v1.NamespaceSpec{
			Finalizers: []v1.FinalizerName{
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

func (n *Namespace) HarborCreate(namespace string) error {
	httpclient.WithHeader("Authorization", "Basic YWRtaW46SGFyYm9yMTIzNDU=")
	params := &HarborNamespace{
		ProjectName: namespace,
		Metadata: Metadata{
			Public: "true",
		},
	}
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	res, err := httpclient.PostJson("http://harbor.cloud2go.cn/api/projects", b)
	if err != nil {
		return err
	}
	code := res.StatusCode
	log.Infof("create harbor return code : %v", code)
	return nil
}
