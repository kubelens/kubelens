/*
MIT License

Copyright (c) 2020 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package k8sv1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kubelens/kubelens/api/config"
)

func getProjectSlug(value string) string {
	if len(config.C.ProjectSlugRegex) > 0 {
		// Example value name: some-pod-name-slug-1234-6fd8676889-xsksc
		// the "slug-1234" piece is the slug mapping to the the id of the deploy template.
		// Using this could allow building of a link to the deployment.
		regex := regexp.MustCompile(config.C.ProjectSlugRegex)
		return regex.FindString(value)
	}
	return ""
}

func getFriendlyAppName(labels map[string]string, defaultName string) (name string) {
	if labels != nil {
		// if no app label
		for _, search := range config.C.AppNameLabelKeys {
			if len(labels[search]) > 0 {
				return labels[search]
			}
		}
	}

	return defaultName
}

// getDeployerLink tries to build a link to the tool that deploys the application. See doc.go for more info.
func getDeployerLink(value string) string {
	slug := getProjectSlug(value)
	if len(slug) > 0 {
		return fmt.Sprintf("%s%s", config.C.DeployerLink, slug)
	}

	return ""
}

func toLabelSelectorString(labelSelector map[string]string) (labelSelectorString string) {
	for k, v := range labelSelector {
		labelSelectorString += fmt.Sprintf("%s=%s,", k, v)
	}
	return strings.TrimSuffix(labelSelectorString, ",")
}
