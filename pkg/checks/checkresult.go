package checks

import (
	"fmt"

	"git.i.sedorz.net/infrastructure/icinga/check-k8s/pkg/icinga"
)

type (
	// ServiceCheckResult interface for service check results
	ServiceCheckResult interface {
		Name() string
		State() icinga.ServiceState
		Message() string
	}

	serviceCheckResultImpl struct {
		name         string
		serviceState icinga.ServiceState
		message      string
	}

	// HostCheckResult interface for host check results
	HostCheckResult interface {
		Name() string
		State() icinga.HostState
		Message() string
	}

	hostCheckResult struct {
		name      string
		hostState icinga.HostState
		message   string
	}
)

// NewServiceCheckResult creates a new instance of ServiceCheckResult
func NewServiceCheckResult(name string, state icinga.ServiceState, message string) ServiceCheckResult {
	return &serviceCheckResultImpl{name, state, message}
}

// NewServiceCheckResultOk creates a new instance of ServiceCheckResult and set result to ServiceStateOk
func NewServiceCheckResultOk(name string) ServiceCheckResult {
	return &serviceCheckResultImpl{name, icinga.ServiceStateOk, "everything ok"}
}

// NewServiceCheckResultOkMessage creates a new instance of ServiceCheckResult and set result to ServiceStateOk
func NewServiceCheckResultOkMessage(name string, message string) ServiceCheckResult {
	return &serviceCheckResultImpl{name, icinga.ServiceStateOk, message}
}

func (r *serviceCheckResultImpl) Name() string {
	return r.name
}

func (r *serviceCheckResultImpl) State() icinga.ServiceState {
	return r.serviceState
}

func (r *serviceCheckResultImpl) Message() string {
	return r.message
}

func (r *serviceCheckResultImpl) String() string {
	return fmt.Sprintf("{name: %s serviceState: %s message: %s}", r.name, r.serviceState, r.message)
}

// NewHostCheckResult creates a new instance of HostCheckResult
func NewHostCheckResult(name string, state icinga.HostState, message string) HostCheckResult {
	return &hostCheckResult{name, state, message}
}

// NewHostCheckResultOk creates a new instance of HostCheckResult and set result to HostStateUp
func NewHostCheckResultOk(name string) HostCheckResult {
	return &hostCheckResult{name, icinga.HostStateUp, "everything ok"}
}

// NewHostCheckResultOkMessage creates a new instance of HostCheckResult and set result to HostStateUp
func NewHostCheckResultOkMessage(name string, message string) HostCheckResult {
	return &hostCheckResult{name, icinga.HostStateUp, message}
}

func (r *hostCheckResult) Name() string {
	return r.name
}

func (r *hostCheckResult) State() icinga.HostState {
	return r.hostState
}

func (r *hostCheckResult) Message() string {
	return r.message
}

type (
	// ServiceCheckResults contains multiple results for service checks
	ServiceCheckResults interface {
		All() []ServiceCheckResult
		Add(ServiceCheckResult)
		Get(string) ServiceCheckResult
		CalculateServiceState() icinga.ServiceState
	}

	serviceCheckResultsImpl struct {
		results map[string]ServiceCheckResult
	}
)

// NewServiceCheckResults creates a new instance of ServiceCheckResults
func NewServiceCheckResults() ServiceCheckResults {
	return &serviceCheckResultsImpl{make(map[string]ServiceCheckResult)}
}

// Add adds a element to the set
func (r *serviceCheckResultsImpl) Add(result ServiceCheckResult) {
	r.results[result.Name()] = result
}

// Contains return if a element is in the set
func (r *serviceCheckResultsImpl) Contains(name string) bool {
	_, ok := r.results[name]
	return ok
}

// Get gets a element from the set
func (r *serviceCheckResultsImpl) Get(name string) ServiceCheckResult {
	return r.results[name]
}

// All returns all values from the set
func (r *serviceCheckResultsImpl) All() []ServiceCheckResult {
	values := []ServiceCheckResult{}
	for _, value := range r.results {
		values = append(values, value)
	}
	return values
}

// CalculateServiceState calculates the service state for multiple checks
func (r *serviceCheckResultsImpl) CalculateServiceState() icinga.ServiceState {
	calculatedState := icinga.ServiceStateOk
	for _, serviceStateResult := range r.results {
		if serviceStateResult.State() > calculatedState {
			calculatedState = serviceStateResult.State()
		}
	}
	return calculatedState
}
