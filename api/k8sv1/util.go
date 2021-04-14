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
	"strings"

	"github.com/kubelens/kubelens/api/config"
)

func labelsContainSelector(selector string, labels map[string]string) bool {
	if labels == nil || len(labels) == 0 || len(selector) == 0 {
		return false
	}

	if strings.EqualFold(labels[config.C.LabelKeyLink], selector) {
		return true
	}

	// default to true if no labels are provided to allow all
	return false
}

func getLinkedName(labels map[string]string) string {
	return labels[config.C.LabelKeyLink]
}
