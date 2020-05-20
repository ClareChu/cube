package kube

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"hidevops.io/cube/manager/pkg/aggregate/client"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewNamespace(t *testing.T) {
	ns := "fp2005130205477160000002542269"
	clientSet, err := client.GetDefaultK8sClientSet()
	assert.Equal(t, err, nil)
	n := NewNamespace(clientSet)
	namespace := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: ns,
		},
		Spec: v1.NamespaceSpec{
			Finalizers: []v1.FinalizerName{
			},
		},
	}
	_, err = n.Create(namespace)
	assert.Equal(t, err, nil)
	nss, err := n.Get(ns)
	out, err := yaml.Marshal(nss)
	fmt.Printf("--- m dump: %s", string(out))
	nss.ObjectMeta.Annotations["openshift.io/sa.scc.uid-range"] = "0/0"
	err = n.Update(nss)
	assert.Equal(t, err, nil)
}


func TestNewNamespace1(t *testing.T) {
	ns := "fp2005130205477160000002542269"
	clientSet, err := client.GetDefaultK8sClientSet()
	assert.Equal(t, err, nil)
	n := NewNamespace(clientSet)
	nss, err := n.Get(ns)
	out, err := yaml.Marshal(nss)
	fmt.Printf("--- m dump: %s", string(out))
	nss.ObjectMeta.Annotations["openshift.io/sa.scc.uid-range"] = "0/0"
	assert.Equal(t, err, nil)
}
