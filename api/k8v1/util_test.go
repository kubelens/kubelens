package k8v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kubelens/kubelens/api/config"
)

func TestSvcNameFromPod(t *testing.T) {
	config.Set("../testdata/mock_config.json")
	a := "app-projects-1234-dev-6fd8676889-xsksc"
	b := "app-projects-1234-prod-6fd8676889-xsksc"
	c := "app-notmatchregex"

	ar := svcNameFromPod(a)
	br := svcNameFromPod(b)
	cr := svcNameFromPod(c)

	fmt.Printf("\n\nAR: %s - %v\n\n", ar, strings.Compare(ar, "app-projects-1234-dev"))
	assert.Equal(t, "app-projects-1234-dev", ar)
	assert.Equal(t, "app-projects-1234-prod", br)
	assert.Equal(t, "", cr)
}

func TestGetProjectSlug(t *testing.T) {
	config.Set("../testdata/mock_config.json")
	a1 := "app-projects-1234-dev-6fd8676889-xsksc"

	assert.Equal(t, "projects-1234", getProjectSlug(a1))
}
