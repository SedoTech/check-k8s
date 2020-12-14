package api

import (
	"fmt"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// GetSecretOptions options to get a deployment
	GetSecretOptions struct {
		Name      string
		Namespace string
	}
)

var secrets = make(map[string]*v1.Secret)

// GetSecret returns a single k8s Secret resource
func GetSecret(client kubernetes.Interface, options GetSecretOptions) (*v1.Secret, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if secret, found := secrets[key]; found {
		return secret, nil
	}
	resource, err := client.CoreV1().Secrets(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	secrets[key] = resource
	return resource, nil
}

type (
	// GetSecretsOptions options to get a deployment
	GetSecretsOptions struct {
		Namespace string
	}
)

// GetSecrets returns a SecretList
func GetSecrets(client kubernetes.Interface, options GetSecretsOptions) (*v1.SecretList, error) {
	resources, err := client.CoreV1().Secrets(options.Namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resources, nil
}
