package k8sv1

import (
	"strings"
	"sync"

	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSetOverviews returns a list of daemonsets given filter options
func (k *Client) DaemonSetOverviews(options DaemonSetOptions) (daemonsets []DaemonSetOverview, apiErr *errs.APIError) {
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

	list, err := clientset.AppsV1().DaemonSets(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	if list != nil && len(list.Items) > 0 {
		daemonsets = make([]DaemonSetOverview, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			if options.UserRole.HasDaemonSetAccess(item.GetLabels()) {
				go func(index int, ds appsv1.DaemonSet) {
					defer wg.Done()

					var labelSelector map[string]string
					// this shouldn't be null, but default to regular labels if it is.
					if ds.Spec.Selector != nil {
						labelSelector = ds.Spec.Selector.MatchLabels
					} else if len(options.LabelSelector) > 0 {
						labelSelector = options.LabelSelector
					} else {
						labelSelector = ds.GetLabels()
					}

					name := getFriendlyAppName(
						ds.GetLabels(),
						ds.GetName(),
					)

					dso := DaemonSetOverview{
						FriendlyName:           name,
						Name:                   ds.GetName(),
						Namespace:              ds.GetNamespace(),
						LabelSelector:          labelSelector,
						CurrentNumberScheduled: int(ds.Status.CurrentNumberScheduled),
						DesiredNumberScheduled: int(ds.Status.DesiredNumberScheduled),
						NumberAvailable:        int(ds.Status.NumberAvailable),
						NumberMisscheduled:     int(ds.Status.NumberMisscheduled),
						NumberReady:            int(ds.Status.NumberReady),
						NumberUnavailable:      int(ds.Status.NumberUnavailable),
						UpdatedNumberScheduled: int(ds.Status.UpdatedNumberScheduled),
						Conditions:             ds.Status.Conditions,
					}

					// add deployments per daemon set for ease of display by client since deployments are really
					// specific to certian K8s kinds.
					if options.UserRole.HasDeploymentAccess(ds.GetLabels()) {
						deployments, err := k.DeploymentOverviews(DeploymentOptions{
							LabelSelector: labelSelector,
							Namespace:     ds.GetNamespace(),
							UserRole:      options.UserRole,
							Logger:        options.Logger,
						})
						// just trace the error and move on, shouldn't be critical.
						if err != nil {
							klog.Trace()
						}

						if len(deployments) > 0 {
							dso.AddDeploymentOverviews(&deployments)
						}
					}

					// get this list of config maps here to ensure correct namespace. This should perform
					// just fine, but if not, this list could be generated outsid of this routine.
					cmList, err := clientset.CoreV1().ConfigMaps(item.GetNamespace()).List(lo)
					// just trace the error and move on, shouldn't be critical.
					if err != nil {
						klog.Trace()
					}

					if cmList != nil && options.UserRole.HasConfigMapAccess(item.GetLabels()) {
						dsConfigMaps := []v1.ConfigMap{}

						for _, cm := range cmList.Items {
							for cmLblKey, cmLblValue := range cm.GetLabels() {
								// Only look for a match on the selector lables. Not sure of a better way
								// to ensure that the found config map is the one tied to the service.
								if strings.EqualFold(labelSelector[cmLblKey], cmLblValue) {
									dsConfigMaps = append(dsConfigMaps, cm)
									// break out if found since there could be multiple
									// selector labels to match on. one find should be unique enough?
									break
								}
							}
						}

						if len(dsConfigMaps) > 0 {
							dso.AddConfigMaps(&dsConfigMaps)
						}
					}

					daemonsets[index] = dso
				}(i, item)
			}
		}

		wg.Wait()
	}

	return daemonsets, nil
}

// DaemonSetAppInfos returns basic info for all daemon sets found for a given namespace.
func (k *Client) DaemonSetAppInfos(options DaemonSetOptions) (info []AppInfo, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.AppsV1().DaemonSets(options.Namespace).List(metav1.ListOptions{IncludeUninitialized: true})

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
			go func(index int, ds appsv1.DaemonSet) {
				defer wg.Done()
				if err != nil {
					klog.Trace()
					errors = append(errors, errs.InternalServerError(err.Error()))
				}

				var labelSelector map[string]string
				// this shouldn't be null, but default to regular labels if it is.
				if ds.Spec.Selector != nil {
					labelSelector = ds.Spec.Selector.MatchLabels
				} else {
					labelSelector = ds.GetLabels()
				}

				name := getFriendlyAppName(
					ds.GetLabels(),
					ds.GetName(),
				)

				info[index] = AppInfo{
					FriendlyName:  name,
					Name:          ds.GetName(),
					Namespace:     ds.GetNamespace(),
					Kind:          "DaemonSet",
					LabelSelector: labelSelector,
				}
			}(i, item)
		}
		wg.Wait()
	}

	if len(errors) > 0 {
		return info, errs.ListToInternalServerError(errors)
	}

	return info, nil
}
