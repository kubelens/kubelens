package k8sv1

import (
	"strings"
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

	if len(options.LabelSelector) > 0 {
		lo.LabelSelector = toLabelSelectorString(options.LabelSelector)
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
			go func(index int, service v1.Service) {
				defer wg.Done()

				svc := ServiceOverview{
					FriendlyName: getFriendlyAppName(
						item.GetLabels(),
						item.GetName(),
					),
					Selector:     service.Spec.Selector,
					DeployerLink: getDeployerLink(service.GetName()),
					Name:         service.GetName(),
					Namespace:    service.GetNamespace(),
				}

				if options.Detailed {
					svc.AddDetail(&service.Spec, &service.Status)
				}

				// add deployments per service for ease of display by client since deployments are really
				// specific to certian K8s kinds.
				if options.Detailed && options.UserRole.HasDeploymentAccess(item.GetLabels()) {
					deployments, err := k.DeploymentOverviews(DeploymentOptions{
						LabelSelector: svc.Selector,
						Namespace:     service.GetNamespace(),
						UserRole:      options.UserRole,
						Logger:        options.Logger,
					})
					// just trace the error and move on, shouldn't be critical.
					if err != nil {
						klog.Trace()
					}
					svc.AddDeploymentOverviews(deployments)
				}

				// get this list of config maps here to ensure correct namespace. This should perform
				// just fine, but if not, this list could be generated outsid of this routine.
				cmList, err := clientset.CoreV1().ConfigMaps(item.GetNamespace()).List(lo)
				// just trace the error and move on, shouldn't be critical.
				if err != nil {
					klog.Trace()
				}

				if cmList != nil && options.UserRole.HasConfigMapAccess(item.GetLabels()) {
					serviceConfigMaps := []v1.ConfigMap{}

					for _, cm := range cmList.Items {
						for cmLblKey, cmLblValue := range cm.GetLabels() {
							// Only look for a match on the selector lables. Not sure of a better way
							// to ensure that the found config map is the one tied to the service.
							if strings.EqualFold(item.Spec.Selector[cmLblKey], cmLblValue) {
								serviceConfigMaps = append(serviceConfigMaps, cm)
								// break out if found since there could be multiple
								// selector labels to match on. one find should be unique enough?
								break
							}
						}
					}

					if options.Detailed && len(serviceConfigMaps) > 0 {
						svc.AddConfigMaps(&serviceConfigMaps)
					}
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

// ServiceAppInfos returns basic info for all services found for a given namespace.
func (k *Client) ServiceAppInfos(options ServiceOptions) (info []AppInfo, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.CoreV1().Services(options.Namespace).List(metav1.ListOptions{IncludeUninitialized: true})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	errors := []*errs.APIError{}

	if list != nil {
		info = make([]AppInfo, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			go func(index int, svc v1.Service) {
				defer wg.Done()

				name := getFriendlyAppName(
					svc.GetLabels(),
					svc.GetName(),
				)

				info[index] = AppInfo{
					FriendlyName:  name,
					Name:          svc.GetName(),
					Namespace:     svc.GetNamespace(),
					Kind:          "Service",
					LabelSelector: svc.Spec.Selector,
				}
			}(i, item)
		}
		wg.Wait()
	}

	if len(errors) > 0 {
		if len(errors) > 0 {
			return info, errs.ListToInternalServerError(errors)
		}
	}

	return info, nil
}
