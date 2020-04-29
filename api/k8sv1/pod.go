package k8sv1

import (
	"sync"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	return 100
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

// PodDetail returns details for a pod
func (k *Client) PodDetail(options PodDetailOptions) (po *PodDetail, apiErr *errs.APIError) {
	po = &PodDetail{}

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	pod, err := clientset.CoreV1().Pods(options.Namespace).Get(options.Name, metav1.GetOptions{
		IncludeUninitialized: true,
	})

	if !options.UserCanAccess(pod.GetLabels()) {
		return nil, errs.Forbidden()
	}

	// po.Set(*pod, options.UserCanAccess)
	var st string
	if pod.Status.StartTime != nil {
		st = pod.Status.StartTime.String()
	}

	// remove environment variables for containers
	// that might use env vars for secrets. This could be added back
	// auth roles and such is figured out better.
	spec := &pod.Spec
	for i, v := range spec.Containers {
		// Add all container names for easier searching
		po.ContainerNames = append(po.ContainerNames, v.Name)

		if !options.UserRole.HasEnvVarsAccess(pod.GetLabels()) {
			spec.Containers[i].Env = nil
		}
	}

	// add the current pod
	po.Name = pod.GetName()
	po.Namespace = pod.GetNamespace()
	po.HostIP = pod.Status.HostIP
	po.PodIP = pod.Status.PodIP
	po.StartTime = st
	po.Phase = pod.Status.Phase
	po.PhaseMessage = pod.Status.Message
	po.ContainerStatus = pod.Status.ContainerStatuses
	po.Status = pod.Status
	po.Spec = *spec

	return po, nil
}

// PodOverview returns an overview of pods related to an application
func (k *Client) PodOverview(options PodOverviewOptions) (po *PodOverview, apiErr *errs.APIError) {
	if !options.UserCanAccess() {
		return nil, errs.Forbidden()
	}

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.CoreV1().Pods(options.Namespace).List(metav1.ListOptions{
		LabelSelector:        toLabelSelectorString(options.LabelSelector),
		IncludeUninitialized: true,
		Limit:                options.GetLimit(),
	})

	if err != nil {
		return nil, errs.InternalServerError(err.Error())
	}

	po = &PodOverview{
		PodInfo: []*PodInfo{},
	}

	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	for i, pod := range list.Items {
		go func(index int, pod v1.Pod) {
			defer wg.Done()
			if options.UserRole.HasNamespaceAccess(pod.GetNamespace()) &&
				options.UserRole.HasPodAccess(pod.GetLabels()) &&
				options.UserRole.Matches(pod.GetLabels(), &options.AppName) {
				// set common overivew fields on first pass
				if index == 0 {
					// set app overview and initialize pod quickview
					po.Name = pod.GetName()
					po.Namespace = pod.GetNamespace()
					po.ClusterName = pod.GetClusterName()
					po.DeployerLink = getDeployerLink(pod.GetName())
				}

				var st string
				if pod.Status.StartTime != nil {
					st = pod.Status.StartTime.String()
				}

				// remove environment variables for containers
				// that might use env vars for secrets. This could be added back
				// auth roles and such is figured out better.
				spec := &pod.Spec
				containerNames := []string{}
				for i, v := range spec.Containers {
					// Add all container names for easier searching
					containerNames = append(containerNames, v.Name)

					if !options.UserRole.HasEnvVarsAccess(pod.GetLabels()) {
						spec.Containers[i].Env = nil
					}
				}

				pi := &PodInfo{
					Name:         pod.GetName(),
					Namespace:    pod.GetNamespace(),
					HostIP:       pod.Status.HostIP,
					PodIP:        pod.Status.PodIP,
					StartTime:    st,
					Phase:        string(pod.Status.Phase),
					PhaseMessage: pod.Status.Message,
					Conditions:   pod.Status.Conditions,
				}

				for _, container := range pod.Spec.Containers {
					pi.Images = append(pi.Images, Image{
						ContainerName: container.Name,
						Name:          container.Image,
					})
				}

				// add the current pod
				po.PodInfo = append(po.PodInfo, pi)
			}
		}(i, pod)
	}

	wg.Wait()

	return po, nil
}
