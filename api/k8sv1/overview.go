package k8sv1

import (
	"context"
	"strings"
	"sync"

	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OverviewOptions contains fields used for filtering when retrieving application overiew(s).
type OverviewOptions struct {
	// the value from the label "app=NAME", corresponds to config.LabelKeyLink
	LinkedName string `json:"linkedName"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	// logger instance
	Logger klog.Logger
	// Context .
	Context context.Context
}

type Overview struct {
	LinkedName  string               `json:"linkedName,omitempty"`
	Namespace   string               `json:"namespace,omitempty"`
	DaemonSets  []DaemonSetOverview  `json:"daemonSets,omitempty"`
	Deployments []DeploymentOverview `json:"deployments,omitempty"`
	Jobs        []JobOverview        `json:"jobs,omitempty"`
	Pods        []PodOverview        `json:"pods,omitempty"`
	ReplicaSets []ReplicaSetOverview `json:"replicaSets,omitempty"`
	Services    []ServiceOverview    `json:"services,omitempty"`
	ConfigMaps  []ConfigMapOverview  `json:"configMaps,omitempty"`
}

// Overview returns a Overview given filter options
func (k *Client) Overview(options OverviewOptions) (overview *Overview, apiErr *errs.APIError) {
	// DaemonSets
	dss, _ := k.DaemonSets(DaemonSetOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})
	// Deployments
	dps, _ := k.Deployments(DeploymentOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})
	// Jobs
	jbs, _ := k.Jobs(JobOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})
	// Pods
	povs, _ := k.Pods(PodOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})
	// ReplicaSets
	rss, _ := k.ReplicaSets(ReplicaSetOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})
	// Services
	svcs, _ := k.Services(ServiceOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})

	// ConfigMaps
	cms, _ := k.ConfigMaps(ConfigMapOptions{
		Namespace:  options.Namespace,
		LinkedName: options.LinkedName,
		Logger:     options.Logger,
		Context:    options.Context,
	})

	overview = &Overview{
		LinkedName:  options.LinkedName,
		Namespace:   options.Namespace,
		DaemonSets:  dss,
		Deployments: dps,
		Jobs:        jbs,
		Pods:        povs,
		ReplicaSets: rss,
		Services:    svcs,
		ConfigMaps:  cms,
	}

	return overview, nil
}

// Pods returns a list ofPods given filter options
func (k *Client) Overviews(options OverviewOptions) (overviews []Overview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(options.Context, metav1.ListOptions{})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(len(namespaces.Items))

	nsOverviews := make([][]Overview, len(namespaces.Items))

	for i, namespace := range namespaces.Items {
		nsOverviews[i] = []Overview{}

		go func(index int, ns v1.Namespace) {
			defer wg.Done()
			// DaemonSets
			dss, _ := k.DaemonSets(DaemonSetOptions{
				Namespace: ns.Namespace,
				Logger:    options.Logger,
				Context:   options.Context,
			})

			for _, ds := range dss {
				nsOverviews[index] = append(nsOverviews[index], Overview{
					LinkedName: ds.LinkedName,
					Namespace:  ds.Namespace,
				})
			}

			// Jobs
			jbs, _ := k.Jobs(JobOptions{
				Namespace: ns.Namespace,
				Logger:    options.Logger,
				Context:   options.Context,
			})

			for _, jb := range jbs {
				nsOverviews[index] = append(nsOverviews[index], Overview{
					LinkedName: jb.LinkedName,
					Namespace:  jb.Namespace,
				})
			}

			// Pods
			povs, _ := k.Pods(PodOptions{
				Namespace: ns.Namespace,
				Logger:    options.Logger,
				Context:   options.Context,
			})

			for _, pov := range povs {
				nsOverviews[index] = append(nsOverviews[index], Overview{
					LinkedName: pov.LinkedName,
					Namespace:  pov.Namespace,
				})
			}
		}(i, namespace)
	}

	wg.Wait()

	overviews = []Overview{}

	for _, nsovs := range nsOverviews {
		for _, nsov := range nsovs {
			if len(nsov.LinkedName) > 0 {
				found := false
				for _, ov := range overviews {
					if strings.EqualFold(nsov.LinkedName, ov.LinkedName) && strings.EqualFold(nsov.Namespace, ov.Namespace) {
						found = true
					}
				}
				if !found {
					overviews = append(overviews, Overview{
						LinkedName: nsov.LinkedName,
						Namespace:  nsov.Namespace,
					})
				}
			}
		}
	}

	return overviews, nil
}
