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
	"strings"

	"github.com/kubelens/kubelens/api/config"
)

func getLinkedName(labels map[string]string) string {
	return labels[config.C.LabelKeyLink]
}

func generateLabelSelector(value string) string {
	return fmt.Sprintf("%s=%s", config.C.LabelKeyLink, value)
}

// Hard coded check for now, but there really shouldn't be any
// secrets exposed unless the application has them stored in an
// insecure way such as a configMap or through environment variables.
func stringContainsSensitiveInfo(toCheck string) bool {
	l := strings.ToLower(toCheck)
	if strings.Contains(l, "pass") ||
		strings.Contains(l, "key") ||
		strings.Contains(l, "secret") {
		return true
	}
	return false
}
