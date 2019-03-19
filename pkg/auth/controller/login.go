package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/mio/pkg/auth/service"
	"net/http"
	"time"
)

type UserRequest struct {
	model.RequestBody
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	PrivateToken string `json:"private_token"`
	Url          string `json:"url"`
	Uid          string `json:"uid"`
	Token        string `json:"token"`
}

type loginController struct {
	at.RestController
	token   jwt.Token
	loginService service.LoginService
}

func init() {
	app.Register(newLoginController)
}

func newLoginController(token jwt.Token, loginService service.LoginService) *loginController {
	return &loginController{
		token:   token,
		loginService: loginService,
	}
}

func (c *loginController) Post(request *UserRequest) (response model.Response, err error) {
	log.Debug("loginController.Login")
	response = new(model.BaseResponse)
	privateToken, uid, _, err := c.loginService.GetSession(request.Username, request.Password)
	if err != nil {
		response.SetCode(http.StatusNonAuthoritativeInfo)
		log.Error("username and password error ")
		return
	}
	u, err := c.loginService.GetUser(privateToken)
	token, err := c.token.Generate(jwt.Map{
		"username":      request.Username,
		"private_token": privateToken,
		"password":      request.Password,
		"uid":           uid,
	}, 10, time.Hour)
	data := map[string]interface{}{
		"token": token,
		"user":  u,
	}
	response.SetData(data)
	return
}
