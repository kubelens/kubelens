package k8sv1

import (
	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

// AppOptions .
type AppOptions struct {
	// user roles
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// App represents an application within Kubernetes.
type App struct {
	// name of the application
	Name string `json:"name"`
	// the namespace of the app
	Namespace string `json:"namespace"`
	// kind of application, e.g. Service, DaemonSet
	Kind string `json:"kind"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// deployer link if any
	DeployerLink string `json:"deployerLink,omitempty"`
}

// AppOverviewOptions .
type AppOverviewOptions struct {
	// name of the  application
	AppName string `json:"appname"`
	// namespace of the app
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// include detail
	Detailed bool `json:"detailed"`
	// user roles
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// AppOverview .
type AppOverview struct {
	PodOverviews       PodOverview         `json:"podOverviews,omitempty"`
	ServiceOverviews   []ServiceOverview   `json:"serviceOverviews,omitempty"`
	DaemonSetOverviews []DaemonSetOverview `json:"daemonSetOverviews,omitempty"`
}

// AppInfo .
type AppInfo struct {
	// friendly name of app
	FriendlyName string
	// actual name of app
	Name string
	// namspace the app is in
	Namespace string
	// the kind of app, e.g. Service, DaemonSet, etc.
	Kind string
	// the LabelSelector of the object.
	LabelSelector map[string]string
}

// PodOverviewOptions contains fields used for filtering when retrieving application overiew(s).
type PodOverviewOptions struct {
	// name of the application
	AppName string `json:"appname"`
	// name of the pod
	PodName string `json:"podname"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
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
	// the name of the app
	Name string `json:"name"`
	// the namespace (if supported)
	Namespace string `json:"namespace,omitempty"`
	// the cluster name (if supported)
	ClusterName string `json:"clusterName,omitempty"`
	// the link to the tool that deploys the application
	DeployerLink string `json:"deployerLink,omitempty"`
	// the pods associated with the application
	PodInfo []*PodInfo `json:"pods,omitempty"`
}

// PodDetailOptions .
type PodDetailOptions struct {
	// the name of the pod
	Name string `json:"name"`
	// the name of the container in the pod. Defaults to the application name for
	// backwards compatibility (original functionality)
	ContainerName string `json:"containerName"`
	// the namespace of the pod
	Namespace string `json:"namespace"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// Valid .
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
	// individual container names, there will always be at least 1.
	ContainerNames []string `json:"containerNames"`
}

// Image - fields pertaining to container image within a pod
type Image struct {
	Name          string `json:"name"`
	ContainerName string `json:"containerName"`
}

// PodInfo contains an aggregated/condensed view of the pod details.
type PodInfo struct {
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
	Phase string `json:"phase,omitempty"`
	// message description of the phase
	PhaseMessage string `json:"phaseMessage,omitempty"`
	// container images
	Images []Image `json:"images"`
	// pod conditions
	Conditions []v1.PodCondition `json:"conditions"`
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
	// The name of the container to get logs from.
	ContainerName string `json:"containerName"`
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
	FriendlyName        string               `json:"friendlyName"`
	DeployerLink        string               `json:"deployerLink,omitempty"`
	Name                string               `json:"name"`
	Namespace           string               `json:"namespace"`
	Type                v1.ServiceType       `json:"type"`
	Selector            map[string]string    `json:"selector,omitempty"`
	ClusterIP           string               `json:"clusterIP,omitempty"`
	LoadBalancerIP      string               `json:"loadBalancerIP,omitempty"`
	ExternalIPs         []string             `json:"externalIPs,omitempty"`
	Ports               []v1.ServicePort     `json:"ports,omitempty"`
	Spec                *v1.ServiceSpec      `json:"spec,omitempty"`
	Status              *v1.ServiceStatus    `json:"status,omitempty"`
	ConfigMaps          *[]v1.ConfigMap      `json:"configMaps,omitempty"`
	DeploymentOverviews []DeploymentOverview `json:"deploymentOverviews,omitempty"`
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

// AddConfigMaps sets the ConfigMaps value. Normally wouldn't have a pointer for a slice,
// but this allows to easily return empty for the client.
func (s *ServiceOverview) AddConfigMaps(cms *[]v1.ConfigMap) {
	s.ConfigMaps = cms
}

// AddDeploymentOverviews sets the DeploymentOverviews value. Adding separately for ease of
// checking access before hand.
func (s *ServiceOverview) AddDeploymentOverviews(d []DeploymentOverview) {
	s.DeploymentOverviews = d
}

// ServiceOptions contains fields used for filtering when retrieving services
type ServiceOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
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
	if len(a.LabelSelector) > 0 {
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

// DeploymentOptions contains fields used for filtering when retrieving deployments
type DeploymentOptions struct {
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// DeploymentOverview contains aggregated fields for a deployment
type DeploymentOverview struct {
	// the name of the application
	FriendlyName string `json:"friendlyName"`
	// the name of the deployment
	Name string `json:"name"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// the resource version of
	ResourceVersion string `json:"resourceVersion"`
	// replicas is the number of desired pods
	Replicas int `json:"replicas"`
	// total number of non-terminated pods targeted by this deployment that have the desired template spec.
	UpdatedReplicas int `json:"updatedReplicas"`
	// number of ready replicas
	ReadyReplicas int `json:"readyReplicas"`
	// number of available replicas
	AvailableReplicas int `json:"availableReplicas"`
	// number of unavailable replicas
	UnavailableReplicas int `json:"unavailableReplicas"`
	// represents the latest available observations of a deployment's current state.
	DeploymentConditions []appsv1.DeploymentCondition `json:"deploymentConditions"`
}

// DaemonSetOptions contains fields used for filtering when retrieving daemon sets
type DaemonSetOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// DaemonSetOverview .
type DaemonSetOverview struct {
	// the name of the application
	FriendlyName string `json:"friendlyName"`
	// the name of the deployment
	Name string `json:"name"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector          map[string]string `json:"labelSelector"`
	CurrentNumberScheduled int               `json:"currentNumberScheduled"`
	DesiredNumberScheduled int               `json:"desiredNumberScheduled"`
	NumberAvailable        int               `json:"numberAvailable"`
	NumberMisscheduled     int               `json:"numberMisscheduled"`
	NumberReady            int               `json:"numberReady"`
	NumberUnavailable      int               `json:"numberUnavailable"`
	UpdatedNumberScheduled int               `json:"updatedNumberScheduled"`
	// represents the latest available observations of a deployment's current state.
	Conditions          []appsv1.DaemonSetCondition `json:"conditions"`
	ConfigMaps          *[]v1.ConfigMap             `json:"configMaps,omitempty"`
	DeploymentOverviews *[]DeploymentOverview       `json:"deploymentOverviews,omitempty"`
}

// AddConfigMaps sets the ConfigMaps value. Normally wouldn't have a pointer for a slice,
// but this allows to easily return empty for the client.
func (d *DaemonSetOverview) AddConfigMaps(cms *[]v1.ConfigMap) {
	d.ConfigMaps = cms
}

// AddDeploymentOverviews sets the DeploymentOverviews value. Adding separately for ease of
// checking access before hand.
func (d *DaemonSetOverview) AddDeploymentOverviews(dp *[]DeploymentOverview) {
	d.DeploymentOverviews = dp
}
