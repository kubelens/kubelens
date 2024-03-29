package k8sv1

import (
	"context"
	"fmt"

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

	dsl := clientset.AppsV1().DaemonSets(options.Namespace)

	if dsl == nil {
		return nil, nil
	}
	list, err := dsl.List(options.Context, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", options.Name),
	})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			return &DaemonSetOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				DaemonSet:  &item,
			}, nil
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

	dsl := clientset.AppsV1().DaemonSets(options.Namespace)

	if dsl == nil {
		return nil, nil
	}

	lo := metav1.ListOptions{}
	if len(options.LinkedName) > 0 {
		lo.LabelSelector = generateLabelSelector(options.LinkedName)
	}
	list, err := dsl.List(options.Context, lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			overviews = append(overviews, DaemonSetOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				DaemonSet:  &item,
			})
		}
	}
	return overviews, nil
}
