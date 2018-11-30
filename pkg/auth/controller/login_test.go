// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/hioak/starter/scm"
	"hidevops.io/mio/pkg/auth/service/mocks"
	"net/http"
	"os"
	"testing"
)

func TestUserLogin(t *testing.T) {

	io.EnsureWorkDir(1, "config/application.yml")
	os.Setenv("SCM_URL", "https://github.com")
	loginService := new(mocks.LoginService)
	jwtToken := jwt.NewJwtToken(&jwt.Properties{
		PrivateKeyPath: "config/ssl/app.rsa",
		PublicKeyPath:  "config/ssl/app.rsa.pub",
	})
	login := newLoginController(jwtToken, loginService)
	testApp := web.NewTestApp(t, login).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	loginService.On("GetSession", "https://github.com", "1", "2").Return("", 1, "", nil)
	loginService.On("GetUser", "https://github.com", "").Return(&scm.User{}, nil)
	user := &UserRequest{
		Username: "1",
		Password: "2",
	}
	t.Run("should pass with jwt token", func(t *testing.T) {
		testApp.Post("/login").WithJSON(user).Expect().Status(http.StatusOK)
	})
}
