package k8sv1

import (
	"context"
	"fmt"

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

	list, err := dpl.List(options.Context, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", options.Name),
	})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			return &DeploymentOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				Deployment: &item,
			}, nil
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

	lo := metav1.ListOptions{}
	if len(options.LinkedName) > 0 {
		lo.LabelSelector = generateLabelSelector(options.LinkedName)
	}
	list, err := dpl.List(options.Context, lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			overviews = append(overviews, DeploymentOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				Deployment: &item,
			})
		}
	}

	return overviews, nil
}
