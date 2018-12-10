package service

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/utils/replacer"
	"os"
	"reflect"
	"strings"
)

const (
	ApplicationId  = "a824457ddf48f93c4fc1315ed9ecf4a22576ab54734ce920906c402accc1a704"
	Secret         = "41fcb4021e19d22a7205164c29443acfc10c525aadc150dd878d610f4df91eed"
	CallbackUrl    = "https://devcloud.vpclub.cn/hidevopsio/admin/index.html"
	BaseUrl        = "${SCM_URL}/oauth/authorize?client_id=${ApplicationId}&redirect_uri=${CallbackUrl}&response_type=code"
	AccessTokenUrl = "${SCM_URL}/oauth/token?client_id=${ApplicationId}&redirect_uri=${CallbackUrl}&client_secret=${Secret}&code=${Code}&grant_type=authorization_code"
)


type OauthService interface {
	GetAuthURL(a *Auth) string
}

type OauthServiceImpl struct {

}

func init() {
	app.Register(newOauthService)
}

func newOauthService() OauthService {
	return &OauthServiceImpl{
	}
}

type Auth struct {
	ApplicationId string
	Secret        string
	AuthURL       string
	CallbackUrl   string
}

func (o *OauthServiceImpl) GetAuthURL(a *Auth) string {
	t := replacer.GetMatches(BaseUrl)
	baseUrl := ReplaceEnv(BaseUrl, t, a)
	return baseUrl
}

func ReplaceEnv(source string, rr [][]string, t interface{}) string {
	for _, r := range rr {
		base := os.Getenv(r[1])
		if base == "" {
			immutable := reflect.ValueOf(t).Elem()
			base = immutable.FieldByName(r[1]).String()
		}
		source = strings.Replace(source, r[0], base, -1)
	}
	return source
}

