package k8sv1

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Log holds fields related to log output
type Log struct {
	// the name of the pod
	Pod string `json:"pod"`
	// the log output
	Output string `json:"output"`
}

// LogOptions contains fields used for filtering when retrieving application logs
type LogOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// The label key used for the application name, ex. app=some-app-name
	PodName string `json:"podname"`
	// The name of the container to get logs from.
	ContainerName string `json:"containerName"`
	// follow enables streaming
	Follow bool `json:"follow"`
	// tail logs from line. If a stream request, this is ignored.
	Tail int64 `json:"tail"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// Valid validates LogOptions fields
func (a *LogOptions) Valid() *errs.APIError {
	if !a.UserRole.HasNamespaceAccess(a.Namespace) {
		return errs.Unauthorized()
	}

	if !a.ValidPodName() {
		return errs.ValidationError("podname must be provided when getting logs")
	}

	if !a.ValidNamespace() {
		return errs.ValidationError("namespace must be provided when getting logs")
	}

	return nil
}

// ValidNamespace retuns true if the length LogOptions.Namespace is > 0
func (a *LogOptions) ValidNamespace() bool {
	if len(a.Namespace) > 0 {
		return true
	}
	return false
}

// ValidPodName retuns true if the length LogOptions.PodName is > 0
func (a *LogOptions) ValidPodName() bool {
	if len(a.PodName) > 0 {
		return true
	}
	return false
}

// GetTailLines returns PodOverviewOptions.Tail > 0 || default (32, why not)
func (a *LogOptions) GetTailLines() int64 {
	if a.Tail > 0 {
		return a.Tail
	}
	return 100
}

// UserCanAccess validates the user has access
func (a *LogOptions) UserCanAccess(labels map[string]string) bool {
	if !a.UserRole.HasPodAccess(labels) || !a.UserRole.HasLogAccess(labels) || !a.UserRole.Matches(labels, nil) {
		return false
	}
	return true
}

// Logs returns a list of all logs for pods
func (k *Client) Logs(options LogOptions) (logs Log, apiErr *errs.APIError) {
	// get logs
	rc, apiErr := k.ReadLogs(options)

	// add a message to the output for user if errror
	if apiErr != nil {
		logs.Output = fmt.Sprintf("Could not read request stream when retrieving logs for %s/%s", options.Namespace, options.PodName)
		return logs, apiErr
	}

	defer rc.Close()

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, rc)

	if err != nil {
		return logs, errs.InternalServerError(err.Error())
	}

	o := buf.String()

	if len(o) > 0 {
		logs.Output = buf.String()
	}

	return logs, nil
}

// ReadLogs returns an io.ReadCloser to live stream logs for a pod. Error codes will be the same
// as standard http error codes, but using the values directly so this package doesn't need to import http.
func (k *Client) ReadLogs(options LogOptions) (rc io.ReadCloser, apiErr *errs.APIError) {
	if apiErr = options.Valid(); apiErr != nil {
		return nil, apiErr
	}

	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list := clientset.
		CoreV1().
		Pods(options.Namespace)

	// get pod for labels to check against role
	if pd, err := list.Get(options.PodName, metav1.GetOptions{}); pd != nil {
		if err != nil {
			klog.Trace()
			return nil, errs.InternalServerError(err.Error())
		}

		labels := pd.GetLabels()
		if !options.UserCanAccess(labels) {
			klog.Trace()
			return nil, errs.Forbidden()
		}
	}

	tail := options.Tail
	// ensure tail set to 1 line since we would be streaming
	if options.Follow {
		tail = 1
	}
	// start stream at last line
	req := list.GetLogs(options.PodName, &v1.PodLogOptions{
		Container: options.ContainerName,
		TailLines: &tail,
		Follow:    options.Follow,
	})

	stream, err := req.Stream()

	if err != nil {
		stream.Close()
		return nil, errs.InternalServerError(err.Error())
	}

	return stream, nil
}
