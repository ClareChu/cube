package cr

import (
	"hidevops.io/cube/operator/client"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineConfig struct {
	clientSet versioned.Interface
	Resource  string
}

const (
	PipelineConfigResource = "pipelineconfigs"
	IDEName                = "ide"
	SonarName              = "sonar"
	IDEPATHName            = "ide-path"
	Namespace              = "hidevopsio"
	//Namespace              = "cloudtogo-system"
)

func NewPipelineConfig(clientSet versioned.Interface) CubeManagerInterface {
	return &PipelineConfig{
		clientSet: clientSet,
		Resource:  PipelineConfigResource,
	}
}

func (ap *PipelineConfig) create() {
	ide := &v1alpha1.PipelineConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      IDEName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.PipelineSpec{
			Branch:  "master",
			Profile: "dev",
			Version: "v1",
			Events: []v1alpha1.Events{
				{
					EventTypes: "imageStream",
					Name:       IDEName,
				},
				/*		{
						EventTypes: "volume",
						Name:       IDEName,
					},*/
				{
					EventTypes: "deploy",
					Name:       IDEName,
				},
				{
					EventTypes: "service",
					Name:       IDEName,
				},
				{
					EventTypes: "gateway",
					Name:       IDEName,
				},
				{
					EventTypes: "callback",
					Name:       IDEName,
				},
			},
		},
	}
	// create ide pipeline config 流水线
	//client.CreateRestClient(ap.clientSet, ide, ap.Resource)
	client.CreatPipelineConfig(ap.clientSet, ide)
	idePath := &v1alpha1.PipelineConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      IDEPATHName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.PipelineSpec{
			Branch:  "master",
			Profile: "dev",
			Version: "v1",
			Events: []v1alpha1.Events{
				{
					EventTypes: "imageStream",
					Name:       IDEName,
				},
				{
					EventTypes: "volume",
					Name:       IDEName,
				},
				{
					EventTypes: "deploy",
					Name:       IDEName,
				},
				{
					EventTypes: "service",
					Name:       IDEName,
				},
				{
					EventTypes: "gateway",
					Name:       IDEPATHName,
				},
				{
					EventTypes: "callback",
					Name:       IDEName,
				},
			},
		},
	}
	// create ide path pipeline config 流水线
	//client.CreateRestClient(ap.clientSet, idePath, ap.Resource)
	client.CreatPipelineConfig(ap.clientSet, idePath)

	sonar := &v1alpha1.PipelineConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      SonarName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.PipelineSpec{
			Profile: "dev",
			Version: "v1",
			Events: []v1alpha1.Events{
				{
					EventTypes: "imageStream",
					Name:       IDEName,
				},
				{
					EventTypes: "deploy",
					Name:       "sonar",
				},
				{
					EventTypes: "service",
					Name:       "sonar",
				},
				{
					EventTypes: "callback",
					Name:       IDEName,
				},
			},
		},
	}
	// create sonar pipeline config 流水线
	//client.CreateRestClient(ap.clientSet, sonar, ap.Resource)
	client.CreatPipelineConfig(ap.clientSet, sonar)
}
