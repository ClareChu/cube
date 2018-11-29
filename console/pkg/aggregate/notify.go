package aggregate

import (
	"fmt"
	"github.com/iris-contrib/go.uuid"
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NotifyAggregate interface {
	Create(notify *v1alpha1.Notify) (*v1alpha1.Notify, error)
}

type Notify struct {
	NotifyAggregate
	notify *mio.Notify
}

func init() {
	app.Register(NewNotifyService)
}

func NewNotifyService(notify *mio.Notify) NotifyAggregate {
	return &Notify{
		notify: notify,
	}
}

func (n *Notify) Create(notify *v1alpha1.Notify) (*v1alpha1.Notify, error) {
	log.Infof("create notify :%v", notify)
	return n.notify.Create(notify)
}

func (n *Notify) CreateNotify(name, namespace, profile, sourceCode string) (notify *v1alpha1.Notify, err error) {
	log.Info("create notify ")
	uid, err := uuid.NewV4()
	name = fmt.Sprintf("%s-%s", name, uid)
	notifyTemplate, err := n.notify.Get(sourceCode, constant.TemplateDefaultNamespace)
	log.Infof("create templates :%v", notifyTemplate)
	if err != nil {
		return
	}
	notify = new(v1alpha1.Notify)
	notify.TypeMeta = v1.TypeMeta{
		Kind:       constant.NotifyKind,
		APIVersion: constant.NotifyApiVersion,
	}
	notify.ObjectMeta = v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			constant.App:  name,
			constant.Name: name,
		},
	}
	notify.Spec = notifyTemplate.Spec
	config, err := n.notify.Create(notify)
	log.Infof("create notify config : %v", config)
	if err != nil {
		return
	}
	//notify.Spec.Roles
	//builder.SendEMail()
	return
}
