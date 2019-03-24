package service

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	docker_types "github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"hidevops.io/cube/node/protobuf"
	"hidevops.io/cube/node/service/mock"
	cube_fake "hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/docker"
	"hidevops.io/hioak/starter/docker/fake"
	"io"
	"os"
	"strings"
	"testing"
)

//go:generate mockgen -destination mock/mock_build.go -package mock hidevops.io/cube/node/pkg/service BuildConfigService

func TestBuild(t *testing.T) {

	defer os.RemoveAll("Dockerfile")

	projectName := "demo"

	var cmdList []*protobuf.BuildCommand
	cmd1 := &protobuf.BuildCommand{
		CodeType:    "",
		CommandName: "pwds",
		Params:      []string{},
	}

	cmd2 := &protobuf.BuildCommand{
		ExecType: "script",
		Script: `if [[ $? == 0 ]]; then
          echo "Build Successful."
        else
          echo "Build Failed!"
          exit 1
        fi`,
	}

	cmdList = append(cmdList, cmd1)
	cmdList = append(cmdList, cmd2)

	compileRequest := protobuf.CompileRequest{
		CompileCmd: cmdList,
	}

	//
	imageBuildRequest := protobuf.ImageBuildRequest{
		App:        projectName,
		S2IImage:   "FROM ubuntu:16.04",
		Tags:       []string{"test:0.1"},
		DockerFile: []string{},
	}

	//
	imagePushRequest := protobuf.ImagePushRequest{
		Tags: []string{"test:0.1"},
	}

	cli, err := fake.NewClient()
	assert.Equal(t, nil, err)
	buildConfigimpl := buildConfigServiceImpl{
		imageClient: &docker.ImageClient{Client: cli},
	}

	mockCtl := gomock.NewController(t)
	m := mock.NewMockBuildConfigService(mockCtl)

	t.Run("should err is nil in the clone code", func(t *testing.T) {
		m.EXPECT().Clone(nil).Return("", nil)
		dstPath, err := m.Clone(nil)
		assert.Equal(t, "", dstPath)
		assert.Equal(t, nil, err)
	})

	t.Run("should err in the compile", func(t *testing.T) {
		err := buildConfigimpl.Compile(&compileRequest)
		assert.NotEqual(t, nil, err)

	})

	t.Run("should err not nil in the image build", func(t *testing.T) {

		cli.On("ImageBuild", nil, nil,
			nil).Return(types.ImageBuildResponse{}, errors.New("1"))
		err = buildConfigimpl.ImageBuild(&imageBuildRequest)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should err not nil in the image push", func(t *testing.T) {

		var i io.ReadCloser
		cli.On("ImagePush", nil, nil, nil).Return(i, errors.New("1"))
		err = buildConfigimpl.ImagePush(&imagePushRequest)
		assert.NotEqual(t, nil, err)
	})
}

func TestBuildConfigServiceImpl_Compile(t *testing.T) {

	os.Setenv("CODE_TYPE", ".")

	b := new(buildConfigServiceImpl)
	cr := &protobuf.CompileRequest{

		CompileCmd: []*protobuf.BuildCommand{{CodeType: Command, CommandName: "pwd"}, {ExecType: "script", Script: "pwd"}},
	}

	err := b.Compile(cr)
	assert.Equal(t, nil, err)
}

func TestBuildConfigServiceImpl_Clone(t *testing.T) {
	b := new(buildConfigServiceImpl)
	sp := &protobuf.SourceCodePullRequest{}

	b.Clone(sp, func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {

		return nil, nil
	})
}

func TestBuildConfigServiceImpl_Clone2(t *testing.T) {

	b := new(buildConfigServiceImpl)
	sourceCodePullRequest := &protobuf.SourceCodePullRequest{
		Url:      "http://gitlab.vpclub:8022/wanglulu/clone-test01.git",
		Branch:   "master",
		DstDir:   "/Users/mac/.gvm/pkgsets/go1.10/vpcloud/src/hidevops.io/",
		Username: "",
		Password: "",
		Token:    "test",
	}
	_, err := b.Clone(sourceCodePullRequest, func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
		return nil, nil
	})
	assert.NotEqual(t, nil, err)
}

func TestBuildConfigServiceImplCreateImage(t *testing.T) {
	clientSet := cube_fake.NewSimpleClientset().CubeV1alpha1()
	name := "hello-world"
	namespace := "demo"
	tag := "v1"
	image := cube.NewImageStream(clientSet)
	b := newBuildService(nil, image)
	imageSummary := docker_types.ImageSummary{
		RepoDigests: []string{
			"aaaaaaa@aaa", "",
		},
		RepoTags: []string{
			"kkkkk",
		},
	}
	err := b.CreateImage(name, namespace, tag, imageSummary)
	assert.Equal(t, nil, err)
	err = b.CreateImage(name, namespace, tag, imageSummary)
	assert.Equal(t, nil, err)
	s := strings.Split("docker-registry-default.app.vpclub.io/demo/hello-world@sha256:485a3c93699c107dbe6d8a265a75d282b0bc767b5780f60d45a18a923689cee2", "@")
	fmt.Sprint(s)
}
