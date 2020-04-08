package v1

import (
	"strings"
)

// Req .
//
// swagger:model
type Req struct {
	// the name of the application
	Appname string `json:"appname,omitempty"`
	// the name of a pod
	PodName string `json:"podName,omitempty"`
	// the name of a container in the pod
	ContainerName string `json:"containerName,omitempty"`
	// the namespace to filter
	Namespace string `json:"namespace,omitempty"`
	// the label selector to search, example: ?labelSelector="app=some-app,app.kubernetes.io/name=app"
	LabelSelector string `json:"labelSelector"`
	// include more detail
	Detailed bool `json:"detailed"`
	// the number of lines to grab from the output
	Tail int `json:"tail,omitempty"`
}

// LabelSelectorMap returns the field LabelSelector as a map[string]string
func (r Req) LabelSelectorMap() (labelSelector map[string]string) {
	labelSelector = map[string]string{}

	if len(r.LabelSelector) > 0 {
		selectors := strings.Split(r.LabelSelector, ",")
		for _, pair := range selectors {
			labels := strings.Split(pair, "=")
			labelSelector[labels[0]] = labels[1]
		}
	}

	return labelSelector
}

// App .
type App struct {
	// name of the application
	Name string `json:"name"`
	// the namespace of the app
	Namespace string `json:"namespace"`
	// kind of application, e.g. Service, DaemonSet
	Kind string `json:"kind"`
	// the label selector to match kinds
	LabelSelector string `json:"labelSelector"`
	// deployer link if any
	DeployerLink string `json:"deployerLink,omitempty"`
}
