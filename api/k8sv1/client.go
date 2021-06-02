/*
MIT License

Copyright (c) 2020 The KubeLens Authors

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

package k8sv1

import (
	"io"

	"github.com/kubelens/kubelens/api/errs"
)

const defaultErrorMessage = "Error retrieving container info, please contact your admin."

// Clienter is the interface for Client
type Clienter interface {
	// SanityCheck tries to list pods. if it can't, return will be error, else nil.
	// really only used for a sanity/health check.
	SanityCheck() (apiErr *errs.APIError)
	// Overview returns a list of things running in kubernetes, determined by searching deployments for each namespace.
	// For each namespace, Kubernetes Kinds are searched for the type of application, e.g. Service, DaemonSet, etc.
	Overview(options OverviewOptions) (overviews *Overview, apiErr *errs.APIError)
	// Overviews returns an list of application overviews with high level info such as name & namespace
	Overviews(options OverviewOptions) (overviews []Overview, apiErr *errs.APIError)
	// Pod returns the pod found by name or labels.
	Pod(options PodOptions) (overview *PodOverview, apiErr *errs.APIError)
	// Pods returns a list of pods given filter options
	Pods(options PodOptions) (overviews []PodOverview, apiErr *errs.APIError)
	// Service returns the service found by name or labels.
	Service(options ServiceOptions) (overview *ServiceOverview, apiErr *errs.APIError)
	// Services returns a list of services given filter options
	Services(options ServiceOptions) (overviews []ServiceOverview, apiErr *errs.APIError)
	// Deployment returns the deployment found by name or labels.
	Deployment(options DeploymentOptions) (overview *DeploymentOverview, apiErr *errs.APIError)
	// Deployments returns a list of deployments given filter options
	Deployments(options DeploymentOptions) (overviews []DeploymentOverview, apiErr *errs.APIError)
	// DaemonSet returns the daemonset found by name or labels.
	DaemonSet(options DaemonSetOptions) (overview *DaemonSetOverview, apiErr *errs.APIError)
	// DaemonSets returns a list of daemonsets given filter options
	DaemonSets(options DaemonSetOptions) (overviews []DaemonSetOverview, apiErr *errs.APIError)
	// Job returns the job found by name or labels.
	Job(options JobOptions) (overview *JobOverview, apiErr *errs.APIError)
	// Jobs returns a list of jobs given filter options
	Jobs(options JobOptions) (overviews []JobOverview, apiErr *errs.APIError)
	// ReplicaSet returns the replicaset found by name or labels.
	ReplicaSet(options ReplicaSetOptions) (overview *ReplicaSetOverview, apiErr *errs.APIError)
	// ReplicaSets returns a list of replicasets given filter options
	ReplicaSets(options ReplicaSetOptions) (overviews []ReplicaSetOverview, apiErr *errs.APIError)
	// Logs returns a list of all logs for pods
	Logs(options LogOptions) (logs Log, apiErr *errs.APIError)
	// ReadLogs returns an io.ReadCloser to live stream logs for a pod
	ReadLogs(options LogOptions) (rc io.ReadCloser, apiErr *errs.APIError)
}

// Client is the wrapper for kubernetes go client commands
type Client struct {
	wrapper Wrapper
}

// New returns a new instance of Client
func New(w Wrapper) Clienter {
	return &Client{w}
}
