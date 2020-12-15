package endpoints

import (
	"fmt"

	"github.com/SedoTech/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"github.com/SedoTech/check-k8s/pkg/checks/api"
)

type (
	// CheckEndpoints interface to check a deployment
	CheckEndpoints interface {
		CheckAvailableAddresses(CheckAvailableAddressesOptions) icinga.Result
		CheckAll(CheckAllOptions) icinga.Results
	}

	checkEndpointsImpl struct {
		Client    kubernetes.Interface
		Name      string
		Namespace string
	}
)

// NewCheckEndpoints creates a new instance of CheckEndpoints
func NewCheckEndpoints(client kubernetes.Interface, name string, namespace string) CheckEndpoints {
	return &checkEndpointsImpl{Client: client, Name: name, Namespace: namespace}
}

// CheckAllOptions contains options needed to run all deployment checks
type CheckAllOptions struct {
	Client                         kubernetes.Interface
	CheckAvailableAddressesOptions CheckAvailableAddressesOptions
}

// CheckAll runs all tests and returns an instance of ServiceCheckResults
func (c *checkEndpointsImpl) CheckAll(options CheckAllOptions) icinga.Results {
	results := icinga.NewResults()
	results.Add(c.CheckAvailableAddresses(options.CheckAvailableAddressesOptions))
	return results
}

// CheckAvailableAddressesOptions contains options needed to run CheckAvailableAddresses check
type CheckAvailableAddressesOptions struct {
	ThresholdWarning  string
	ThresholdCritical string
}

// CheckAvailableAddresses checks if the deployment has a minimum of available replicas
func (c *checkEndpointsImpl) CheckAvailableAddresses(options CheckAvailableAddressesOptions) icinga.Result {
	name := "Endpoints.AvailableAddresses"

	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}

	endpoints, err := api.GetEndpoints(c.Client, api.GetEndpointsOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult("GetEndpoints", icinga.ServiceStatusUnknown, fmt.Sprintf("cant't get endpoint: %v", err))
	}

	if len(endpoints.Subsets) != 1 {
		return icinga.NewResult("GetEndpoints", icinga.ServiceStatusUnknown, fmt.Sprintf("endpoint has %d subsets, dont know what to do", len(endpoints.Subsets)))
	}

	addresses := len(endpoints.Subsets[0].Addresses)
	status := statusCheck.CheckInt(addresses)
	message := fmt.Sprintf("endpoint has %v available address(ess)", addresses)

	return icinga.NewResult(name, status, message)
}
