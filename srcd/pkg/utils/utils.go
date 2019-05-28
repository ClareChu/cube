package utils

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"hidevops.io/cube/srcd/pkg/entity"
	utilsio "hidevops.io/hiboot/pkg/utils/io"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ResourceString struct {
	XMLName      xml.Name `xml:"project"`
	ModelVersion string   `xml:"modelVersion"`
	GroupId      string   `xml:"groupId"`
	ArtifactId   string   `xml:"artifactId"`
	Version      string   `xml:"version"`
	Packaging    string   `xml:"packaging"`
}

func ExecCommand(commandName string, params []string) (*exec.Cmd, *bufio.Reader, error) {
	cmd := exec.Command(commandName, params...)

	fmt.Println("$ ", cmd.Args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	return cmd, bufio.NewReader(stdout), nil
}

func CloneBYCMD(clone *entity.Clone) (string, error) {

	//git clone -b "分支" --depth=1 xxx.git "指定目录"
	urls := strings.Split(clone.Url, "//")

	projectName := utilsio.Filename(clone.Url)
	projectName = utilsio.Basename(projectName)
	projectPath := filepath.Join(clone.DstDir, projectName)

	if _, err := os.Stat(projectPath); err == nil {
		return "", fmt.Errorf("file %s already exists", projectPath)
	}

	url := clone.Url
	if clone.Token != "" {
		url = fmt.Sprintf("%s//oauth2:%s@%s", urls[0], clone.Token, urls[1])
	}
	Params := []string{"clone", "-b", clone.Branch, "--depth=1", url, projectPath}

	cmd, bufioReader, err := ExecCommand("git", Params)
	if err != nil {
		return "", err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return "", err
	}

	for {
		line, err2 := bufioReader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	if err = cmd.Wait(); err != nil {
		return "", err
	}

	if _, err := os.Stat(projectPath); err != nil {
		fmt.Printf("git clone err: %v", err)
		return "", err
	}

	return projectPath, nil
}

func GetCurrentDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pwd)
	return pwd

}