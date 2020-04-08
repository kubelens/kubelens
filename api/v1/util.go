package v1

import (
	"fmt"
	"strings"

	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
)

func appsResponse(inApps []k8sv1.App) (apps []App) {
	apps = []App{}
	for _, item := range inApps {
		apps = append(apps, App{
			Name:          item.Name,
			Namespace:     item.Namespace,
			Kind:          item.Kind,
			LabelSelector: toLabelSelectorString(item.LabelSelector),
			DeployerLink:  item.DeployerLink,
		})
	}
	return apps
}

func toLabelSelectorString(labelSelector map[string]string) (labelSelectorString string) {
	for k, v := range labelSelector {
		labelSelectorString += fmt.Sprintf("%s=%s,", k, v)
	}
	return strings.TrimSuffix(labelSelectorString, ",")
}
