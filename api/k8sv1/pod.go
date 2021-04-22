package k8sv1

import (
	"context"
	"fmt"
	"sync"

	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodOptions contains fields used for filtering when retrieving application overiew(s).
type PodOptions struct {
	// name of the pod
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	// Limit the number of pod summaries to return
	// Use function GetLimit to get the default limit or this overriden value.
	Limit int64 `json:"linit"`
	// logger instance
	Logger klog.Logger
	// Context .
	Context context.Context
}

// GetLimit returns Limit || default
func (a *PodOptions) GetLimit() int64 {
	if a.Limit > 0 {
		return a.Limit
	}
	return 100
}

type PodOverview struct {
	/// the name
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace
	Namespace string `json:"namespace"`
	// the full Pod
	Pod *v1.Pod `json:"pod,omitempty"`
}

// Pod returns a Pod given filter options
func (k *Client) Pod(options PodOptions) (overview *PodOverview, apiErr *errs.APIError) {

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	pds := clientset.CoreV1().Pods(options.Namespace)

	if pds == nil {
		return nil, nil
	}

	list, err := pds.List(options.Context, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", options.Name),
	})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	if list != nil && len(list.Items) > 0 {
		for i, item := range list.Items {
			go func(index int, pod v1.Pod) {
				defer wg.Done()

				// try to catch sensitive values in environment variables.
				// TODO is there a better way?
				for ci, c := range pod.Spec.Containers {
					env := []v1.EnvVar{}

					for _, e := range c.Env {
						if !stringContainsSensitiveInfo(e.Name) {
							env = append(env, e)
						}
					}
					pod.Spec.Containers[ci].Env = env
				}

				overview = &PodOverview{
					Name:       pod.Name,
					LinkedName: getLinkedName(pod.Labels),
					Namespace:  pod.Namespace,
					Pod:        &pod,
				}
			}(i, item)
		}
	}

	wg.Wait()

	return overview, nil
}

// Pods returns a list ofPods given filter options
func (k *Client) Pods(options PodOptions) (overviews []PodOverview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	pds := clientset.CoreV1().Pods(options.Namespace)

	if pds == nil {
		return nil, nil
	}

	lo := metav1.ListOptions{}
	if len(options.LinkedName) > 0 {
		lo.LabelSelector = generateLabelSelector(options.LinkedName)
	}
	list, err := pds.List(options.Context, lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	if list != nil && len(list.Items) > 0 {
		overviews = make([]PodOverview, len(list.Items))

		for i, item := range list.Items {
			go func(index int, pod v1.Pod) {
				defer wg.Done()
				overviews[index] = PodOverview{
					Name:       pod.Name,
					LinkedName: getLinkedName(pod.Labels),
					Namespace:  pod.Namespace,
					Pod:        &pod,
				}
			}(i, item)
		}
	}

	wg.Wait()

	return overviews, nil
}
