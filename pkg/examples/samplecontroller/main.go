package samplecontroller

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	samplecontrollerv1alpha1 "hidevops.io/mio/pkg/examples/samplecontroller/pkg/apis/samplecontroller/v1alpha1"
	clientset "hidevops.io/mio/pkg/examples/samplecontroller/pkg/client/clientset/versioned"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

type sampleController struct {
	restConfig *kube.RestConfig
}

func newSampleController(restConfig *kube.RestConfig) *sampleController {
	return &sampleController{
		restConfig: restConfig,
	}
}

const (
	crdName string = "foos.samplecontroller.k8s.io"
	fooName string = "example-foo"
)

func (c *sampleController) TestCreateCRD() error {

	clientSet, err := apiextensionsclient.NewForConfig(c.restConfig.Config)
	if err != nil {
		log.Debugf("Error building example clientset: %s", err.Error())
		return err
	}

	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: crdName},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "samplecontroller.k8s.io",
			Version: "v1alpha1",
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: "foos",
				Kind:   "Foo",
			},
		},
	}

	if _, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd); err != nil {
		log.Debugf("Error building example clientset: %s", err)
		//return err
	}

	return nil
}

func (c *sampleController) TestDeleteCRD() error {

	clientSet, err := apiextensionsclient.NewForConfig(c.restConfig.Config)
	if err != nil {
		log.Debugf("Error building example clientset: %s", err.Error())
		return err
	}

	if err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crdName, &metav1.DeleteOptions{}); err != nil {
		log.Debugf("Error building example clientset: %s", err)
		return err
	}

	return nil
}

func (c *sampleController) TestFooClient() error {
	//create foo crd
	exampleClient, err := clientset.NewForConfig(c.restConfig.Config)
	if err != nil {
		log.Debugf("Error building example clientset: %s", err)
		return err
	}
	foo := exampleClient.SamplecontrollerV1alpha1().Foos("default")

	// create foo
	replica := int32(1)
	f := &samplecontrollerv1alpha1.Foo{
		ObjectMeta: metav1.ObjectMeta{Name: "example-foo"},
		Spec: samplecontrollerv1alpha1.FooSpec{
			DeploymentName: fooName,
			Replicas:       &replica,
		},
	}

	if _, err := foo.Create(f); err != nil {
		log.Debugf("Error : %v", err)
		//return err
	}

	// list
	fooList, err := foo.List(metav1.ListOptions{})
	if err != nil {
		log.Debugf("Error : %v", err)
		return err
	}
	if len(fooList.Items) == 0 {
		err := fmt.Errorf("list foo faile.")
		log.Debugf("Error : %v", err)
		return err
	}

	// get one
	fooStruct, err := foo.Get("example-foo", metav1.GetOptions{})
	if err != nil {
		log.Debugf("Error : %v", err)
		return err
	}

	//update
	replica = int32(2)
	fooStruct.Spec.Replicas = &replica
	if _, err := foo.Update(fooStruct); err != nil {
		log.Debugf("Error : %v", err)
		return err
	}

	//watch
	watchFoo, err := foo.Watch(metav1.ListOptions{})
	if err != nil {
		log.Debugf("Error : %v", err)
		return err
	}

	if _, ok := <-watchFoo.ResultChan(); !ok {
		err := fmt.Errorf("foo watch faile")
		log.Debugf("Error : %v", err)
		return err
	}

	// delete
	if err := foo.Delete(fooName, &metav1.DeleteOptions{}); err != nil {
		log.Debugf("Error : %v", err)
		return err
	}
	return nil
}

func init() {
	app.Component(newSampleController)
	log.SetLevel(log.DebugLevel)
}

func TestCrd(t *testing.T) {
	testApp := web.NewTestApplication(t).(app.ApplicationContext)

	sampleController := testApp.GetInstance("sampleController").(*sampleController)

	t.Run("|should register crd", func(t *testing.T) {
		err := sampleController.TestCreateCRD()
		assert.Equal(t, nil, err)
	})

	time.Sleep(time.Second * 2)
	t.Run("|should create crd", func(t *testing.T) {
		err := sampleController.TestFooClient()
		assert.Equal(t, nil, err)
	})

	time.Sleep(time.Second * 2)
	t.Run("|should delete crd foo", func(t *testing.T) {
		err := sampleController.TestDeleteCRD()
		assert.Equal(t, nil, err)
	})

}
