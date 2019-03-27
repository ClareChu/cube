package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	build := new(mocks.BuildAggregate)
	appInfo := newBuildController(build)
	build.On("Create", &v1alpha1.BuildConfig{}, "", "v1").Return(nil, nil)
	b := &v1alpha1.Build{ObjectMeta: metav1.ObjectMeta{
		Name:      "a",
		Namespace: "s",
	}}
	build.On("DeleteNode", b).Return(nil)
	testApplication := web.NewTestApp(t, appInfo).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	t.Run("should pass with jwt token", func(t *testing.T) {
		testApplication.Post("/build").WithJSON(&command.BuildConfig{}).Expect().Status(http.StatusOK)
	})

	t.Run("should pass with jwt token", func(t *testing.T) {
		testApplication.Get("/build/name/a/namespace/s").Expect().Status(http.StatusOK)
	})
}
