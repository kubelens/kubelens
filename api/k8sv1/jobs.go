package k8sv1

import (
	"sync"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobOptions contains fields used for filtering when retrieving jobs
type JobOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// JobOverview .
type JobOverview struct {
	// the name of the application
	FriendlyName string `json:"friendlyName"`
	// the name of the deployment
	Name string `json:"name"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// Represents time when the job was acknowledged by the job controller.
	StartTime string `json:"startTime"`
	// Represents time when the job was completed. It is not guaranteed to
	CompletionTime string `json:"completionTime"`
	// The number of actively running pods.
	Active int `json:"active"`
	// The number of pods which reached phase Succeeded.
	Succeeded int `json:"succeeded"`
	// The number of pods which reached phase Failed.
	Failed int `json:"failed"`
	// The latest available observations of an object's current state.
	Conditions          []batchv1.JobCondition `json:"conditions"`
	ConfigMaps          *[]v1.ConfigMap        `json:"configMaps,omitempty"`
	DeploymentOverviews *[]DeploymentOverview  `json:"deploymentOverviews,omitempty"`
}

// AddConfigMaps sets the ConfigMaps value. Normally wouldn't have a pointer for a slice,
// but this allows to easily return empty for the client.
func (j *JobOverview) AddConfigMaps(cms *[]v1.ConfigMap) {
	j.ConfigMaps = cms
}

// AddDeploymentOverviews sets the DeploymentOverviews value. Adding separately for ease of
// checking access before hand.
func (j *JobOverview) AddDeploymentOverviews(dp *[]DeploymentOverview) {
	j.DeploymentOverviews = dp
}

// JobOverviews returns a list of jobs given filter options
func (k *Client) JobOverviews(options JobOptions) (jobs []JobOverview, apiErr *errs.APIError) {
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

	list, err := clientset.BatchV1().Jobs(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		jobs = make([]JobOverview, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			if options.UserRole.HasJobAccess(item.GetLabels()) {
				go func(index int, j batchv1.Job) {
					defer wg.Done()

					var labelSelector map[string]string
					// this shouldn't be null, but default to regular labels if it is.
					if j.Spec.Selector != nil {
						labelSelector = j.Spec.Selector.MatchLabels
					} else if len(options.LabelSelector) > 0 {
						labelSelector = options.LabelSelector
					} else {
						labelSelector = j.GetLabels()
					}

					name := getFriendlyAppName(
						j.GetLabels(),
						j.GetName(),
					)

					jo := JobOverview{
						FriendlyName:  name,
						Name:          j.GetName(),
						Namespace:     j.GetNamespace(),
						LabelSelector: labelSelector,
						Active:        int(j.Status.Active),
						Succeeded:     int(j.Status.Succeeded),
						Failed:        int(j.Status.Failed),
						Conditions:    j.Status.Conditions,
					}

					if j.Status.StartTime != nil {
						jo.StartTime = j.Status.StartTime.String()
					}

					if j.Status.CompletionTime != nil {
						jo.CompletionTime = j.Status.CompletionTime.String()
					}

					// add deployments per daemon set for ease of display by client since deployments are really
					// specific to certian K8s kinds.
					if options.UserRole.HasDeploymentAccess(j.GetLabels()) {
						deployments, err := k.DeploymentOverviews(DeploymentOptions{
							LabelSelector: labelSelector,
							Namespace:     j.GetNamespace(),
							UserRole:      options.UserRole,
							Logger:        options.Logger,
						})
						// just trace the error and move on, shouldn't be critical.
						if err != nil {
							klog.Trace()
						}

						if len(deployments) > 0 {
							jo.AddDeploymentOverviews(&deployments)
						}
					}

					if options.UserRole.HasConfigMapAccess(item.GetLabels()) {
						cms, err := k.ConfigMaps(ConfigMapOptions{
							Namespace:     j.GetNamespace(),
							LabelSelector: labelSelector,
						})

						// just trace the error and move on, shouldn't be critical.
						if err != nil {
							klog.Trace()
						}

						if len(cms) > 0 {
							jo.AddConfigMaps(&cms)
						}
					}

					jobs[index] = jo
				}(i, item)
			}
		}

		wg.Wait()
	}

	return jobs, nil
}

// JobAppInfos returns basic info for all jobs found for a given namespace.
func (k *Client) JobAppInfos(options JobOptions) (info []AppInfo, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.BatchV1().Jobs(options.Namespace).List(metav1.ListOptions{IncludeUninitialized: true})

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
			go func(index int, j batchv1.Job) {
				defer wg.Done()
				if err != nil {
					klog.Trace()
					errors = append(errors, errs.InternalServerError(err.Error()))
				}

				var labelSelector map[string]string
				// this shouldn't be null, but default to regular labels if it is.
				if j.Spec.Selector != nil {
					labelSelector = j.Spec.Selector.MatchLabels
				} else {
					labelSelector = j.GetLabels()
				}

				name := getFriendlyAppName(
					j.GetLabels(),
					j.GetName(),
				)

				info[index] = AppInfo{
					FriendlyName:  name,
					Name:          j.GetName(),
					Namespace:     j.GetNamespace(),
					Kind:          "Job",
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
