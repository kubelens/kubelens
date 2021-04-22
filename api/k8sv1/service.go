package k8sv1

import (
	"context"
	"fmt"

	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceOptions contains fields used for filtering when retrieving services
type ServiceOptions struct {
	/// the name (optional)
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the labels to match kinds
	Labels map[string]string `json:"labels"`
	// logger instance
	Logger klog.Logger
	// Context .
	Context context.Context
}

type ServiceOverview struct {
	/// the name
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace
	Namespace string `json:"namespace"`
	// the full Service
	Service *v1.Service `json:"service,omitempty"`
}

// Service returns a Service given filter options
func (k *Client) Service(options ServiceOptions) (overview *ServiceOverview, apiErr *errs.APIError) {

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	svcs := clientset.CoreV1().Services(options.Namespace)

	if svcs == nil {
		return nil, nil
	}

	list, err := svcs.List(options.Context, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", options.Name),
	})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			return &ServiceOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				Service:    &item,
			}, nil
		}
	}
	return overview, nil
}

// Services returns a list ofServices given filter options
func (k *Client) Services(options ServiceOptions) (overviews []ServiceOverview, apiErr *errs.APIError) {
	overviews = []ServiceOverview{}
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	svcs := clientset.CoreV1().Services(options.Namespace)

	if svcs == nil {
		return nil, nil
	}

	lo := metav1.ListOptions{}
	if len(options.LinkedName) > 0 {
		lo.LabelSelector = generateLabelSelector(options.LinkedName)
	}
	list, err := svcs.List(options.Context, lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			overviews = append(overviews, ServiceOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				Service:    &item,
			})
		}
	}
	return overviews, nil
}
