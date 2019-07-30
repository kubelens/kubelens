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

package rbac

import (
	"strings"

	"github.com/kubelens/kubelens/api/config"
)

// Role represents RBAC groups
type Role struct {
	Viewers        bool     `json:"viewers"`
	Operators      bool     `json:"operators"`
	MatchPrefix    string   `json:"matchPrefix"`
	MatchSplitChar string   `json:"matchSplitChar"`
	MatchLabels    []string `json:"matchLabels"`
	Exclusions     []string `json:"exclusions"`
}

// Matches returns a boolean given matching rules for labels and appname
func (r Role) Matches(labels map[string]string, appname *string) bool {
	if config.C.EnableRBAC {
		for key, value := range labels {
			// obviously cant access if excluded
			if !r.InExclusions(value) && !r.InExclusions(key) {

				// if filtering by appname, only add pods matching that "appname label" value if
				// not in exclusions list. break out and move to next value if found.
				if appname != nil && len(*appname) > 0 {
					if !r.InExclusions(*appname) && strings.Contains(*appname, value) {
						return true
					}
				} else {
					// operators automatically get any pod not in exclusions. break if found
					if r.Operators {
						return true
					}

					// viewers automatically get any app label not in exclusions
					if r.Viewers && strings.EqualFold(key, "app") {
						return true
					}
				}
			}
		}
		return false
	}
	return true
}

// CompareLabels - if MatchLabels has values, we will assume RBAC includes
// scoping results to labels matched with the user allowed apps
func (r Role) CompareLabels(labels map[string]string, exact bool) bool {
	if config.C.EnableRBAC {
		if r.Operators {
			if labels != nil {
				for _, label := range labels {
					if r.InExclusions(label) {
						return false
					}
				}
			}
			return true
		}

		// short ciruit for null labels for those not in operators role
		if labels == nil {
			return false
		}

		if len(r.MatchLabels) > 0 {
			canAccess := false
			// check labels against user allowed list
			for _, mlbl := range r.MatchLabels {
				// without a splitter each char will be an array element.
				if len(r.MatchSplitChar) == 0 {
					return false
				}
				// ex. app=appgroup
				lblParts := strings.Split(mlbl, r.MatchSplitChar)
				// find a matching label and return found to grant access
				if value, ok := labels[lblParts[0]]; ok && !r.InExclusions(value) {
					// if exact, must match the label value as well (application name typically)
					if exact {
						if strings.EqualFold(value, lblParts[1]) {
							canAccess = true
						}
					} else {
						canAccess = true
					}
				}
			}
			// default
			return canAccess
		}
		// skip if not scoping by labels
		return true
	}
	return true
}

// InExclusions checks to see if exclusions list CONTAINS value
func (r Role) InExclusions(value string) bool {
	if config.C.EnableRBAC {
		if len(r.Exclusions) > 0 {
			for _, exclude := range r.Exclusions {
				// Using contains allows the blacklist to contain more than just exact label matches.
				// This allows exclusions for namespaces that might have an environment prefix, e.g. dev-mongo
				// which the blacklist can have a value of "mongo" to exclude regardless of environment.
				if ok := strings.Contains(value, exclude); ok {
					return true
				}
			}
		}

		return false
	}
	return true
}

// GetMatchLabels returns role.MatchLabels
func (r Role) GetMatchLabels() []string {
	return r.MatchLabels
}
