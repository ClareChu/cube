package service

import (
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/scm"
	"hidevops.io/hioak/starter/scm/gitlab"
	"os"
)

var apiVersion = "/api/v3"

type LoginService interface {
	GetSession(baseUrl, username, password string) (string, int, string, error)
	GetUser(baseUrl, accessToken string) (*scm.User, error)
}

type LoginServiceImpl struct {
	session *gitlab.Session
	user    *gitlab.User
}

func init() {
	app.Register(newLoginService)
	av := os.Getenv("API_VERSION")
	if av != "" {
		apiVersion = av
	}
}

func newLoginService(session *gitlab.Session, user *gitlab.User) LoginService {
	return &LoginServiceImpl{
		session: session,
		user:    user,
	}
}

func (l *LoginServiceImpl) GetSession(baseUrl, username, password string) (string, int, string, error) {
	log.Debugf("get session baseUrl: %v , username: %v, password:%v", baseUrl, username, password)
	baseUrl = baseUrl + apiVersion
	err := l.session.GetSession(baseUrl, username, password)
	return l.session.GetToken(), l.session.GetId(), "", err
}

func (l *LoginServiceImpl) GetUser(baseUrl, accessToken string) (*scm.User, error) {
	baseUrl = baseUrl + apiVersion
	user, err := l.user.GetUser(baseUrl, accessToken)
	return user, err
}
