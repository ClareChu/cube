package constant

const (
	BuildKind              = "Build"
	BuildApiVersion        = "Build.cube.io/v1alpha1"
	CLONE                  = "clone"
	COMPILE                = "compile"
	BuildImage             = "buildImage"
	PushImage              = "pushImage"
	DeployNode             = "deployNode"
	CreateService          = "createService"
	DeleteDeployment       = "deleteDeployment"
	Complete               = "complete"
	Ending                 = "ending"
	PipelineName           = "pipelineName"
	TimeoutSeconds   int64 = 60 * 60
	DefaultTag             = "latest"
	DefaultPort            = "7575"
)
