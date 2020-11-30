package aggregate

import (
	"encoding/json"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"hidevops.io/cube/manager/pkg/aggregate/client"
	nss "hidevops.io/cube/manager/pkg/aggregate/ns"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/crypto/base64"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
)

type NamespaceAggregate interface {
	InitNamespace(namespace string, tls bool) error
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
	secret     *service.Secret
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

func (n *Namespace) InitNamespace(namespace string, tls bool) (err error) {
	if os.Getenv("OCP") == "openshift" {
		err := n.CreateOcp(namespace)
		if err != nil {
			return nil
		}
	} else {
		err := n.Create(namespace)
		if err != nil {
			return nil
		}
	}
	if tls {
		err = n.secret.CreateTraefikSecret(namespace)
		if err != nil {
			return err
		}
	}
	err = n.HarborCreate(namespace)
	return err
}

const HarborCreateNamespaceApi = "api/projects"

var (
	// 支持openshift 使pod uid发生改变
	NamespacesAnnotation = map[string]string{
		"openshift.io/sa.scc.uid-range": "0/0",
	}
	AnnotationKey   = "openshift.io/sa.scc.uid-range"
	AnnotationValue = "0/0"
)

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
			log.Infof("create namespace err :%v", err)
		}
	}
	return err
}

//todo Create 暂时支持openshift 做的改变 巨坑
func (n *Namespace) CreateOcp(ns string) error {
	namespace := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: ns,
		},
		Spec: v1.NamespaceSpec{
			Finalizers: []v1.FinalizerName{
			},
		},
	}
	options := meta_v1.GetOptions{}
	_, err := n.namespace.Get(ns, options)
	if err != nil {
		clientSet, err := client.GetDefaultK8sClientSet()
		newNamespace := nss.NewNamespace(clientSet)
		_, err = newNamespace.Create(namespace)
		if err != nil {
			log.Errorf("create namespace err :%v", err)
		}
		time.Sleep(5 * time.Second)
		name, err := newNamespace.Get(ns)
		if err != nil {
			return err
		}
		name.ObjectMeta.Annotations["openshift.io/sa.scc.uid-range"] = "0/0"
		err = newNamespace.Update(name)
		return err
	}
	time.Sleep(7 * time.Second)
	return err
}

func (n *Namespace) HarborCreate(namespace string) error {
	config, err := n.configMaps.Get(constant.DockerConstant, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("docker config map not found :%v", err)
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
	go func() {
		res, err := httpclient.PostJson(fmt.Sprintf("http://%s/%s", config.Data[constant.DockerRegistry], HarborCreateNamespaceApi), b)
		if err != nil {
			return
		}
		code := res.StatusCode
		log.Infof("create harbor return code : %v", code)
	}()
	return nil
}
