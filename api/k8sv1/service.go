package k8v1

import (
	"sync"

	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceOverviews returns a list of services given filter options
func (k *Client) ServiceOverviews(options ServiceOptions) (svco []ServiceOverview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	lo := metav1.ListOptions{
		IncludeUninitialized: true,
	}

	if len(options.LabelSearch) > 0 {
		lo.LabelSelector = options.LabelSearch
	}

	list, err := clientset.CoreV1().Services(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	svco = make([]ServiceOverview, len(list.Items))
	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	for i, item := range list.Items {
		// check access at label level
		if options.UserCanAccess(item.GetLabels()) {
			name, labelKey := getAppName(
				item.GetLabels(),
				"",
				getDefaultSearchLabel(item.Spec.Selector),
				item.GetName(),
			)

			go func(index int, service v1.Service) {
				defer wg.Done()

				svc := ServiceOverview{
					AppName: Name{
						LabelKey: labelKey,
						Value:    name,
					},
					DeployerLink: getDeployerLink(service.GetName()),
					Name:         service.GetName(),
					Namespace:    service.GetNamespace(),
				}

				if options.Detailed {
					svc.AddDetail(&service.Spec, &service.Status)
				}

				svco[index] = svc
			}(i, item)
		} else {
			wg.Done()
		}
	}

	wg.Wait()

	return svco, nil
}
