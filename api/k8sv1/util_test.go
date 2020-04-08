package k8sv1

import (
	"testing"

	"github.com/kubelens/kubelens/api/config"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectSlug(t *testing.T) {
	config.Set("../testdata/mock_config.json")
	a1 := "app-projects-1234-dev-6fd8676889-xsksc"

	assert.Equal(t, "projects-1234", getProjectSlug(a1))
}

func TestGetProjectSlugMissing(t *testing.T) {
	config.Set("../testdata/mock_config2.json")
	a1 := "app-projects-1234-dev-6fd8676889-xsksc"

	assert.Len(t, getProjectSlug(a1), 0)
}

func TestGetDeployerLinkOnly(t *testing.T) {
	config.Set("../testdata/mock_config2.json")

	assert.Equal(t, "https://test-deployer.com", getDeployerLink(""))
}

func TestGetDeployerLinkWithSlug(t *testing.T) {
	config.Set("../testdata/mock_config.json")

	assert.Equal(t, "https://test-deployer.com/projects-1234", getDeployerLink("app-projects-1234-dev"))
}

func TestGetFriendlyAppNameDefault(t *testing.T) {
	config.Set("../testdata/mock_config.json")

	l := map[string]string{}
	l["notamatch"] = "random-label"

	// should just get the first match
	n := getFriendlyAppName(l, "default-name")
	assert.Equal(t, "default-name", n)
}

func TestGetFriendlyAppNameMatchedt(t *testing.T) {
	config.Set("../testdata/mock_config.json")

	l := map[string]string{}
	l["name"] = "real-app-name"

	// should just get the first match
	n := getFriendlyAppName(l, "auto-generated-name-1234")
	assert.Equal(t, "real-app-name", n)
}
