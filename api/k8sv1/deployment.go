package k8sv1

import (
	"github.com/kubelens/kubelens/api/errs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	klog "github.com/kubelens/kubelens/api/log"
)

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
