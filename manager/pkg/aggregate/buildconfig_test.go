package aggregate

/*
func TestBuildConfigTemplate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	configMaps := new(mocks.ConfigMapsAggregate)
	buildConfig := cube.NewBuildConfig(clientSet)
	buildAggregate := new(mocks.BuildAggregate)
	//buildConfigAggregate := NewBuildConfigService(buildConfig, buildAggregate, configMaps)
	bc := &command.BuildConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: "demo",
		},
	}
	build1 := &v1alpha1.BuildConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.BuildConfigKind,
			APIVersion: constant.BuildConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
			Labels: map[string]string{
				"CodeType": "java",
			},
		},
		Spec: v1alpha1.BuildSpec{
			Tags: []string{"/demo/hello-world"},
			App:  "hello-world",
		},
		Status: v1alpha1.BuildConfigStatus{
			LastVersion: 1,
		},
	}
	_, err := buildConfigAggregate.Template(bc)
	assert.Equal(t, nil, err)
	config := &corev1.ConfigMap{
		Data: map[string]string{
			constant.Username : "",
			constant.Password: "",
			constant.DockerRegistry: "",
		},
	}
	configMaps.On("Get", "docker", "hidevopsio").Return(config, nil)
	buildAggregate.On("Create", build1, "hello-world-1", "v1").Return(nil, nil)
	param := &command.PipelineReqParams{}
	_, err = buildConfigAggregate.Create(param)
}
*/
