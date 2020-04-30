package k8sv1

import (
	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapOptions .
type ConfigMapOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// ConfigMaps returns a list of config maps. If LabelSelector is provided, only config map labels matching that selector
// will be returned
func (k *Client) ConfigMaps(options ConfigMapOptions) (configMaps []v1.ConfigMap, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	lo := metav1.ListOptions{
		IncludeUninitialized: true,
	}

	// get this list of config maps here to ensure correct namespace. This should perform
	// just fine, but if not, this list could be generated outsid of this routine.
	cmList, err := clientset.CoreV1().ConfigMaps(options.Namespace).List(lo)
	// just trace the error and move on, shouldn't be critical.
	if err != nil {
		klog.Trace()
	}

	if cmList != nil {
		jConfigMaps := []v1.ConfigMap{}

		for _, cm := range cmList.Items {
			cmLabels := cm.GetLabels()
			if options.UserRole.HasConfigMapAccess(cmLabels) && labelsContainSelector(options.LabelSelector, cmLabels) {
				jConfigMaps = append(jConfigMaps, cm)
			}
		}

		return jConfigMaps, nil
	}

	return nil, nil
}
