package k8v1

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
		LabelSelector:        options.LabelSelector,
		IncludeUninitialized: true,
	}

	// could match on deployment name, but this allows a bit more flexibility
	list, err := clientset.AppsV1().Deployments(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			if options.UserRole.HasDeploymentAccess(item.GetLabels()) {
				deployments = append(deployments, DeploymentOverview{
					Name:                 item.GetName(),
					Namespace:            item.GetNamespace(),
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
