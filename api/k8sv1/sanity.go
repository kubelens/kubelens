package k8sv1

import (
	"context"

	"github.com/kubelens/kubelens/api/errs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SanityCheck tries to list pods. if it can't, return will be false, else true.
// really only used for a sanity/health check.
func (k *Client) SanityCheck() (apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		return errs.InternalServerError(err.Error())
	}

	var tm int64 = 5

	list, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{
		Limit:          1,
		TimeoutSeconds: &tm,
	})

	if err != nil {
		return errs.InternalServerError(err.Error())
	}

	if len(list.Items) == 0 {
		return errs.ValidationError("no service kinds found")
	}

	return nil
}
