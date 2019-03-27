package builder

import (
	"hidevops.io/cube/manager/pkg/constant"
	"net/smtp"
	"strings"
)

func SendEMail(toUsers []string, pipelineName string) error {
	auth := smtp.PlainAuth("", constant.EmailUsername, constant.EmailPassword, constant.EmailHost)
	nickname := "test"
	user := constant.EmailUsername
	subject := constant.EmailSubject
	contentType := "Content-Type: text/plain; charset=UTF-8"
	body := constant.EmailBody
	msg := []byte("To: " + strings.Join(toUsers, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail(constant.EmailAddr, auth, user, toUsers, msg)
	return err
}
