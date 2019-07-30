package k8v1

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/kubelens/kubelens/api/errs"
)

type mockk8 struct{}

func (m *mockk8) SanityCheck() (success bool) {
	return true
}

func (m *mockk8) PodOverview(options PodOverviewOptions) (po *PodOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return po, errs.InternalServerError("GetApps Test Error")
	}

	return &PodOverview{
		Name: Name{
			LabelKey: options.AppNameLabelKey,
			Value:    "test",
		},
		Namespace: "default",
		PodDetails: []*PodDetail{
			&PodDetail{
				Name: "testpod",
			},
		},
	}, nil
}

func (m *mockk8) Logs(options LogOptions) (logs Log, apiErr *errs.APIError) {
	if options.Namespace == "bad2" {
		logs = Log{
			Pod:    options.PodName,
			Output: "No logs returned from K8",
		}
		return logs, errs.InternalServerError("Logs Test Error")
	}

	return Log{
		Pod:    options.PodName,
		Output: "some output",
	}, nil
}

func (m *mockk8) ReadLogs(options LogOptions) (rc io.ReadCloser, apiErr *errs.APIError) {
	stringReader := strings.NewReader("message\n")
	stringReadCloser := ioutil.NopCloser(stringReader)
	return stringReadCloser, nil
}

func (m *mockk8) ServiceOverviews(options ServiceOptions) (svco []ServiceOverview, apiErr *errs.APIError) {
	if options.Namespace == "bad" {
		return nil, errs.InternalServerError("ServiceOverviews Test Error")
	}

	return []ServiceOverview{
		ServiceOverview{
			Name: "service-name",
		},
	}, nil
}

func setupClient(ns, n string, fail, innerFail bool) Clienter {
	w := &mockWrapper{
		namespace: ns,
		appname:   n,
		fail:      fail,
		innerFail: innerFail,
	}
	return New(w)
}
