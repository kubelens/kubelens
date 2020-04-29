package rbac

import (
	"testing"

	"github.com/kubelens/kubelens/api/config"
	"github.com/stretchr/testify/assert"
)

func TestHasEnvVarsAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasEnvVarsAccess(labels)

	assert.True(t, result)
}

func TestHasConfigMapAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasConfigMapAccess(labels)

	assert.True(t, result)
}

func TestHasPodAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasPodAccess(labels)

	assert.True(t, result)
}

func TestHasLogAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasLogAccess(labels)

	assert.True(t, result)
}

func TestHasServiceAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasLogAccess(labels)

	assert.True(t, result)
}

func TestHasDeploymentAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasDeploymentAccess(labels)

	assert.True(t, result)
}

func TestHasDaemonSetAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasDaemonSetAccess(labels)

	assert.True(t, result)
}

func TestHasJobAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasJobAccess(labels)

	assert.True(t, result)
}

func TestHasReplicaSetAccess_Operators_HasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   true,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasReplicaSetAccess(labels)

	assert.True(t, result)
}
