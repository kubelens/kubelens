package k8sv1

import (
	"context"
	"fmt"

	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobOptions contains fields used for filtering when retrieving jobs
type JobOptions struct {
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

type JobOverview struct {
	/// the name
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace
	Namespace string `json:"namespace"`
	// the full configmap
	Job *batchv1.Job `json:"job,omitempty"`
}

// Job returns a Job given filter options
func (k *Client) Job(options JobOptions) (overview *JobOverview, apiErr *errs.APIError) {

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	jbs := clientset.BatchV1().Jobs(options.Namespace)

	if jbs == nil {
		return nil, nil
	}

	list, err := jbs.List(options.Context, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", options.Name),
	})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			return &JobOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				Job:        &item,
			}, nil
		}
	}
	return overview, nil
}

// Jobs returns a list ofJobs given filter options
func (k *Client) Jobs(options JobOptions) (overviews []JobOverview, apiErr *errs.APIError) {
	overviews = []JobOverview{}
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	jbs := clientset.BatchV1().Jobs(options.Namespace)

	if jbs == nil {
		return nil, nil
	}

	lo := metav1.ListOptions{}
	if len(options.LinkedName) > 0 {
		lo.LabelSelector = generateLabelSelector(options.LinkedName)
	}
	list, err := jbs.List(options.Context, lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		for _, item := range list.Items {
			overviews = append(overviews, JobOverview{
				Name:       item.Name,
				LinkedName: getLinkedName(item.Labels),
				Namespace:  item.Namespace,
				Job:        &item,
			})
		}
	}
	return overviews, nil
}
