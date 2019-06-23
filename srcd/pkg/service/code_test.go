package service

import (
	"gopkg.in/src-d/go-git.v4"
	"gotest.tools/assert"
	"hidevops.io/cube/srcd/pkg/entity"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestCodeClone(t *testing.T) {
	code := new(Code)
	c := &entity.Clone{
		Url:       "https://github.com/ClareChu/51well-video",
		Username:  "ClareChu",
		Password:  "lei13971368720",
		Namespace: "ClareChu",
		Name:      "51well-video",
	}
	g := git.PlainClone
	_, err := code.Clone(c, g)
	assert.Equal(t, err, nil)
}


func TestCheck(t *testing.T) {
	code := new(Code)
	l, err := code.Check("enry", nil)
	assert.Equal(t, err, nil)
	m := ToMap(l)
	log.Info(m)
}