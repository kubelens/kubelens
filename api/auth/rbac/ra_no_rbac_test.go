package rbac

import (
	"testing"

	"github.com/kubelens/kubelens/api/config"
	"github.com/stretchr/testify/assert"
)

func TestHasApplicationAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasApplicationAccess()

	assert.True(t, result)
}

func TestHasNamespaceAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasNamespaceAccess("")

	assert.True(t, result)
}

func TestHasEnvVarsAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasEnvVarsAccess(nil)

	assert.True(t, result)
}

func TestHasPodAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasPodAccess(nil)

	assert.True(t, result)
}

func TestHasLogAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasLogAccess(nil)

	assert.True(t, result)
}

func TestHasDeploymentAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasDeploymentAccess(nil)

	assert.True(t, result)
}

func TestHasServiceAccess_RBACDisabled(t *testing.T) {
	config.C.EnableRBAC = false
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{""},
		Exclusions:  []string{""},
	}
	ra := RoleAssignment{r}

	result := ra.HasServiceAccess(nil)

	assert.True(t, result)
}
