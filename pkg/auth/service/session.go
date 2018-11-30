package service

import (
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/replacer"
)

type Session struct {
	AuthURL       string
	ApplicationId string
	Secret        string
	TokenURL      string
	profileURL    string
	CallbackUrl   string
	Code          string
}

type SessionResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	CreateAt         string `json:"create_at"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewClient(authUrl, tokenUrl, applicationId, callbackUrl, secret string) *Session {
	s := &Session{
		AuthURL:       authUrl,
		TokenURL:      tokenUrl,
		ApplicationId: applicationId,
		CallbackUrl:   callbackUrl,
		Secret:        secret,
	}
	return s
}

func (session *Session) GetAccessToken(code string) (*SessionResponse, error) {
	log.Info("session GetAccessToken code : ", code)
	session.Code = code
	t := replacer.GetMatches(AccessTokenUrl)
	baseUrl := ReplaceEnv(AccessTokenUrl, t, session)
	sessionResponse := &SessionResponse{}
	_, err := Client("POST", baseUrl, sessionResponse)
	log.Info(sessionResponse)
	return sessionResponse, err
}
