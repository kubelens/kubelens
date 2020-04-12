package fakes

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/kubelens/kubelens/api/errs"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
)

// K8V1 .
type K8V1 struct {
	fail *bool
}

// SanityCheck .
func (m *K8V1) SanityCheck() (apiErr *errs.APIError) {
	if m.fail != nil && *m.fail {
		return errs.InternalServerError("Sanity check error")
	}
	return nil
}

// Apps .
func (m *K8V1) Apps(options k8sv1.AppOptions) (apps []k8sv1.App, apiErr *errs.APIError) {
	if m.fail != nil && *m.fail {
		return nil, errs.InternalServerError("Apps Test Error")
	}

	return []k8sv1.App{
		k8sv1.App{
			Name:      "test-service",
			Namespace: "default",
			Kind:      "Service",
		},
		k8sv1.App{
			Name:      "test-daemonset",
			Namespace: "default",
			Kind:      "DaemonSet",
		},
	}, nil
}

// AppOverview .
func (m *K8V1) AppOverview(options k8sv1.AppOverviewOptions) (ao *k8sv1.AppOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return ao, errs.InternalServerError("AppOverview Test Error")
	}

	return &k8sv1.AppOverview{
		PodOverviews: k8sv1.PodOverview{
			Name:      "testpod",
			Namespace: "default",
			PodInfo: []*k8sv1.PodInfo{
				&k8sv1.PodInfo{
					Name: "testpod",
				},
			},
		},
		ServiceOverviews: []k8sv1.ServiceOverview{
			k8sv1.ServiceOverview{
				Name: "service-name",
			},
		},
	}, nil
}

// PodDetail .
func (m *K8V1) PodDetail(options k8sv1.PodDetailOptions) (po *k8sv1.PodDetail, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return po, errs.InternalServerError("PodDetail Test Error")
	}

	return &k8sv1.PodDetail{
		Name: "testpod",
	}, nil
}

// PodOverview .
func (m *K8V1) PodOverview(options k8sv1.PodOverviewOptions) (po *k8sv1.PodOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return po, errs.InternalServerError("GetApps Test Error")
	}

	return &k8sv1.PodOverview{
		Name:      "test",
		Namespace: "default",
		PodInfo: []*k8sv1.PodInfo{
			&k8sv1.PodInfo{
				Name: "testpod",
			},
		},
	}, nil
}

// Logs .
func (m *K8V1) Logs(options k8sv1.LogOptions) (logs k8sv1.Log, apiErr *errs.APIError) {
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
func (m *K8V1) ReadLogs(options k8sv1.LogOptions) (rc io.ReadCloser, apiErr *errs.APIError) {
	stringReader := strings.NewReader("message\n")
	stringReadCloser := ioutil.NopCloser(stringReader)
	return stringReadCloser, nil
}

// ServiceOverviews .
func (m *K8V1) ServiceOverviews(options k8sv1.ServiceOptions) (svco []k8sv1.ServiceOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("ServiceOverviews Test Error")
	}

	return []k8sv1.ServiceOverview{
		k8sv1.ServiceOverview{
			Name: "service-name",
		},
	}, nil
}

// ServiceAppInfos .
func (m *K8V1) ServiceAppInfos(options k8sv1.ServiceOptions) (info []k8sv1.AppInfo, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("ServiceAppInfos Test Error")
	}

	return []k8sv1.AppInfo{
		k8sv1.AppInfo{
			Name: "service-name",
		},
	}, nil
}

// DeploymentOverviews .
func (m *K8V1) DeploymentOverviews(options k8sv1.DeploymentOptions) (deployments []k8sv1.DeploymentOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("DeploymentOverviews Test Error")
	}

	return []k8sv1.DeploymentOverview{
		k8sv1.DeploymentOverview{
			Name: "service-name",
		},
	}, nil
}

// DaemonSetOverviews .
func (m *K8V1) DaemonSetOverviews(options k8sv1.DaemonSetOptions) (daemonsets []k8sv1.DaemonSetOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("DaemonSetOverviews Test Error")
	}

	return []k8sv1.DaemonSetOverview{
		k8sv1.DaemonSetOverview{
			Name: "daemonset-name",
		},
	}, nil
}

// DaemonSetAppInfos .
func (m *K8V1) DaemonSetAppInfos(options k8sv1.DaemonSetOptions) (info []k8sv1.AppInfo, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("DaemonSetAppInfos Test Error")
	}

	return []k8sv1.AppInfo{
		k8sv1.AppInfo{
			Name: "daemonset-name",
		},
	}, nil
}

// JobOverviews .
func (m *K8V1) JobOverviews(options k8sv1.JobOptions) (jobs []k8sv1.JobOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("JobOverviews Test Error")
	}

	return []k8sv1.JobOverview{
		k8sv1.JobOverview{
			Name: "job-name",
		},
	}, nil
}

// JobAppInfos .
func (m *K8V1) JobAppInfos(options k8sv1.JobOptions) (info []k8sv1.AppInfo, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("JobAppInfos Test Error")
	}

	return []k8sv1.AppInfo{
		k8sv1.AppInfo{
			Name: "job-name",
		},
	}, nil
}
