package k8sv1

import (
	"reflect"
	"sync"

	v1 "k8s.io/api/core/v1"

	"github.com/kubelens/kubelens/api/errs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	klog "github.com/kubelens/kubelens/api/log"
)

// Apps returns a list of apps running in kubernetes, determined by searching deployments for each namespace.
// For each namespace, Kubernetes Kinds are searched for the type of application, e.g. Service, DaemonSet, etc.
func (k *Client) Apps(options AppOptions) (apps []*App, apiErr *errs.APIError) {
	apps = []*App{}
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return apps, errs.InternalServerError(err.Error())
	}

	lo := metav1.ListOptions{
		IncludeUninitialized: true,
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(lo)

	if err != nil {
		klog.Trace()
		return apps, errs.InternalServerError(err.Error())
	}

	if namespaces != nil {
		wg := sync.WaitGroup{}
		wg.Add(len(namespaces.Items))

		for _, item := range namespaces.Items {
			go func(ns v1.Namespace) {
				defer wg.Done()

				dos, err := k.DeploymentOverviews(DeploymentOptions{
					Namespace: ns.GetName(),
					UserRole:  options.UserRole,
					Logger:    options.Logger,
				})

				if err != nil {
					klog.Trace()
				}

				wginner := sync.WaitGroup{}
				wginner.Add(len(dos))

				for _, do := range dos {
					go func(dep DeploymentOverview) {
						defer wginner.Done()

						kind := k.getAppKind(dep.LabelSelector, dep.Namespace)

						if len(kind) > 0 {
							apps = append(apps, &App{
								Name:          dep.FriendlyName,
								Namespace:     dep.Namespace,
								Kind:          kind,
								LabelSelector: dep.LabelSelector,
								DeployerLink:  getDeployerLink(dep.Name),
							})
						}
					}(do)
				}
				wginner.Wait()
			}(item)
		}
		wg.Wait()
	}

	return apps, nil
}

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
			AppName:       options.AppName,
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
	}()

	go func() {
		defer wg.Done()
		sos, err2 = k.ServiceOverviews(ServiceOptions{
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
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

// getAppKind checks for Service, DaemonSet, by LabelSelector and namespace.
// The Kubernetes Kind is returned if found.
func (k *Client) getAppKind(labelSelector map[string]string, namespace string) (kind string) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return ""
	}

	opts := metav1.ListOptions{
		IncludeUninitialized: true,
	}

	svcList, err := clientset.CoreV1().Services(namespace).List(opts)

	if err != nil {
		klog.Trace()
	}

	var svcfound string

	if svcList != nil && len(svcList.Items) > 0 {
		wg1 := sync.WaitGroup{}
		wg1.Add(len(svcList.Items))
		for _, item := range svcList.Items {
			go func(svc v1.Service) {
				defer wg1.Done()
				if reflect.DeepEqual(svc.Spec.Selector, labelSelector) {
					svcfound = "Service"
				}
			}(item)
		}
		wg1.Wait()
	}

	if len(svcfound) > 0 {
		return svcfound
	}

	dsList, err := clientset.AppsV1().DaemonSets(namespace).List(opts)
	if err != nil {
		klog.Trace()
	}

	var dsfound string

	if dsList != nil && len(dsList.Items) > 0 {
		for _, item := range dsList.Items {
			if reflect.DeepEqual(item.Spec.Selector, labelSelector) {
				dsfound = "DaemonSet"
				break
			}
		}
	}

	if len(dsfound) > 0 {
		return dsfound
	}

	return ""
}
