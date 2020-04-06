/*
MIT License

Copyright (c) 2019 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package k8v1

import (
	"io"

	"github.com/kubelens/kubelens/api/errs"
)

const defaultErrorMessage = "Error retrieving container info, please contact your admin."

// Clienter is the interface for Client
type Clienter interface {
	// SanityCheck tries to list pods. if it can't, return will be false, else true.
	// really only used for a sanity/health check.
	SanityCheck() (success bool)
	// AppOverview returns an list of application overviews with high level info such as pods, services, deployments, etc.
	AppOverview(options AppOverviewOptions) (ao *AppOverview, apiErr *errs.APIError)
	// PodDetail returns details for a pod
	PodDetail(options PodDetailOptions) (po *PodDetail, apiErr *errs.APIError)
	// PodOverview returns a list of application overview which is a high level detail of an application.
	PodOverview(options PodOverviewOptions) (po *PodOverview, apiErr *errs.APIError)
	// Logs returns a list of all logs for pods
	Logs(options LogOptions) (logs Log, apiErr *errs.APIError)
	// ReadLogs returns an io.ReadCloser to live stream logs for a pod
	ReadLogs(options LogOptions) (rc io.ReadCloser, apiErr *errs.APIError)
	// ServiceOverviews returns a list of services given filter options
	ServiceOverviews(options ServiceOptions) (svco []ServiceOverview, apiErr *errs.APIError)
	// DeploymentOverviews returns a list of deployments given filter options
	DeploymentOverviews(options DeploymentOptions) (deployments []DeploymentOverview, apiErr *errs.APIError)
}

// Client is the wrapper for kubernetes go client commands
type Client struct {
	wrapper Wrapper
}

// New returns a new instance of Client
func New(w Wrapper) Clienter {
	return &Client{w}
}
