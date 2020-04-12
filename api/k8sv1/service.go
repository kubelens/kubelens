package k8sv1

import (
	"sync"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"

	klog "github.com/kubelens/kubelens/api/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceOverview provides field relating to a service. Can return full values or slimmed down values.
type ServiceOverview struct {
	FriendlyName        string               `json:"friendlyName"`
	DeployerLink        string               `json:"deployerLink,omitempty"`
	Name                string               `json:"name"`
	Namespace           string               `json:"namespace"`
	Type                v1.ServiceType       `json:"type"`
	Selector            map[string]string    `json:"selector,omitempty"`
	ClusterIP           string               `json:"clusterIP,omitempty"`
	LoadBalancerIP      string               `json:"loadBalancerIP,omitempty"`
	ExternalIPs         []string             `json:"externalIPs,omitempty"`
	Ports               []v1.ServicePort     `json:"ports,omitempty"`
	Spec                *v1.ServiceSpec      `json:"spec,omitempty"`
	Status              *v1.ServiceStatus    `json:"status,omitempty"`
	ConfigMaps          *[]v1.ConfigMap      `json:"configMaps,omitempty"`
	DeploymentOverviews []DeploymentOverview `json:"deploymentOverviews,omitempty"`
}

// AddDetail sets additional fields for more detailed overview
func (s *ServiceOverview) AddDetail(spec *v1.ServiceSpec, status *v1.ServiceStatus) {
	s.Type = spec.Type
	s.Ports = spec.Ports
	s.ClusterIP = spec.ClusterIP
	s.LoadBalancerIP = spec.LoadBalancerIP
	s.Selector = spec.Selector
	s.Spec = spec
	s.Status = status
}

// AddConfigMaps sets the ConfigMaps value. Normally wouldn't have a pointer for a slice,
// but this allows to easily return empty for the client.
func (s *ServiceOverview) AddConfigMaps(cms *[]v1.ConfigMap) {
	s.ConfigMaps = cms
}

// AddDeploymentOverviews sets the DeploymentOverviews value. Adding separately for ease of
// checking access before hand.
func (s *ServiceOverview) AddDeploymentOverviews(d []DeploymentOverview) {
	s.DeploymentOverviews = d
}

// ServiceOptions contains fields used for filtering when retrieving services
type ServiceOptions struct {
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the label selector to match kinds
	LabelSelector map[string]string `json:"labelSelector"`
	// include full detail
	Detailed bool `json:"detailed"`
	//users role assignemnt
	UserRole rbac.RoleAssignmenter
	// logger instance
	Logger klog.Logger
}

// ValidNamespace validates the namespace field
func (a *ServiceOptions) ValidNamespace() bool {
	if len(a.Namespace) > 0 {
		return true
	}
	return false
}

// ValidLabelSearch validates the label search field
func (a *ServiceOptions) ValidLabelSearch() bool {
	if len(a.LabelSelector) > 0 {
		return true
	}
	return false
}

// UserCanAccess validates the user has access
func (a *ServiceOptions) UserCanAccess(labels map[string]string) bool {
	if !a.UserRole.HasApplicationAccess() || !a.UserRole.HasServiceAccess(labels) {
		return false
	}
	return true
}

// ServiceOverviews returns a list of services given filter options
func (k *Client) ServiceOverviews(options ServiceOptions) (svco []ServiceOverview, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	lo := metav1.ListOptions{
		IncludeUninitialized: true,
	}

	if len(options.LabelSelector) > 0 {
		lo.LabelSelector = toLabelSelectorString(options.LabelSelector)
	}

	list, err := clientset.CoreV1().Services(options.Namespace).List(lo)

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	svco = make([]ServiceOverview, len(list.Items))
	wg := sync.WaitGroup{}

	wg.Add(len(list.Items))

	for i, item := range list.Items {
		// check access at label level
		if options.UserCanAccess(item.GetLabels()) {
			go func(index int, service v1.Service) {
				defer wg.Done()

				svc := ServiceOverview{
					FriendlyName: getFriendlyAppName(
						item.GetLabels(),
						item.GetName(),
					),
					Selector:     service.Spec.Selector,
					DeployerLink: getDeployerLink(service.GetName()),
					Name:         service.GetName(),
					Namespace:    service.GetNamespace(),
				}

				if options.Detailed {
					svc.AddDetail(&service.Spec, &service.Status)
				}

				// add deployments per service for ease of display by client since deployments are really
				// specific to certian K8s kinds.
				if options.Detailed && options.UserRole.HasDeploymentAccess(item.GetLabels()) {
					deployments, err := k.DeploymentOverviews(DeploymentOptions{
						LabelSelector: svc.Selector,
						Namespace:     svc.Namespace,
						UserRole:      options.UserRole,
						Logger:        options.Logger,
					})
					// just trace the error and move on, shouldn't be critical.
					if err != nil {
						klog.Trace()
					}
					svc.AddDeploymentOverviews(deployments)
				}

				if options.UserRole.HasConfigMapAccess(item.GetLabels()) {
					cms, err := k.ConfigMaps(ConfigMapOptions{
						Namespace:     svc.Namespace,
						LabelSelector: svc.Selector,
					})

					// just trace the error and move on, shouldn't be critical.
					if err != nil {
						klog.Trace()
					}

					if len(cms) > 0 {
						svc.AddConfigMaps(&cms)
					}
				}

				svco[index] = svc
			}(i, item)
		} else {
			wg.Done()
		}
	}

	wg.Wait()

	return svco, nil
}

// ServiceAppInfos returns basic info for all services found for a given namespace.
func (k *Client) ServiceAppInfos(options ServiceOptions) (info []AppInfo, apiErr *errs.APIError) {
	clientset, err := k.wrapper.GetClientSet()

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	list, err := clientset.CoreV1().Services(options.Namespace).List(metav1.ListOptions{IncludeUninitialized: true})

	if err != nil {
		klog.Trace()
		return nil, errs.InternalServerError(err.Error())
	}

	errors := []*errs.APIError{}

	if list != nil {
		info = make([]AppInfo, len(list.Items))

		wg := sync.WaitGroup{}
		wg.Add(len(list.Items))

		for i, item := range list.Items {
			go func(index int, svc v1.Service) {
				defer wg.Done()

				name := getFriendlyAppName(
					svc.GetLabels(),
					svc.GetName(),
				)

				info[index] = AppInfo{
					FriendlyName:  name,
					Name:          svc.GetName(),
					Namespace:     svc.GetNamespace(),
					Kind:          "Service",
					LabelSelector: svc.Spec.Selector,
				}
			}(i, item)
		}
		wg.Wait()
	}

	if len(errors) > 0 {
		if len(errors) > 0 {
			return info, errs.ListToInternalServerError(errors)
		}
	}

	return info, nil
}
