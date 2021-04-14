package k8sv1

import (
	"context"
	"strings"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSetOptions contains fields used for filtering when retrieving daemon sets
type DaemonSetOptions struct {
	// the name of the daemonSet
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
	// Context .
	Context context.Context
}

// DaemonSetOverview .
type DaemonSetOverview struct {
	// the name
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace
	Namespace string `json:"namespace"`
	// the full daemonset
	DaemonSet *appsv1.DaemonSet `json:"daemonSet,omitempty"`
}

// DaemonSet returns a daemonsets given filter options
func (k *Client) DaemonSet(options DaemonSetOptions) (overview *DaemonSetOverview, apiErr *errs.APIError) {

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.AppsV1().DaemonSets(options.Namespace).List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			if options.UserRole.HasDaemonSetAccess(item.Labels) {
				// first check name of deployment, then by labelSelctor
				if strings.EqualFold(item.Name, options.Name) {
					return &DaemonSetOverview{
						Name:       item.Name,
						LinkedName: getLinkedName(item.Labels),
						Namespace:  item.Namespace,
						DaemonSet:  &item,
					}, nil
				}
			}
		}
	}
	return overview, nil
}

// DaemonSet returns a daemonsets given filter options
func (k *Client) DaemonSets(options DaemonSetOptions) (overviews []DaemonSetOverview, apiErr *errs.APIError) {
	overviews = []DaemonSetOverview{}
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.AppsV1().DaemonSets(options.Namespace).List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			if options.UserRole.HasDaemonSetAccess(item.Labels) {
				if len(options.LinkedName) > 0 {
					if labelsContainSelector(options.LinkedName, item.Labels) {
						overviews = append(overviews, DaemonSetOverview{
							Name:       item.Name,
							LinkedName: getLinkedName(item.Labels),
							Namespace:  item.Namespace,
							DaemonSet:  &item,
						})
					}
				} else {
					overviews = append(overviews, DaemonSetOverview{
						Name:       item.Name,
						LinkedName: getLinkedName(item.Labels),
						Namespace:  item.Namespace,
						DaemonSet:  &item,
					})
				}
			}
		}
	}
	return overviews, nil
}
