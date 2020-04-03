package k8v1

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
