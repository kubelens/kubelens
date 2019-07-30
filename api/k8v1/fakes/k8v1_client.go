package fakes

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/kubelens/kubelens/api/errs"
	"github.com/kubelens/kubelens/api/k8v1"
)

type K8V1 struct{}

func (m *K8V1) SanityCheck() (success bool) {
	return true
}

func (m *K8V1) AppOverview(options k8v1.AppOverviewOptions) (ao *k8v1.AppOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return ao, errs.InternalServerError("AppOverview Test Error")
	}

	return &k8v1.AppOverview{
		PodOverviews: k8v1.PodOverview{
			Name: k8v1.Name{
				LabelKey: "app",
				Value:    "testpod",
			},
			Namespace: "default",
			PodDetails: []*k8v1.PodDetail{
				&k8v1.PodDetail{
					Name: "testpod",
				},
			},
		},
		ServiceOverviews: []k8v1.ServiceOverview{
			k8v1.ServiceOverview{
				Name: "service-name",
			},
		},
	}, nil
}

func (m *K8V1) PodDetail(options k8v1.PodDetailOptions) (po *k8v1.PodDetail, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return po, errs.InternalServerError("PodDetail Test Error")
	}

	return &k8v1.PodDetail{
		Name: "testpod",
	}, nil
}

func (m *K8V1) PodOverview(options k8v1.PodOverviewOptions) (po *k8v1.PodOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return po, errs.InternalServerError("GetApps Test Error")
	}

	return &k8v1.PodOverview{
		Name: k8v1.Name{
			LabelKey: options.AppNameLabelKey,
			Value:    "test",
		},
		Namespace: "default",
		PodDetails: []*k8v1.PodDetail{
			&k8v1.PodDetail{
				Name: "testpod",
			},
		},
	}, nil
}

func (m *K8V1) Logs(options k8v1.LogOptions) (logs k8v1.Log, apiErr *errs.APIError) {
	if options.Namespace == "bad2" {
		logs = k8v1.Log{
			Pod:    options.PodName,
			Output: "No logs returned from K8",
		}
		return logs, errs.InternalServerError("Logs Test Error")
	}

	return k8v1.Log{
		Pod:    options.PodName,
		Output: "some output",
	}, nil
}

func (m *K8V1) ReadLogs(options k8v1.LogOptions) (rc io.ReadCloser, apiErr *errs.APIError) {
	stringReader := strings.NewReader("message\n")
	stringReadCloser := ioutil.NopCloser(stringReader)
	return stringReadCloser, nil
}

func (m *K8V1) ServiceOverviews(options k8v1.ServiceOptions) (svco []k8v1.ServiceOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("ServiceOverviews Test Error")
	}

	return []k8v1.ServiceOverview{
		k8v1.ServiceOverview{
			Name: "service-name",
		},
	}, nil
}
