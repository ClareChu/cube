package aggregate

import (
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
)

type SecretAggregate interface {
	Create(secret *command.Secret) error
	Get(name, namespace string) (*corev1.Secret, error)
}

type SecretAggregateImpl struct {
	SecretAggregate
	secretClient *kube.Secret
}

func init() {
	app.Register(NewSecretAggregate)
}

func NewSecretAggregate(secretClient *kube.Secret) SecretAggregate {
	return &SecretAggregateImpl{
		secretClient: secretClient,
	}
}

func (s *SecretAggregateImpl) Create(secret *command.Secret) error {
	log.Info()
	err := s.secretClient.Create(secret.Username, secret.Password, secret.Token, secret.Name, secret.Namespace)
	return err
}

func (s *SecretAggregateImpl) Get(name, namespace string) (*corev1.Secret, error) {
	log.Infof("get secret name:%v, namespace: %v", name, namespace)
	secret, err := s.secretClient.Get(name, namespace)
	return secret, err
}
