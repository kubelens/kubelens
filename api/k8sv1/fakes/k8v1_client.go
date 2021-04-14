package fakes

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/kubelens/kubelens/api/errs"
	"github.com/kubelens/kubelens/api/k8sv1"
)

// K8sV1 .
type K8sV1 struct {
	fail *bool
}

// SanityCheck .
func (m *K8sV1) SanityCheck() (apiErr *errs.APIError) {
	if m.fail != nil && *m.fail {
		return errs.InternalServerError("Sanity check error")
	}
	return nil
}

// Overview .
func (m *K8sV1) Overview(options k8sv1.OverviewOptions) (overview *k8sv1.Overview, apiErr *errs.APIError) {
	if m.fail != nil && *m.fail {
		return nil, errs.InternalServerError("Overview Test Error")
	}

	return &k8sv1.Overview{
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
		DaemonSets: []k8sv1.DaemonSetOverview{
			{
				Name:       options.LinkedName + "-daemonset",
				LinkedName: options.LinkedName,
				Namespace:  options.Namespace,
			},
		},
		Deployments: []k8sv1.DeploymentOverview{
			{
				Name:       options.LinkedName + "-deployment",
				LinkedName: options.LinkedName,
				Namespace:  options.Namespace,
			},
		},
		Jobs: []k8sv1.JobOverview{
			{
				Name:       options.LinkedName + "-job",
				LinkedName: options.LinkedName,
				Namespace:  options.Namespace,
			},
		},
		Pods: []k8sv1.PodOverview{
			{
				Name:       options.LinkedName + "-pod",
				LinkedName: options.LinkedName,
				Namespace:  options.Namespace,
			},
		},
		ReplicaSets: []k8sv1.ReplicaSetOverview{
			{
				Name:       options.LinkedName + "-replicaset",
				LinkedName: options.LinkedName,
				Namespace:  options.Namespace,
			},
		},
		Services: []k8sv1.ServiceOverview{
			{
				Name:       options.LinkedName + "-service",
				LinkedName: options.LinkedName,
				Namespace:  options.Namespace,
			},
		},
	}, nil
}

// Overviews .
func (m *K8sV1) Overviews(options k8sv1.OverviewOptions) (overviews []k8sv1.Overview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("Overviews Test Error")
	}

	return []k8sv1.Overview{
		{
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}

// Pod .
func (m *K8sV1) Pod(options k8sv1.PodOptions) (overview *k8sv1.PodOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overview, errs.InternalServerError("Pod Test Error")
	}

	return &k8sv1.PodOverview{
		Name:       options.Name + "-pod",
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
	}, nil
}

// Pods .
func (m *K8sV1) Pods(options k8sv1.PodOptions) (overviews []k8sv1.PodOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("Pods Test Error")
	}

	return []k8sv1.PodOverview{
		{
			Name:       options.Name + "-pod",
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}

// Logs .
func (m *K8sV1) Logs(options k8sv1.LogOptions) (logs k8sv1.Log, apiErr *errs.APIError) {
	if options.Namespace == "bad2" {
		logs = k8sv1.Log{
			Pod:    options.PodName,
			Output: "No logs returned from K8",
		}
		return logs, errs.InternalServerError("Logs Test Error")
	}

	return k8sv1.Log{
		Pod:    options.PodName,
		Output: "some output",
	}, nil
}

// ReadLogs .
func (m *K8sV1) ReadLogs(options k8sv1.LogOptions) (rc io.ReadCloser, apiErr *errs.APIError) {
	stringReader := strings.NewReader("message\n")
	stringReadCloser := ioutil.NopCloser(stringReader)
	return stringReadCloser, nil
}

// Service .
func (m *K8sV1) Service(options k8sv1.ServiceOptions) (overview *k8sv1.ServiceOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overview, errs.InternalServerError("Service Test Error")
	}

	return &k8sv1.ServiceOverview{
		Name:       options.Name + "-service",
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
	}, nil
}

// Services .
func (m *K8sV1) Services(options k8sv1.ServiceOptions) (overviews []k8sv1.ServiceOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("Services Test Error")
	}

	return []k8sv1.ServiceOverview{
		{
			Name:       options.Name + "-service",
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}

// Deployment .
func (m *K8sV1) Deployment(options k8sv1.DeploymentOptions) (overview *k8sv1.DeploymentOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overview, errs.InternalServerError("Deployment Test Error")
	}

	return &k8sv1.DeploymentOverview{
		Name:       options.Name + "-deployment",
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
	}, nil
}

// Deployments .
func (m *K8sV1) Deployments(options k8sv1.DeploymentOptions) (overviews []k8sv1.DeploymentOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("Deployments Test Error")
	}

	return []k8sv1.DeploymentOverview{
		{
			Name:       options.Name + "-deployment",
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}

// DaemonSet .
func (m *K8sV1) DaemonSet(options k8sv1.DaemonSetOptions) (overview *k8sv1.DaemonSetOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overview, errs.InternalServerError("DaemonSet Test Error")
	}

	return &k8sv1.DaemonSetOverview{
		Name:       options.Name + "-daemonset",
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
	}, nil
}

// DaemonSets .
func (m *K8sV1) DaemonSets(options k8sv1.DaemonSetOptions) (overviews []k8sv1.DaemonSetOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("DaemonSets Test Error")
	}

	return []k8sv1.DaemonSetOverview{
		{
			Name:       options.Name + "-daemonset",
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}

// Job .
func (m *K8sV1) Job(options k8sv1.JobOptions) (overview *k8sv1.JobOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overview, errs.InternalServerError("Job Test Error")
	}

	return &k8sv1.JobOverview{
		Name:       options.Name + "-job",
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
	}, nil
}

// Jobs .
func (m *K8sV1) Jobs(options k8sv1.JobOptions) (overviews []k8sv1.JobOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("DaemonSets Test Error")
	}

	return []k8sv1.JobOverview{
		{
			Name:       options.Name + "-job",
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}

// ReplicaSet .
func (m *K8sV1) ReplicaSet(options k8sv1.ReplicaSetOptions) (overview *k8sv1.ReplicaSetOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overview, errs.InternalServerError("ReplicaSet Test Error")
	}

	return &k8sv1.ReplicaSetOverview{
		Name:       options.Name + "-replicaset",
		LinkedName: options.LinkedName,
		Namespace:  options.Namespace,
	}, nil
}

// ReplicaSets .
func (m *K8sV1) ReplicaSets(options k8sv1.ReplicaSetOptions) (overviews []k8sv1.ReplicaSetOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return overviews, errs.InternalServerError("ReplicaSet Test Error")
	}

	return []k8sv1.ReplicaSetOverview{
		{
			Name:       options.Name + "-replicaset",
			LinkedName: options.LinkedName,
			Namespace:  options.Namespace,
		},
	}, nil
}
