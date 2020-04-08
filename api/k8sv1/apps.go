package k8sv1

import (
	"sync"

	"github.com/kubelens/kubelens/api/errs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	klog "github.com/kubelens/kubelens/api/log"
)

// Apps returns a list of apps running in kubernetes, determined by searching deployments for each namespace.
// For each namespace, Kubernetes Kinds are searched for the type of application, e.g. Service, DaemonSet, etc.
func (k *Client) Apps(options AppOptions) (apps []App, apiErr *errs.APIError) {
	apps = []App{}
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
		for _, item := range namespaces.Items {
			dos, err := k.DeploymentOverviews(DeploymentOptions{
				Namespace: item.GetName(),
				UserRole:  options.UserRole,
				Logger:    options.Logger,
			})

			if err != nil {
				klog.Trace()
				return apps, errs.InternalServerError(err.Message)
			}

			for _, do := range dos {
				kind := k.getAppKind(do.LabelSelector, do.Namespace)

				if len(kind) > 0 {
					apps = append(apps, App{
						Name:          do.FriendlyName,
						Namespace:     do.Namespace,
						Kind:          kind,
						LabelSelector: do.LabelSelector,
						DeployerLink:  getDeployerLink(do.Name),
					})
				}
			}
		}
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

// getAppKind checks for Service, DaemonSet, by name and namespace.
// The Kubernetes Kind is returned if found.
func (k *Client) getAppKind(labelSelector map[string]string, namespace string) (kind string) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return ""
	}

	opts := metav1.ListOptions{
		LabelSelector:        toLabelSelectorString(labelSelector),
		IncludeUninitialized: true,
	}

	svcList, err := clientset.CoreV1().Services(namespace).List(opts)

	if err != nil {
		klog.Trace()
	}

	if svcList != nil && len(svcList.Items) > 0 {
		return "Service"
	}

	dsList, err := clientset.AppsV1().DaemonSets(namespace).List(opts)
	if err != nil {
		klog.Trace()
	}

	if dsList != nil && len(dsList.Items) > 0 {
		return "DaemonSet"
	}

	return ""
}
