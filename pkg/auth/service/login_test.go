package service

import (
	"github.com/magiconair/properties/assert"
	gg "github.com/xanzy/go-gitlab"
	gogitlab "github.com/xanzy/go-gitlab"
	"hidevops.io/hioak/starter/scm/gitlab"
	"hidevops.io/hioak/starter/scm/gitlab/fake"
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
	login := newLoginService(s, u)
	_, err := login.GetUser("", "")
	assert.Equal(t, nil, err)

	_, _, _, err = login.GetSession("", "", "")
	assert.Equal(t, nil, err)
}
