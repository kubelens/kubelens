package k8sv1

import (
	"context"
	"fmt"

	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapOptions .
type ConfigMapOptions struct {
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

type ConfigMapOverview struct {
	/// the name
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace
	Namespace string `json:"namespace"`
	// the full configmap
	ConfigMap *v1.ConfigMap `json:"configMap,omitempty"`
}

// ConfigMap returns a configmap given filter options
func (k *Client) ConfigMap(options ConfigMapOptions) (overview *ConfigMapOverview, apiErr *errs.APIError) {

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.CoreV1().ConfigMaps(options.Namespace).List(options.Context, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", options.Name),
	})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			return &ConfigMapOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				ConfigMap:  &item,
			}, nil
		}
	}
	return overview, nil
}

// ConfigMaps returns a list ofconfigmaps given filter options
func (k *Client) ConfigMaps(options ConfigMapOptions) (overviews []ConfigMapOverview, apiErr *errs.APIError) {
	overviews = []ConfigMapOverview{}
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	cml := clientset.CoreV1().ConfigMaps(options.Namespace)

	if cml == nil {
		return nil, nil
	}

	lo := metav1.ListOptions{}
	if len(options.LinkedName) > 0 {
		lo.LabelSelector = generateLabelSelector(options.LinkedName)
	}
	list, err := cml.List(options.Context, lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			overviews = append(overviews, ConfigMapOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				ConfigMap:  &item,
			})
		}
	}
	return overviews, nil
}
