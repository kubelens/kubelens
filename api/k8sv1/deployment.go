package k8sv1

import (
	"context"
	"strings"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	klog "github.com/kubelens/kubelens/api/log"
)

// DeploymentOptions contains fields used for filtering when retrieving deployments
type DeploymentOptions struct {
	// the name of the deployment
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
	// Context .
	Context context.Context
}

type DeploymentOverview struct {
	// the name of the deployment
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace of the deployment
	Namespace string `json:"namespace"`
	// deployment labels
	Deployment *appsv1.Deployment `json:"deployment,omitempty"`
}

func (k *Client) Deployment(options DeploymentOptions) (overview *DeploymentOverview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	dpl := clientset.AppsV1().Deployments(options.Namespace)

	if dpl == nil {
		return nil, nil
	}

	list, err := dpl.List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			if options.UserRole.HasDeploymentAccess(item.Labels) {
				// first check name of deployment, then by labelSelctor
				if strings.EqualFold(item.Name, options.Name) {
					return &DeploymentOverview{
						Name:       item.Name,
						LinkedName: getLinkedName(item.Labels),
						Namespace:  item.Namespace,
						Deployment: &item,
					}, nil
				}
			}
		}
	}

	return overview, nil
}

// Deployments retrieves all deployments by namespace.
func (k *Client) Deployments(options DeploymentOptions) (overviews []DeploymentOverview, apiErr *errs.APIError) {
	overviews = []DeploymentOverview{}
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	dpl := clientset.AppsV1().Deployments(options.Namespace)

	if dpl == nil {
		return nil, nil
	}

	list, err := dpl.List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			if options.UserRole.HasDeploymentAccess(item.Labels) {
				if len(options.LinkedName) > 0 {
					if labelsContainSelector(options.LinkedName, item.Labels) {
						overviews = append(overviews, DeploymentOverview{
							Name:       item.Name,
							LinkedName: getLinkedName(item.Labels),
							Namespace:  item.Namespace,
							Deployment: &item,
						})
					}
				} else {
					overviews = append(overviews, DeploymentOverview{
						Name:       item.Name,
						LinkedName: getLinkedName(item.Labels),
						Namespace:  item.Namespace,
						Deployment: &item,
					})
				}
			}
		}
	}

	return overviews, nil
}
