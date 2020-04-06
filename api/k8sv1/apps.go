package k8v1

import (
	"fmt"
	"sync"

	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
)

// AppOverview returns an list of application overviews with high level info such as pods, services, deployments, etc.
func (k *Client) AppOverview(options AppOverviewOptions) (ao *AppOverview, apiErr *errs.APIError) {
	var po *PodOverview
	var sos []ServiceOverview
	var err1, err2 *errs.APIError

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		po, err1 = k.PodOverview(PodOverviewOptions{
			AppName:         options.AppName,
			Namespace:       options.Namespace,
			AppNameLabelKey: options.LabelKey,
			UserRole:        options.UserRole,
			Logger:          options.Logger,
		})
	}()

	go func() {
		defer wg.Done()
		sos, err2 = k.ServiceOverviews(ServiceOptions{
			Namespace:     options.Namespace,
			LabelSelector: fmt.Sprintf("%s=%s", options.LabelKey, options.AppName),
			Detailed:      options.Detailed,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
	}()

	wg.Wait()

	if err1 != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err1.Message)
	}

	if err2 != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err2.Message)
	}

	return &AppOverview{
		PodOverviews:     *po,
		ServiceOverviews: sos,
	}, nil
}
