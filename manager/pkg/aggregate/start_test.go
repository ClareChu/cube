package aggregate

import (
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"gopkg.in/src-d/enry.v1"
	"gotest.tools/assert"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/crypto/base64"
	"testing"
)

type FileName struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

func TestGetLanguagesByFilename(t *testing.T) {
	res, err := httpclient.Get("https://api.github.com/repos/ClareChu/harbor/contents/make/install.sh")
	assert.Assert(t, err, nil)
	bodyBytes, err := res.ReadAll()
	file := &FileName{}
	json.Unmarshal(bodyBytes, file)
	_, err = base64.Decode([]byte(file.Content))
	lang := enry.GetLanguage("", nil)
	log.Info(lang)
}