package k8sv1

import (
	"context"
	"strings"
	"sync"

	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicaSetOptions contains fields used for filtering when retrieving daemon sets
type ReplicaSetOptions struct {
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

type ReplicaSetOverview struct {
	/// the name
	Name string `json:"name"`
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// the namespace
	Namespace string `json:"namespace"`
	// the full ReplicaSet
	ReplicaSet *appsv1.ReplicaSet `json:"replicaSet,omitempty"`
}

// ReplicaSet returns a ReplicaSet given filter options
func (k *Client) ReplicaSet(options ReplicaSetOptions) (overview *ReplicaSetOverview, apiErr *errs.APIError) {

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	rsl := clientset.AppsV1().ReplicaSets(options.Namespace)

	if rsl == nil {
		return nil, nil
	}

	list, err := rsl.List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	if list != nil && len(list.Items) > 0 {
		for i, item := range list.Items {
			go func(index int, rs appsv1.ReplicaSet) {
				defer wg.Done()
				// first check name of deployment, then by labelSelctor
				if strings.EqualFold(rs.Name, options.Name) {
					overview = &ReplicaSetOverview{
						Name:       rs.Name,
						LinkedName: getLinkedName(rs.Labels),
						Namespace:  rs.Namespace,
						ReplicaSet: &rs,
					}
				}
			}(i, item)
		}
	}

	wg.Wait()

	return overview, nil
}

// ReplicaSets returns a list of ReplicaSets given filter options
func (k *Client) ReplicaSets(options ReplicaSetOptions) (overviews []ReplicaSetOverview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	rsl := clientset.AppsV1().ReplicaSets(options.Namespace)

	if rsl == nil {
		return nil, nil
	}

	list, err := rsl.List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	var ovs []ReplicaSetOverview

	if list != nil && len(list.Items) > 0 {
		ovs = make([]ReplicaSetOverview, len(list.Items))

		for i, item := range list.Items {
			go func(index int, rs appsv1.ReplicaSet) {
				defer wg.Done()
				if len(options.LinkedName) > 0 {
					if labelsContainSelector(options.LinkedName, rs.Labels) {
						ovs[index] = ReplicaSetOverview{
							Name:       rs.Name,
							LinkedName: getLinkedName(rs.Labels),
							Namespace:  rs.Namespace,
							ReplicaSet: &rs,
						}
					}
				} else {
					ovs[index] = ReplicaSetOverview{
						Name:       rs.Name,
						LinkedName: getLinkedName(rs.Labels),
						Namespace:  rs.Namespace,
						ReplicaSet: &rs,
					}
				}
			}(i, item)
		}
	}

	wg.Wait()

	if ovs != nil && len(ovs) > 0 {
		overviews = []ReplicaSetOverview{}
		for _, ov := range ovs {
			if len(ov.Name) > 0 && len(ov.LinkedName) > 0 {
				overviews = append(overviews, ov)
			}
		}
	}

	return overviews, nil
}
