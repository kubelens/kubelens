package k8sv1

import (
	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	klog "github.com/kubelens/kubelens/api/log"
)

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

// DeploymentOverviews returns a list of deployments given filter options
func (k *Client) DeploymentOverviews(options DeploymentOptions) (deployments []DeploymentOverview, apiErr *errs.APIError) {
	deployments = []DeploymentOverview{}
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

	list, err := clientset.AppsV1().Deployments(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			if options.UserRole.HasDeploymentAccess(item.GetLabels()) {
				var labelSelector map[string]string
				// this shouldn't be null, but default to regular labels if it is.
				if item.Spec.Selector != nil {
					labelSelector = item.Spec.Selector.MatchLabels
				} else if len(options.LabelSelector) > 0 {
					labelSelector = options.LabelSelector
				} else {
					labelSelector = item.GetLabels()
				}

				name := getFriendlyAppName(
					item.GetLabels(),
					item.GetName(),
				)

				deployments = append(deployments, DeploymentOverview{
					FriendlyName:         name,
					Name:                 item.GetName(),
					Namespace:            item.GetNamespace(),
					LabelSelector:        labelSelector,
					ResourceVersion:      item.GetResourceVersion(),
					AvailableReplicas:    int(item.Status.AvailableReplicas),
					ReadyReplicas:        int(item.Status.ReadyReplicas),
					Replicas:             int(item.Status.Replicas),
					UnavailableReplicas:  int(item.Status.UnavailableReplicas),
					UpdatedReplicas:      int(item.Status.UnavailableReplicas),
					DeploymentConditions: item.Status.Conditions,
				})
			}
		}
	}

	return deployments, nil
}
