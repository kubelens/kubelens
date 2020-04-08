package k8sv1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// SanityCheck tries to list pods. if it can't, return will be false, else true.
// really only used for a sanity/health check.
func (k *Client) SanityCheck() (success bool) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		return false
	}

	var tm int64 = 5

	list, err := clientset.CoreV1().Services("").List(metav1.ListOptions{
		Limit:                1,
		IncludeUninitialized: false,
		TimeoutSeconds:       &tm,
	})

	if err != nil || len(list.Items) == 0 {
		return false
	}

	return true
}
