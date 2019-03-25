package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/client-go/kubernetes/fake"
	"net/http"
	"testing"
)

func TestAppInfoControllerGetByNamespaceLabel(t *testing.T) {
	client := fake.NewSimpleClientset()
	pod := kube.NewPod(client)
	appInfo := newAppInfoController(pod)
	testApplication := web.NewTestApp(t, appInfo).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		testApplication.Get("/appInfo/namespace/a/label/s").Expect().Status(http.StatusOK)
	})
}
