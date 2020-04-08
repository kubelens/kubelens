package v1

import (
	"testing"

	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	"github.com/stretchr/testify/assert"
)

func TestAppsResponse(t *testing.T) {
	apps := []k8sv1.App{
		k8sv1.App{
			Name:          "test-service",
			Namespace:     "default",
			Kind:          "Service",
			LabelSelector: map[string]string{"label1": "value1", "label2": "value2"},
		},
		k8sv1.App{
			Name:          "test-daemonset",
			Namespace:     "default",
			Kind:          "DaemonSet",
			LabelSelector: map[string]string{"label1": "value1", "label2": "value2"},
		},
	}

	res := appsResponse(apps)

	for _, item := range res {
		assert.Equal(t, "label1=value1,label2=value2", item.LabelSelector)
	}
}
