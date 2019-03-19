package service

import (
	"github.com/magiconair/properties/assert"
	gg "github.com/xanzy/go-gitlab"
	gogitlab "github.com/xanzy/go-gitlab"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/hioak/starter/scm/gitlab"
	"hidevops.io/hioak/starter/scm/gitlab/fake"
	"hidevops.io/mio/console/pkg/constant"
	fk "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestLoginServiceImplGetSession(t *testing.T) {
	fs := new(fake.SessionService)
	us := new(fake.UsersService)
	cli := &fake.Client{
		SessionService: fs,
		UsersService:   us,
	}
	cli1 := &fake.Client{
		UsersService: us,
	}
	s := gitlab.NewSession(func(url, token string) (client gitlab.ClientInterface) {
		return cli
	})

	u := gitlab.NewUser(func(url, token string) (client gitlab.ClientInterface) {
		return cli1
	})

	user := &gogitlab.User{
		Name: "chulei",
	}
	gs := &gg.Session{
		Username: "chulei",
	}
	gr := new(gg.Response)

	resp := new(gogitlab.Response)
	fs.On("GetSession", nil, nil).Return(gs, gr, nil)
	us.On("CurrentUser", nil).Return(user, resp, nil)
	clientSet := fk.NewSimpleClientset()
	configMaps := kube.NewConfigMaps(clientSet)
	data := map[string]string{
		"API_VERSION":"api/v3",
		"BASE_URL":"https://github.com",
	}
	configMaps.Create(constant.GitlabConstant, constant.TemplateDefaultNamespace, data)
	login := newLoginService(s, u, configMaps)
	_, err := login.GetUser("")
	assert.Equal(t, nil, err)

	_, _, _, err = login.GetSession("", "")
	assert.Equal(t, nil, err)
}
