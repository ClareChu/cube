package aggregate

import (
	"encoding/json"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/crypto/base64"
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
	namespace  *kube.Namespace
	configMaps *kube.ConfigMaps
}

func init() {
	app.Register(NewNamespace)
}

func NewNamespace(namespace *kube.Namespace, configMaps *kube.ConfigMaps) NamespaceAggregate {
	return &Namespace{
		namespace:  namespace,
		configMaps: configMaps,
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

const HarborCreateNamespaceApi = "api/projects"

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
		if err != nil {
			log.Info("create namespace err :%v", err)
		}
	}
	return err
}

func (n *Namespace) HarborCreate(namespace string) error {
	config, err := n.configMaps.Get(constant.DockerConstant, constant.TemplateDefaultNamespace)
	if err != nil {
		return err
	}
	auth := fmt.Sprintf("%s:%s", config.Data[constant.Username], config.Data[constant.Password])
	basic := fmt.Sprintf("%s %s", constant.Basic, base64.EncodeToString(auth))
	httpclient.WithHeader(constant.Authorization, basic)
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
	res, err := httpclient.PostJson(fmt.Sprintf("http://%s/%s", config.Data[constant.DockerRegistry], HarborCreateNamespaceApi), b)
	if err != nil {
		return err
	}
	code := res.StatusCode
	log.Infof("create harbor return code : %v", code)
	return nil
}
