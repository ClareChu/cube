package service

import (
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/hioak/starter/scm"
	"hidevops.io/hioak/starter/scm/gitlab"
	"hidevops.io/mio/console/pkg/constant"
)

type LoginService interface {
	GetSession(username, password string) (string, int, string, error)
	GetUser(accessToken string) (*scm.User, error)
	GetUrl() (baseUrl string, err error)
}

type LoginServiceImpl struct {
	session    *gitlab.Session
	user       *gitlab.User
	configMaps *kube.ConfigMaps
}

func init() {
	app.Register(newLoginService)
}

func newLoginService(session *gitlab.Session, user *gitlab.User, configMaps *kube.ConfigMaps) LoginService {
	return &LoginServiceImpl{
		session:    session,
		user:       user,
		configMaps: configMaps,
	}
}

func (l *LoginServiceImpl) GetSession(username, password string) (string, int, string, error) {
	log.Debugf("get session username: %v", username)
	baseUrl, err := l.GetUrl()
	if err != nil {
		return "", 0, "", err
	}
	err = l.session.GetSession(baseUrl, username, password)
	return l.session.GetToken(), l.session.GetId(), "", err
}

func (l *LoginServiceImpl) GetUser(accessToken string) (*scm.User, error) {
	baseUrl, err := l.GetUrl()
	if err != nil {
		return nil, err
	}
	user, err := l.user.GetUser(baseUrl, accessToken)
	return user, err
}

func (l *LoginServiceImpl) GetUrl() (baseUrl string, err error) {
	configMaps, err := l.configMaps.Get(constant.GitlabConstant, constant.TemplateDefaultNamespace)
	if err != nil {
		return "", err
	}
	apiVersion := configMaps.Data[constant.ApiVersion]
	baseUrl = configMaps.Data[constant.BaseUrl]
	baseUrl = baseUrl + apiVersion
	return
}
