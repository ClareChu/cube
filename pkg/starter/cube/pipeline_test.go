package cube

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPipelineCurd(t *testing.T) {
	name := "test"
	namespace := "demo-dev"
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	config := NewPipeline(clientSet)
	config1 := &v1alpha1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
	}
	result, err := config.Create(config1)
	assert.Equal(t, nil, err)
	assert.Equal(t, name, result.Name)
	option := v1.ListOptions{

		LabelSelector: "app=" + name,
	}
	list, err := config.List(namespace, option)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(list.Items))

	result, err = config.Get(name, namespace)
	assert.Equal(t, nil, err)
	//assert.Equal(t, name, result.Name)

	con := &v1alpha1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	result, err = config.Update(name, namespace, con)
	assert.Equal(t, nil, err)
	assert.Equal(t, name, result.Name)

	err = config.Delete(name, namespace)
	assert.Equal(t, nil, err)

	listOptions := metav1.ListOptions{}
	i, err := config.Watch(listOptions, namespace)
	log.Infof("i: %v", i)
	assert.Equal(t, nil, err)
}
