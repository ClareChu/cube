package service

import (
	"fmt"
	docker_types "github.com/docker/docker/api/types"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"hidevops.io/cube/agent/protobuf"
	"hidevops.io/cube/agent/types"
	pkg_utils "hidevops.io/cube/agent/utils"
	cubev1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/docker"
	scmgit "hidevops.io/hioak/starter/scm/git"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"time"
)

type BuildConfigService interface {
	Clone(sourceCodePullRequest *protobuf.SourceCodePullRequest, cloneFunc scmgit.CloneFunc) (string, error)
	Compile(compileRequest *protobuf.CompileRequest) error
	ImageBuild(imageBuildRequest *protobuf.ImageBuildRequest) error
	ImagePush(imagePushRequest *protobuf.ImagePushRequest) error
	GetImage(imagePushRequest *protobuf.ImagePushRequest) error
	CreateImage(name, namespace, tag string, imageSummary docker_types.ImageSummary) error
}

type buildConfigServiceImpl struct {
	BuildConfigService
	imageClient *docker.ImageClient
	imageStream *cube.ImageStream
}

func init() {
	log.SetLevel(log.DebugLevel)
	app.Register(newBuildService)
}

func newBuildService(imageClient *docker.ImageClient, imageStream *cube.ImageStream) BuildConfigService {
	return &buildConfigServiceImpl{
		imageClient: imageClient,
		imageStream: imageStream,
	}
}

const (
	Latest     = "latest"
	Generation = "1"
)

func (b *buildConfigServiceImpl) Clone(sourceCodePullRequest *protobuf.SourceCodePullRequest, cloneFunc scmgit.CloneFunc) (string, error) {
	log.Infof("git clone url: %s, branch: %s", sourceCodePullRequest.Url, sourceCodePullRequest.Branch)
	if sourceCodePullRequest.Token != "" {
		//CMD
		Path, err := pkg_utils.CloneBYCMD(sourceCodePullRequest)
		if err != nil {
			return "", err
		} else {
			return Path, nil
		}
	}

	//go-git
	passwordAuth := transport.AuthMethod(&http.BasicAuth{
		Username: sourceCodePullRequest.Username,
		Password: sourceCodePullRequest.Password},
	)
	//	tokenAuth := transport.AuthMethod(&http.TokenAuth{Token: sourceCodePullRequest.Password})

	referenceName := fmt.Sprintf("refs/heads/%s", sourceCodePullRequest.Branch)
	codePath, err := scmgit.NewRepository(cloneFunc).Clone(&git.CloneOptions{URL: sourceCodePullRequest.Url,
		ReferenceName:     plumbing.ReferenceName(referenceName),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Depth:             int(sourceCodePullRequest.Depth),
		Auth:              passwordAuth,
	},
		sourceCodePullRequest.DstDir)

	if err != nil {
		log.Infof("clone %s filed:", sourceCodePullRequest.Url)
		os.RemoveAll(codePath)
		return "", err
	}
	log.Infof("clone %s succeed", sourceCodePullRequest.Url)
	return codePath, nil
}

func (b *buildConfigServiceImpl) Compile(compileRequest *protobuf.CompileRequest) error {
	log.Infof("compile start")

	execCommand := func(CommandName string, Params []string) error {
		cmd, bufioReader, err := pkg_utils.ExecCommand(CommandName, Params)
		if err != nil {

			return err
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Start(); err != nil {
			return err
		}
		for {
			line, err2 := bufioReader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}

			log.Infof(line)
		}
		if err = cmd.Wait(); err != nil {
			return err
		}
		return nil
	}
	codeType := os.Getenv("CODE_TYPE")
	if codeType == "" {
		return fmt.Errorf("env CODE_TYPE get filed")
	}

	if types.JAVA == codeType {

		pomXmlInfo, err := pkg_utils.GetPomXmlInfo("pom.xml")
		if err != nil {
			return err
		}

		projectName := fmt.Sprintf("%s-%s.%s", pomXmlInfo.ArtifactId, pomXmlInfo.Version, pomXmlInfo.Packaging)

		log.Infof("project name %v", projectName)
		compileRequest.CompileCmd = append(compileRequest.CompileCmd, &protobuf.BuildCommand{ExecType: string(string(cubev1alpha1.Script)),
			Script: fmt.Sprintf("cp target/%s app.%s", projectName, pomXmlInfo.Packaging),
		})

	}

	for _, cmd := range compileRequest.CompileCmd {
		if cmd.ExecType == string(cubev1alpha1.Script) {
			scriptPath, err := pkg_utils.GenScript(cmd.Script)
			if err != nil {
				return err
			}

			if err := execCommand("chmod", []string{"+x", scriptPath}); err != nil {
				return err
			}

			if err := execCommand("sh", []string{"-c", scriptPath}); err != nil {
				log.Errorf("Error compile filed err: %v", err)
				return err
			}
			os.RemoveAll(scriptPath)
			continue
		}

		if err := execCommand(cmd.CommandName, cmd.Params); err != nil {
			log.Errorf("compile filed err : %v", err)
			return err
		}
	}

	fmt.Print("\n[INFO] compile succeed\n")
	return nil
}

func (b *buildConfigServiceImpl) ImageBuild(imageBuildRequest *protobuf.ImageBuildRequest) error {
	fmt.Printf("\n[INFO] image %v start build:\n", imageBuildRequest.Tags)

	buildImage := &docker.Image{
		Tags:       imageBuildRequest.Tags,
		BuildFiles: pkg_utils.GetBuildFileBYDockerfile(imageBuildRequest.DockerFile),
		Username:   imageBuildRequest.Username,
		Password:   imageBuildRequest.Password,
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte(strings.Join(imageBuildRequest.DockerFile, "\n"))); err != nil {
		return err
	}
	file.Close()
	//defer os.RemoveAll("Dockerfile")

	imageBuildResponse, err := b.imageClient.BuildImage(buildImage)
	if err != nil {
		fmt.Printf("\nError image %v build filed\n", imageBuildRequest.Tags)
		return err
	}
	defer imageBuildResponse.Body.Close()

	if _, err := io.Copy(os.Stdout, imageBuildResponse.Body); err != nil {
		return err
	}

	return nil
}

func (b *buildConfigServiceImpl) ImagePush(imagePushRequest *protobuf.ImagePushRequest) error {

	fmt.Printf("\n[INFO] image %v start push:\n", imagePushRequest)
	fmt.Printf("xxxxxxxxxxx")
	for _, imageName := range imagePushRequest.Tags {
		imageInfo := strings.Split(imageName, ":")
		pushImage := &docker.Image{
			Username:  imagePushRequest.Username,
			Password:  imagePushRequest.Password,
			FromImage: imageInfo[0],
			Tag:       imageInfo[1],
		}

		if err := b.imageClient.PushImage(pushImage); err != nil {
			fmt.Printf("\nError image %s push filed\n", imageName)
			return err
		}

	}
	log.Info("push image success ")
	fmt.Printf("xxxxxxxxxxx")
	err := b.GetImage(imagePushRequest)
	return err
}

func (b *buildConfigServiceImpl) GetImage(imagePushRequest *protobuf.ImagePushRequest) error {
	log.Infof("get image: %v", imagePushRequest.ImageName)
	fmt.Printf("xxxxxxxxxxx")
	imageInfo := strings.Split(imagePushRequest.Tags[0], ":")
	image := &docker.Image{
		FromImage: imageInfo[0],
		Tag:       imageInfo[1],
	}
	imageSummary, err := b.imageClient.GetImage(image)
	if err != nil {
		log.Error("get image is not found ")
		return err
	}
	b.CreateImage(imagePushRequest.ImageName, imagePushRequest.Namespace, imageInfo[1], imageSummary)
	fmt.Printf("xxxxxxxxxxx")
	return nil
}

func (b *buildConfigServiceImpl) CreateImage(name, namespace, tag string, imageSummary docker_types.ImageSummary) error {
	t := time.Now()
	image, err := b.imageStream.Get(name, namespace)
	stream := &cubev1alpha1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
	}
	spec := cubev1alpha1.ImageStreamSpec{
		DockerImageRepository: imageSummary.RepoTags[0],
		Tags: map[string]cubev1alpha1.Tag{
			tag: cubev1alpha1.Tag{
				Created:              t.UTC().Format(time.UnixDate),
				DockerImageReference: imageSummary.RepoDigests[0],
				Generation:           "1",
				Image:                strings.Split(imageSummary.RepoDigests[0], "@")[1],
			},
			Latest: cubev1alpha1.Tag{
				Created:              t.UTC().Format(time.UnixDate),
				DockerImageReference: imageSummary.RepoDigests[0],
				Generation:           "1",
				Image:                strings.Split(imageSummary.RepoDigests[0], "@")[1],
			},
		},
	}
	stream.Spec = spec
	if err != nil {
		log.Errorf("get image: %v", err)
		_, err := b.imageStream.Create(stream)
		return err
	}
	delete(image.Spec.Tags, tag)
	delete(image.Spec.Tags, Latest)
	image.Spec.Tags[tag] = cubev1alpha1.Tag{
		Created:              t.UTC().Format(time.UnixDate),
		DockerImageReference: imageSummary.RepoDigests[0],
		Generation:           Generation,
		Image:                strings.Split(imageSummary.RepoDigests[0], "@")[1],
	}
	image.Spec.Tags[Latest] = cubev1alpha1.Tag{
		Created:              t.UTC().Format(time.UnixDate),
		DockerImageReference: imageSummary.RepoDigests[0],
		Generation:           Generation,
		Image:                strings.Split(imageSummary.RepoDigests[0], "@")[1],
	}
	_, err = b.imageStream.Update(name, namespace, image)
	return err
}
