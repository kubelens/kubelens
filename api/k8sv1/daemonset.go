package k8sv1

import (
	"sync"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

// DaemonSetOverviews returns a list of daemonsets given filter options
func (k *Client) DaemonSetOverviews(options DaemonSetOptions) (daemonsets []DaemonSetOverview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	lo := metav1.ListOptions{
		IncludeUninitialized: true,
	}

	if len(options.LabelSelector) > 0 {
		lo.LabelSelector = toLabelSelectorString(options.LabelSelector)
	}

	list, err := clientset.AppsV1().DaemonSets(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		daemonsets = make([]DaemonSetOverview, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			if options.UserRole.HasDaemonSetAccess(item.GetLabels()) {
				go func(index int, ds appsv1.DaemonSet) {
					defer wg.Done()

					var labelSelector map[string]string
					// this shouldn't be null, but default to regular labels if it is.
					if ds.Spec.Selector != nil {
						labelSelector = ds.Spec.Selector.MatchLabels
					} else if len(options.LabelSelector) > 0 {
						labelSelector = options.LabelSelector
					} else {
						labelSelector = ds.GetLabels()
					}

					name := getFriendlyAppName(
						ds.GetLabels(),
						ds.GetName(),
					)

					dso := DaemonSetOverview{
						FriendlyName:           name,
						Name:                   ds.GetName(),
						Namespace:              ds.GetNamespace(),
						LabelSelector:          labelSelector,
						CurrentNumberScheduled: int(ds.Status.CurrentNumberScheduled),
						DesiredNumberScheduled: int(ds.Status.DesiredNumberScheduled),
						NumberAvailable:        int(ds.Status.NumberAvailable),
						NumberMisscheduled:     int(ds.Status.NumberMisscheduled),
						NumberReady:            int(ds.Status.NumberReady),
						NumberUnavailable:      int(ds.Status.NumberUnavailable),
						UpdatedNumberScheduled: int(ds.Status.UpdatedNumberScheduled),
						Conditions:             ds.Status.Conditions,
					}

					// add deployments per daemon set for ease of display by client since deployments are really
					// specific to certian K8s kinds.
					deployments, err := k.DeploymentOverviews(DeploymentOptions{
						LabelSelector: labelSelector,
						Namespace:     ds.GetNamespace(),
						UserRole:      options.UserRole,
						Logger:        options.Logger,
					})
					// just trace the error and move on, shouldn't be critical.
					if err != nil {
						klog.Trace()
					}

					if len(deployments) > 0 {
						dso.AddDeploymentOverviews(&deployments)
					}

					cms, err := k.ConfigMaps(ConfigMapOptions{
						Namespace:     ds.GetNamespace(),
						LabelSelector: labelSelector,
						Logger:        options.Logger,
						UserRole:      options.UserRole,
					})

					// just trace the error and move on, shouldn't be critical.
					if err != nil {
						klog.Trace()
					}

					if len(cms) > 0 {
						dso.AddConfigMaps(&cms)
					}

					daemonsets[index] = dso
				}(i, item)
			}
		}

		wg.Wait()
	}

	ret := []DaemonSetOverview{}
	for _, item := range daemonsets {
		if len(item.Name) > 0 {
			ret = append(ret, item)
		}
	}

	return ret, nil
}

// DaemonSetAppInfos returns basic info for all daemon sets found for a given namespace.
func (k *Client) DaemonSetAppInfos(options DaemonSetOptions) (info []AppInfo, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.AppsV1().DaemonSets(options.Namespace).List(metav1.ListOptions{IncludeUninitialized: true})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	errors := []*errs.APIError{}

	if list != nil {
		info = make([]AppInfo, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			go func(index int, ds appsv1.DaemonSet) {
				defer wg.Done()

				var labelSelector map[string]string
				// this shouldn't be null, but default to regular labels if it is.
				if ds.Spec.Selector != nil {
					labelSelector = ds.Spec.Selector.MatchLabels
				} else {
					labelSelector = ds.GetLabels()
				}

				name := getFriendlyAppName(
					ds.GetLabels(),
					ds.GetName(),
				)

				info[index] = AppInfo{
					FriendlyName:  name,
					Name:          ds.GetName(),
					Namespace:     ds.GetNamespace(),
					Kind:          "DaemonSet",
					LabelSelector: labelSelector,
				}
			}(i, item)
		}
		wg.Wait()
	}

	if len(errors) > 0 {
		return info, errs.ListToInternalServerError(errors)
	}

	return info, nil
}
