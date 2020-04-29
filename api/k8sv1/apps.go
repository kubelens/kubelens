package k8sv1

import (
	"sync"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	klog "github.com/kubelens/kubelens/api/log"
)

// AppOptions .
type AppOptions struct {
	// user roles
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// App represents an application within Kubernetes.
type App struct {
	// name of the application
	Name string `json:"name"`
	// the namespace of the app
	Namespace string `json:"namespace"`
	// kind of application, e.g. Service, DaemonSet
	Kind string `json:"kind"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// deployer link if any
	DeployerLink string `json:"deployerLink,omitempty"`
}

// AppOverviewOptions .
type AppOverviewOptions struct {
	// name of the  application
	AppName string `json:"appname"`
	// namespace of the app
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// include detail
	Detailed bool `json:"detailed"`
	// user roles
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// AppOverview .
type AppOverview struct {
	PodOverviews        PodOverview          `json:"podOverviews,omitempty"`
	ServiceOverviews    []ServiceOverview    `json:"serviceOverviews,omitempty"`
	DaemonSetOverviews  []DaemonSetOverview  `json:"daemonSetOverviews,omitempty"`
	JobOverviews        []JobOverview        `json:"jobOverviews,omitempty"`
	ReplicaSetOverviews []ReplicaSetOverview `json:"replicaSetOverviews,omitempty"`
}

// AppInfo .
type AppInfo struct {
	// friendly name of app
	FriendlyName string
	// actual name of app
	Name string
	// namspace the app is in
	Namespace string
	// the kind of app, e.g. Service, DaemonSet, etc.
	Kind string
	// the LabelSelector of the object.
	LabelSelector map[string]string
}

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
		appInfos := []*[]AppInfo{}

		wg := sync.WaitGroup{}
		// adding 3x since we are using 4 go routines per iteration
		wg.Add(len(namespaces.Items) * 4)

		errors := []*errs.APIError{}

		for i, item := range namespaces.Items {
			ns := item.GetName()

			// get services
			go func(index int, namespace string) {
				defer wg.Done()

				svcInfos, err := k.ServiceAppInfos(ServiceOptions{Namespace: namespace})

				if err != nil {
					klog.Trace()
					errors = append(errors, err)
				}

				appInfos = append(appInfos, &svcInfos)
			}(i, ns)

			// get daemonsets
			go func(index int, namespace string) {
				defer wg.Done()

				dsInfos, err := k.DaemonSetAppInfos(DaemonSetOptions{Namespace: namespace})

				if err != nil {
					klog.Trace()
					errors = append(errors, err)
				}

				appInfos = append(appInfos, &dsInfos)
			}(i, ns)

			// get jobs
			go func(index int, namespace string) {
				defer wg.Done()

				jobInfos, err := k.JobAppInfos(JobOptions{Namespace: namespace})

				if err != nil {
					klog.Trace()
					errors = append(errors, err)
				}

				appInfos = append(appInfos, &jobInfos)
			}(i, ns)

			// get replicasets
			go func(index int, namespace string) {
				defer wg.Done()

				jobInfos, err := k.ReplicaSetAppInfos(ReplicaSetOptions{Namespace: namespace})

				if err != nil {
					klog.Trace()
					errors = append(errors, err)
				}

				appInfos = append(appInfos, &jobInfos)
			}(i, ns)
		}

		wg.Wait()

		if len(errors) > 0 {
			return apps, errs.ListToInternalServerError(errors)
		}

		for _, a := range appInfos {
			for _, b := range *a {
				apps = append(apps, App{
					Name:          b.FriendlyName,
					Namespace:     b.Namespace,
					Kind:          b.Kind,
					LabelSelector: b.LabelSelector,
					DeployerLink:  getDeployerLink(b.Name),
				})
			}
		}
	}

	return apps, nil
}

// AppOverview returns an list of application overviews with high level info such as pods, services, deployments, etc.
func (k *Client) AppOverview(options AppOverviewOptions) (ao *AppOverview, apiErr *errs.APIError) {
	var po *PodOverview
	var sos []ServiceOverview
	var dsos []DaemonSetOverview
	var jsos []JobOverview
	var rsos []ReplicaSetOverview
	errors := []*errs.APIError{}

	wg := sync.WaitGroup{}

	wg.Add(5)

	go func() {
		defer wg.Done()
		var err *errs.APIError
		po, err = k.PodOverview(PodOverviewOptions{
			AppName:       options.AppName,
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer wg.Done()
		var err *errs.APIError
		sos, err = k.ServiceOverviews(ServiceOptions{
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
			Detailed:      options.Detailed,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer wg.Done()
		var err *errs.APIError
		dsos, err = k.DaemonSetOverviews(DaemonSetOptions{
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer wg.Done()
		var err *errs.APIError
		jsos, err = k.JobOverviews(JobOptions{
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer wg.Done()
		var err *errs.APIError
		rsos, err = k.ReplicaSetOverviews(ReplicaSetOptions{
			Namespace:     options.Namespace,
			LabelSelector: options.LabelSelector,
			UserRole:      options.UserRole,
			Logger:        options.Logger,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	if len(errors) > 0 {
		return nil, errs.ListToInternalServerError(errors)
	}

	return &AppOverview{
		PodOverviews:        *po,
		ServiceOverviews:    sos,
		DaemonSetOverviews:  dsos,
		JobOverviews:        jsos,
		ReplicaSetOverviews: rsos,
	}, nil
}
