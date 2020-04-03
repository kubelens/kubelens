package v1

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
	// the app name label key, e.g. "app" in label app=appname
	LabelKey string `json:"labelKey"`
	// include more detail
	Detailed bool `json:"detailed"`
	// the number of lines to grab from the output
	Tail int `json:"tail,omitempty"`
}
