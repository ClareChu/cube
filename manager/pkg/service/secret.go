package service

import (
	"encoding/base64"
	"hidevops.io/cube/manager/pkg/service/client"
	"hidevops.io/hiboot/pkg/app"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretService interface {
}

type Secret struct {
	SecretService
}

func init() {
	app.Register(NewSecret)
}

func NewSecret() SecretService {
	return &Secret{
	}
}

func (s *Secret) CreateTraefikSecret(namespace string) error {
	clientSet, err := client.GetDefaultK8sClientSet()
	if err != nil {
		return err
	}
	_, err = clientSet.CoreV1().Secrets(namespace).Get("traefik-cert", metav1.GetOptions{})
	if err == nil {
		return err
	}
	crt, err := ioutil.ReadFile("/ssl/tls.crt")
	if err != nil {
		return err
	}
	enCrt := base64.StdEncoding.EncodeToString(crt)
	key, err := ioutil.ReadFile("/ssl/tls.key")
	enKey := base64.StdEncoding.EncodeToString(key)
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traefik-cert",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"tls.crt": []byte(enCrt),
			"tls.key": []byte(enKey),
		},
	}
	secret, err = clientSet.CoreV1().Secrets(namespace).Create(secret)
	return err
}
