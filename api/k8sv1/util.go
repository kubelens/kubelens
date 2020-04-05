/*
MIT License

Copyright (c) 2019 The KubeLens Authors

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

package k8v1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kubelens/kubelens/api/config"
)

func getProjectSlug(pod string) string {
	if len(config.C.ProjectSlugRegex) > 0 {
		// Example pod name: some-pod-name-slug-1234-6fd8676889-xsksc
		// the "slug-1234" piece is the slug mapping to the the id of the deploy template.
		// Using this could allow building of a link to the deployment.
		regex := regexp.MustCompile(config.C.ProjectSlugRegex)
		return regex.FindString(pod)
	}
	return ""
}

func getAppName(labels map[string]string, appnameLabelKey, defaultLabelKey, defaultName string) (name string, labelKey string) {
	if labels != nil {
		// first check for the app label
		if len(appnameLabelKey) > 0 {
			if len(labels[appnameLabelKey]) > 0 {
				return labels[appnameLabelKey], appnameLabelKey
			}
		}
		// if no app label
		for _, search := range config.C.DefaultSearchLabels {
			if len(labels[search]) > 0 {
				return labels[search], search
			}
		}
	}

	return defaultName, defaultLabelKey
}

// should this be a new config value, or is there a better
// way to handle this?
func getDefaultSearchLabel(selector map[string]string) string {
	// first try to get the app name by spec.selector, which should never be empty, but check just in case
	if selector != nil {
		for _, dlbl := range config.C.DefaultSearchLabels {
			for name, slbl := range selector {
				if strings.Contains(slbl, dlbl) {
					return name
				}
			}
		}
	}
	// try to get the first item in default search labels
	if len(config.C.DefaultSearchLabels) > 0 {
		return config.C.DefaultSearchLabels[0]
	}
	// all else fails, just use "app" which seems to be very common
	return "app"
}

// getDeployerLink tries to build a link to the tool that deploys the application. See doc.go for more info.
func getDeployerLink(value string) string {
	slug := getProjectSlug(value)
	if len(slug) > 0 {
		return fmt.Sprintf("%s%s", config.C.DeployerLink, slug)
	}

	if len(config.C.DeployerLink) > 0 {
		return config.C.DeployerLink
	}

	return ""
}
