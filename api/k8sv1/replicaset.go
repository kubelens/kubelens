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

// ReplicaSetOptions contains fields used for filtering when retrieving daemon sets
type ReplicaSetOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// ReplicaSetOverview .
type ReplicaSetOverview struct {
	// the name of the application
	FriendlyName string `json:"friendlyName"`
	// the name of the deployment
	Name string `json:"name"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector        map[string]string `json:"labelSelector"`
	AvailableReplicas    int               `json:"availableReplicas"`
	FullyLabeledReplicas int               `json:"fullyLabeledReplicas"`
	ReadyReplicas        int               `json:"readyReplicas"`
	Replicas             int               `json:"replicas"`
	// represents the latest available observations of a deployment's current state.
	Conditions          []appsv1.ReplicaSetCondition `json:"conditions"`
	ConfigMaps          *[]v1.ConfigMap              `json:"configMaps,omitempty"`
	DeploymentOverviews *[]DeploymentOverview        `json:"deploymentOverviews,omitempty"`
}

// AddConfigMaps sets the ConfigMaps value. Normally wouldn't have a pointer for a slice,
// but this allows to easily return empty for the client.
func (d *ReplicaSetOverview) AddConfigMaps(cms *[]v1.ConfigMap) {
	d.ConfigMaps = cms
}

// AddDeploymentOverviews sets the DeploymentOverviews value. Adding separately for ease of
// checking access before hand.
func (d *ReplicaSetOverview) AddDeploymentOverviews(dp *[]DeploymentOverview) {
	d.DeploymentOverviews = dp
}

// ReplicaSetOverviews returns a list of replicasets given filter options
func (k *Client) ReplicaSetOverviews(options ReplicaSetOptions) (replicasets []ReplicaSetOverview, apiErr *errs.APIError) {
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

	list, err := clientset.AppsV1().ReplicaSets(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		replicasets = make([]ReplicaSetOverview, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			if options.UserRole.HasDaemonSetAccess(item.GetLabels()) {
				go func(index int, rs appsv1.ReplicaSet) {
					defer wg.Done()

					var labelSelector map[string]string
					// this shouldn't be null, but default to regular labels if it is.
					if rs.Spec.Selector != nil {
						labelSelector = rs.Spec.Selector.MatchLabels
					} else if len(options.LabelSelector) > 0 {
						labelSelector = options.LabelSelector
					} else {
						labelSelector = rs.GetLabels()
					}

					name := getFriendlyAppName(
						rs.GetLabels(),
						rs.GetName(),
					)

					rso := ReplicaSetOverview{
						FriendlyName:         name,
						Name:                 rs.GetName(),
						Namespace:            rs.GetNamespace(),
						LabelSelector:        labelSelector,
						AvailableReplicas:    int(rs.Status.AvailableReplicas),
						FullyLabeledReplicas: int(rs.Status.FullyLabeledReplicas),
						ReadyReplicas:        int(rs.Status.ReadyReplicas),
						Replicas:             int(rs.Status.Replicas),
						Conditions:           rs.Status.Conditions,
					}

					// add deployments per daemon set for ease of display by client since deployments are really
					// specific to certian K8s kinds.
					if options.UserRole.HasDeploymentAccess(rs.GetLabels()) {
						deployments, err := k.DeploymentOverviews(DeploymentOptions{
							LabelSelector: labelSelector,
							Namespace:     rs.GetNamespace(),
							UserRole:      options.UserRole,
							Logger:        options.Logger,
						})
						// just trace the error and move on, shouldn't be critical.
						if err != nil {
							klog.Trace()
						}

						if len(deployments) > 0 {
							rso.AddDeploymentOverviews(&deployments)
						}
					}

					if options.UserRole.HasConfigMapAccess(item.GetLabels()) {
						cms, err := k.ConfigMaps(ConfigMapOptions{
							Namespace:     rs.GetNamespace(),
							LabelSelector: labelSelector,
						})

						// just trace the error and move on, shouldn't be critical.
						if err != nil {
							klog.Trace()
						}

						if len(cms) > 0 {
							rso.AddConfigMaps(&cms)
						}
					}

					replicasets[index] = rso
				}(i, item)
			}
		}

		wg.Wait()
	}

	return replicasets, nil
}

// ReplicaSetAppInfos returns basic info for all replica sets found for a given namespace.
func (k *Client) ReplicaSetAppInfos(options ReplicaSetOptions) (info []AppInfo, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.AppsV1().ReplicaSets(options.Namespace).List(metav1.ListOptions{IncludeUninitialized: true})

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
			go func(index int, rs appsv1.ReplicaSet) {
				defer wg.Done()
				if err != nil {
					klog.Trace()
					errors = append(errors, errs.InternalServerError(err.Error()))
				}

				var labelSelector map[string]string
				// this shouldn't be null, but default to regular labels if it is.
				if rs.Spec.Selector != nil {
					labelSelector = rs.Spec.Selector.MatchLabels
				} else {
					labelSelector = rs.GetLabels()
				}

				name := getFriendlyAppName(
					rs.GetLabels(),
					rs.GetName(),
				)

				info[index] = AppInfo{
					FriendlyName:  name,
					Name:          rs.GetName(),
					Namespace:     rs.GetNamespace(),
					Kind:          "ReplicaSet",
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
