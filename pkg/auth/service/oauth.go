package service

import (
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/utils/replacer"
	"hidevops.io/hioak/starter/kube"
	"os"
	"reflect"
	"strings"
)

const (
	ApplicationId  = "2c0296a841f08a0cd08f5c40bef100f39b311c019cd2f833d8e7ecdaf61014ef"
	Secret         = "c99a788458ea0d7cdedb1405da8222fd0cc8b7367faa94715d901a1daba3da8c"
	CallbackUrl    = "http://localhost:8081/oauth/code"
	OauthUrl       = "${SCM_URL}/oauth/authorize?client_id=${ApplicationId}&redirect_uri=${CallbackUrl}&response_type=code"
	AccessTokenUrl = "${SCM_URL}/oauth/token?client_id=${ApplicationId}&redirect_uri=${CallbackUrl}&client_secret=${Secret}&code=${Code}&grant_type=authorization_code"
)

type OauthService interface {
	GetAuthURL() (error, string)
	GetConfiguration() (auth *Auth, err error)
}

type OauthServiceImpl struct {
	configMaps *kube.ConfigMaps
}

func init() {
	app.Register(newOauthService)
}

func newOauthService(configMaps *kube.ConfigMaps) OauthService {
	return &OauthServiceImpl{
		configMaps: configMaps,
	}
}

type Auth struct {
	ApplicationId string
	Secret        string
	AuthURL       string
	CallbackUrl   string
}

func (o *OauthServiceImpl) GetAuthURL() (error, string) {
	t := replacer.GetMatches(OauthUrl)
	auth, err := o.GetConfiguration()
	if err != nil {
		return err, ""
	}
	baseUrl := ReplaceEnv(OauthUrl, t, auth)
	return nil, baseUrl
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

func (o *OauthServiceImpl) GetConfiguration() (auth *Auth, err error) {
	configMaps, err := o.configMaps.Get(constant.GitlabConstant, constant.TemplateDefaultNamespace)
	auth = &Auth{
		ApplicationId: configMaps.Data[constant.ApplicationId],
		Secret:        configMaps.Data[constant.Secret],
		CallbackUrl:   configMaps.Data[constant.CallbackUrl],
	}
	return
}
