package utils

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/agent/protobuf"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestTestStart(t *testing.T) {
	sourceCodeTestRequest := &protobuf.CommandRequest{
		CommandList: []*protobuf.Command{{
			CommandName: "pwd",
			Params:      []string{},
		}, {ExecType: "script", Script: "pwd"}},
	}
	err := StartCmd(sourceCodeTestRequest)

	assert.Equal(t, nil, err)
}

func TestStartCmd(t *testing.T) {
	sourceCodeTestRequest := &protobuf.TestsRequest{
		TestCmd: []*protobuf.TestCommand{{
			CommandName: "pwd",
			Params:      []string{},
		}, {ExecType: "script", Script: "pwd"}},
	}
	err := TestStart(sourceCodeTestRequest)

	assert.Equal(t, nil, err)
}

func TestGetCurrentDirectory(t *testing.T) {
	dir := GetCurrentDirectory()
	log.Info(dir)
}