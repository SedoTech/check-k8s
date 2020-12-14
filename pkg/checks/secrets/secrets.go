package secrets

import (
	"fmt"
	"strings"

	"github.com/benkeil/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"SedoTech/check-k8s/pkg/checks/api"
	"SedoTech/check-k8s/pkg/utils"
)

type (
	// CheckSecrets interface to check a deployment
	CheckSecrets interface {
		CheckExists(CheckExistsOptions) icinga.Result
	}

	checkSecretsImpl struct {
		Client    kubernetes.Interface
		Namespace string
	}
)

// NewCheckSecrets creates a new instance of CheckSecrets
func NewCheckSecrets(client kubernetes.Interface, namespace string) CheckSecrets {
	return &checkSecretsImpl{Client: client, Namespace: namespace}
}

// CheckExistsOptions contains options needed to run CheckExists check
type CheckExistsOptions struct {
	SecretsDefined []string
}

// CheckExists checks if the deployment has a minimum of available replicas
func (c *checkSecretsImpl) CheckExists(options CheckExistsOptions) icinga.Result {
	name := "Secrets.Exists"

	// first check if all defined secrets are existent
	missing := []string{}
	for _, secretName := range options.SecretsDefined {
		_, err := api.GetSecret(c.Client, api.GetSecretOptions{Name: secretName, Namespace: c.Namespace})
		if err != nil {
			missing = append(missing, secretName)
		}
	}

	// then check if there are secrets in the namespace, we forgot to check
	odd := []string{}
	allSecrets, err := api.GetSecrets(c.Client, api.GetSecretsOptions{Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult("GetSecrets", icinga.ServiceStatusUnknown, fmt.Sprintf("cant't list secrets: %v", err))
	}
	for _, s := range allSecrets.Items {
		if !utils.Contains(s.Name, options.SecretsDefined) && !strings.HasPrefix(s.Name, "default-token-") {
			odd = append(odd, s.Name)
		}
	}

	// return a warning if we found secrets that are not checked
	// return a critical if a secret is not installed
	status := icinga.ServiceStatusOk
	message := icinga.DefaultSuccessMessage
	messageMissing := ""
	messageOdd := ""
	if len(odd) > 0 {
		status = icinga.ServiceStatusWarning
		messageOdd = fmt.Sprintf("installed but not checked: %v", odd)
	}
	if len(missing) > 0 {
		status = icinga.ServiceStatusCritical
		messageMissing = fmt.Sprintf("missing: %v ", missing)
	}

	if len(odd) > 0 || len(missing) > 0 {
		message = fmt.Sprintf("%v%v", messageMissing, messageOdd)
	}

	return icinga.NewResult(name, status, message)
}
