package service

import (
	"bufio"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"hidevops.io/cube/srcd/pkg/entity"
	"hidevops.io/cube/srcd/pkg/utils"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	scmgit "hidevops.io/hioak/starter/scm/git"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CodeService interface {
	Get() (err error)
	Clone(clone *entity.Clone, cloneFunc scmgit.CloneFunc) (codePath string, err error)
	Pull() (err error)
	CheckWorkspace() (err error)
	Check(commandName string, param []string) (lines []string, err error)
}

type Code struct {
	CodeService
}

func init() {
	app.Register(NewCodeService)
}

func NewCodeService() CodeService {
	return &Code{}
}

func (c *Code) Get() (err error) {
	return nil
}

func (c *Code) Clone(clone *entity.Clone, cloneFunc scmgit.CloneFunc) (codePath string, err error) {
	log.Infof("git clone url: %s, branch: %s", clone.Url, clone.Branch)
	if clone.Token != "" {
		//CMD
		codePath, err := utils.CloneBYCMD(clone)
		if err != nil {
			return "", err
		} else {
			return codePath, nil
		}
	}

	//go-git
	passwordAuth := transport.AuthMethod(&http.BasicAuth{
		Username: clone.Username,
		Password: clone.Password},
	)
	codePath, err = scmgit.NewRepository(cloneFunc).Clone(&git.CloneOptions{URL: clone.Url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              passwordAuth,
	},
		clone.DstDir)

	if err != nil {
		log.Infof("clone %s filed:", clone.Url)
		os.RemoveAll(codePath)
		return "", err
	}
	log.Infof("clone %s succeed", clone.Url)
	return codePath, nil
}

func (c *Code) Pull() (err error) {
	return nil
}

func (c *Code) Check(commandName string, params []string) (lines []string, err error) {
	cmd := exec.Command(commandName, params...)

	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		lines = append(lines, line)
	}

	cmd.Wait()
	return lines, nil
}

func ToMap(lines []string) (l map[string]string) {
	l = map[string]string{}
	for _, line := range lines {
		s := strings.Split(line, "\t")
		key := s[0]
		value := strings.Split(s[1], "\n")[0]
		l[key] = value
	}
	return l
}
