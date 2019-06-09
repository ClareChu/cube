package cube

import (
	"fmt"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	cubev1alpha1 "hidevops.io/cube/pkg/client/clientset/versioned/typed/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type App struct {
	clientSet cubev1alpha1.CubeV1alpha1Interface
}

func NewApp(clientSet cubev1alpha1.CubeV1alpha1Interface) *App {
	return &App{
		clientSet: clientSet,
	}
}

func (a *App) Create(app *v1alpha1.App) (config *v1alpha1.App, err error) {
	log.Debugf("app create : %v", app.Name)
	config, err = a.clientSet.Apps(app.Namespace).Create(app)
	if err != nil {
		return nil, err
	}
	return
}

func (a *App) Get(name, namespace string) (config *v1alpha1.App, err error) {
	log.Info(fmt.Sprintf("get app %s in namespace %s", name, namespace))
	result, err := a.clientSet.Apps(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *App) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete app %s in namespace %s", name, namespace))
	err := a.clientSet.Apps(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (a *App) Update(name, namespace string, config *v1alpha1.App) (*v1alpha1.App, error) {
	log.Info(fmt.Sprintf("update app %s in namespace %s", name, namespace))
	result, err := a.clientSet.Apps(namespace).Update(config)
	return result, err
}

func (a *App) List(namespace string, option v1.ListOptions) (*v1alpha1.AppList, error) {
	log.Info(fmt.Sprintf("list app in namespace %s", namespace))
	result, err := a.clientSet.Apps(namespace).List(option)
	return result, err
}

func (a *App) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s app，in the namespace %s", listOptions.LabelSelector, namespace))

	listOptions.Watch = true

	w, err := a.clientSet.Apps(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
