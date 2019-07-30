package k8v1

import (
	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
)

type AppOverviewOptions struct {
	// name of the  application
	AppName string `json:"appname"`
	// namespace of the app
	Namespace string `json:"namespace"`
	// the label key used to tag the appication
	LabelKey string `json:"labelKey"`
	// include detail
	Detailed bool `json:"detailed"`
	// user roles
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

type AppOverview struct {
	PodOverviews     PodOverview       `json:"podOverviews,omitempty"`
	ServiceOverviews []ServiceOverview `json;"serviceOverviews,omitempty"`
}

// Name holds fiedls for a name. This can be used to match the label key with the value(name) of the
// app/pod/component/etc
type Name struct {
	// the label key associated with the name. this can be used to search
	// labels matching the key
	LabelKey string `json:"labelKey"`
	// the actual name
	Value string `json:"value"`
}

// PodOverviewOptions contains fields used for filtering when retrieving application overiew(s).
type PodOverviewOptions struct {
	// name of the application
	AppName string `json:"appname"`
	// name of the pod
	PodName string `json:"podname"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	// The label key used for the application name, ex. app=some-app-name
	AppNameLabelKey string `json:"appNameLabelKey"`
	// Limit the number of pod summaries to return
	// Use function GetLimit to get the default limit or this overriden value.
	Limit int64 `json:"linit"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// ValidNamespace retuns true if the length self.Namespace is > 0
func (a *PodOverviewOptions) ValidNamespace() bool {
	if len(a.Namespace) > 0 {
		return true
	}
	return false
}

// GetLimit returns PodOverviewOptions.Limit > 0 || default (32, why not)
func (a *PodOverviewOptions) GetLimit() int64 {
	if a.Limit > 0 {
		return a.Limit
	}
	return 32
}

// UserCanAccess validates the user has access
func (a *PodOverviewOptions) UserCanAccess() bool {
	if !a.UserRole.HasApplicationAccess() || !a.UserRole.HasNamespaceAccess(a.Namespace) {
		return false
	}
	return true
}

// PodOverview is meant to hold higher level application information
type PodOverview struct {
	// the name of the application
	Name Name `json:"name"`
	// the namespace (if supported)
	Namespace string `json:"namespace,omitempty"`
	// the cluster name (if supported)
	ClusterName string `json:"clusterName,omitempty"`
	// the link to the tool that deploys the application
	DeployerLink string `json:"deployerLink,omitempty"`
	// the pods associated with the application
	PodDetails []*PodDetail `json:"podDetails,omitempty"`
}

type PodDetailOptions struct {
	// the name of the pod
	Name string `json:"name"`
	// the namespace of the pod
	Namespace string `json:"namespace"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

func (p *PodDetailOptions) Valid() bool {
	return len(p.Name) > 0
}

// UserCanAccess validates the user has access
func (p *PodDetailOptions) UserCanAccess(labels map[string]string) bool {
	if !p.UserRole.HasPodAccess(labels) || !p.UserRole.HasNamespaceAccess(p.Namespace) {
		return false
	}
	return true
}

// PodDetail holds more detailed information for a pod
type PodDetail struct {
	// name of the pod
	Name string `json:"name"`
	// the namespace of the pod
	Namespace string `json:"namespace"`
	// the pod's host ip
	HostIP string `json:"hostIP,omitempty"`
	// the ip of the pod
	PodIP string `json:"podIP,omitempty"`
	// time when the pod started
	StartTime string `json:"startTime,omitempty"`
	// phase is the "phase" status of the pod
	Phase v1.PodPhase `json:"phase,omitempty"`
	// message description of the phase
	PhaseMessage string `json:"phaseMessage,omitempty"`
	// status of the container
	ContainerStatus []v1.ContainerStatus `json:"containerStatus,omitempty"`
	// the status object of the pod
	Status v1.PodStatus `json:"status,omitempty"`
	// the pod's spec
	Spec v1.PodSpec `json:"spec,omitempty"`
	// the pod's log output
	Log *Log `json:"log,omitempty"`
}

// Log holds fields related to log output
type Log struct {
	// the name of the pod
	Pod string `json:"pod"`
	// the log output
	Output string `json:"output"`
}

// LogOptions contains fields used for filtering when retrieving application logs
type LogOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// The label key used for the application name, ex. app=some-app-name
	PodName string `json:"podname"`
	// follow enables streaming
	Follow bool `json:"follow"`
	// tail logs from line. If a stream request, this is ignored.
	Tail int64 `json:"tail"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// Valid validates LogOptions fields
func (a *LogOptions) Valid() *errs.APIError {
	if !a.UserRole.HasNamespaceAccess(a.Namespace) {
		return errs.Unauthorized()
	}

	if !a.ValidPodName() {
		return errs.ValidationError("podname must be provided when getting logs")
	}

	if !a.ValidNamespace() {
		return errs.ValidationError("namespace must be provided when getting logs")
	}

	return nil
}

// ValidNamespace retuns true if the length LogOptions.Namespace is > 0
func (a *LogOptions) ValidNamespace() bool {
	if len(a.Namespace) > 0 {
		return true
	}
	return false
}

// ValidPodName retuns true if the length LogOptions.PodName is > 0
func (a *LogOptions) ValidPodName() bool {
	if len(a.PodName) > 0 {
		return true
	}
	return false
}

// GetTailLines returns PodOverviewOptions.Tail > 0 || default (32, why not)
func (a *LogOptions) GetTailLines() int64 {
	if a.Tail > 0 {
		return a.Tail
	}
	return 100
}

// UserCanAccess validates the user has access
func (a *LogOptions) UserCanAccess(labels map[string]string) bool {
	if !a.UserRole.HasPodAccess(labels) || !a.UserRole.HasLogAccess(labels) || !a.UserRole.Matches(labels, nil) {
		return false
	}
	return true
}

// ServiceOverview provides field relating to a service. Can return full values or slimmed down values.
type ServiceOverview struct {
	AppName        Name              `json:"appName"`
	DeployerLink   string            `json:"deployerLink,omitempty"`
	Name           string            `json:"name"`
	Namespace      string            `json:"namespace"`
	Type           v1.ServiceType    `json:"type"`
	Selector       map[string]string `json:"selector,omitempty"`
	ClusterIP      string            `json:"clusterIP,omitempty"`
	LoadBalancerIP string            `json:"loadBalancerIP,omitempty"`
	ExternalIPs    []string          `json:"externalIPs,omitempty"`
	Ports          []v1.ServicePort  `json:"ports,omitempty"`
	Spec           *v1.ServiceSpec   `json:"spec,omitempty"`
	Status         *v1.ServiceStatus `json:"status,omitempty"`
}

// AddDetail sets additional fields for more detailed overview
func (s *ServiceOverview) AddDetail(spec *v1.ServiceSpec, status *v1.ServiceStatus) {
	s.Type = spec.Type
	s.Ports = spec.Ports
	s.ClusterIP = spec.ClusterIP
	s.LoadBalancerIP = spec.LoadBalancerIP
	s.Selector = spec.Selector
	s.Spec = spec
	s.Status = status
}

// ServiceOptions contains fields used for filtering when retrieving services
type ServiceOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label search used for the application name, ex. app=some-app-name, component=some-component
	LabelSearch string `json:"labelSearch"`
	// include full detail
	Detailed bool `json:"detailed"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// ValidNamespace validates the namespace field
func (a *ServiceOptions) ValidNamespace() bool {
	if len(a.Namespace) > 0 {
		return true
	}
	return false
}

// ValidLabelSearch validates the label search field
func (a *ServiceOptions) ValidLabelSearch() bool {
	if len(a.LabelSearch) > 0 {
		return true
	}
	return false
}

// UserCanAccess validates the user has access
func (a *ServiceOptions) UserCanAccess(labels map[string]string) bool {
	if !a.UserRole.HasApplicationAccess() || !a.UserRole.HasServiceAccess(labels) {
		return false
	}
	return true
}
