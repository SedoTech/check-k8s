package configmaps

import (
	"fmt"

	"github.com/SedoTech/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"github.com/SedoTech/check-k8s/pkg/checks/api"
	"github.com/SedoTech/check-k8s/pkg/utils"
)

type (
	// CheckConfigMaps interface to check a deployment
	CheckConfigMaps interface {
		CheckExists(CheckExistsOptions) icinga.Result
	}

	checkConfigMapsImpl struct {
		Client    kubernetes.Interface
		Namespace string
	}
)

// NewCheckConfigMaps creates a new instance of CheckConfigMaps
func NewCheckConfigMaps(client kubernetes.Interface, name string, namespace string) CheckConfigMaps {
	return &checkConfigMapsImpl{Client: client, Namespace: namespace}
}

// CheckExistsOptions contains options needed to run CheckExists check
type CheckExistsOptions struct {
	ConfigMapsDefined []string
}

// CheckExists checks if the deployment has a minimum of available replicas
func (c *checkConfigMapsImpl) CheckExists(options CheckExistsOptions) icinga.Result {
	name := "ConfigMaps.Exists"

	// first check if all defined secrets are existent
	missing := []string{}
	for _, secretName := range options.ConfigMapsDefined {
		_, err := api.GetConfigMap(c.Client, api.GetConfigMapOptions{Name: secretName, Namespace: c.Namespace})
		if err != nil {
			missing = append(missing, secretName)
		}
	}

	// then check if there are secrets in the namespace, we forgot to check
	odd := []string{}
	allConfigMaps, err := api.GetConfigMaps(c.Client, api.GetConfigMapsOptions{Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult("GetConfigMaps", icinga.ServiceStatusUnknown, fmt.Sprintf("cant't list secrets: %v", err))
	}
	for _, s := range allConfigMaps.Items {
		if !utils.Contains(s.Name, options.ConfigMapsDefined) {
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
