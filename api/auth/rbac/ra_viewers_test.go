package rbac

import (
	"testing"

	"github.com/kubelens/kubelens/api/config"
	"github.com/stretchr/testify/assert"
)

func TestHasApplicationAccess_Viewers_NoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     false,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{"kube-system"},
	}
	ra := RoleAssignment{r}

	result := ra.HasApplicationAccess()

	assert.False(t, result)
}

func TestHasNamespaceAccess_Viewers_NoNamespace(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{"kube-system"},
	}
	ra := RoleAssignment{r}

	result := ra.HasNamespaceAccess("")

	assert.True(t, result)
}

func TestHasNamespaceAccess_Viewers_NoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{"kube-system"},
	}
	ra := RoleAssignment{r}

	result := ra.HasNamespaceAccess("kube-system")

	assert.False(t, result)
}

func TestHasNamespaceAccess_Viewers_NotInExclusionsHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{"kube-system"},
	}
	ra := RoleAssignment{r}

	result := ra.HasNamespaceAccess("kube-public")

	assert.True(t, result)
}

func TestHasEnvVarsAccess_Viewers_MissingMatchLabels(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasEnvVarsAccess(labels)

	assert.False(t, result)
}

func TestHasEnvVarsAccess_Viewers_UnMatchedLabelNoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasEnvVarsAccess(labels)

	assert.False(t, result)
}

func TestHasEnvVarsAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasEnvVarsAccess(labels)

	assert.True(t, result)
}

func TestHasConfigMapAccess_Viewers_MissingMatchLabels(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasConfigMapAccess(labels)

	assert.False(t, result)
}

func TestHasConfigMapAccess_Viewers_UnMatchedLabelNoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasConfigMapAccess(labels)

	assert.False(t, result)
}

func TestHasConfigMapAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasConfigMapAccess(labels)

	assert.True(t, result)
}

func TestHasPodAccess_Viewers_MissingMatchLabelsNotValid(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasPodAccess(labels)

	assert.False(t, result)
}

func TestHasPodAccessAccess_Viewers_UnMatchedLabelViewAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
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

func TestHasPodAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasPodAccess(labels)

	assert.True(t, result)
}

func TestHasLogAccess_Viewers_UnMatchedLabelNoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasLogAccess(labels)

	assert.False(t, result)
}

func TestHasLogAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasLogAccess(labels)

	assert.True(t, result)
}

func TestHasLogAccess_Viewers_MissingMatchLabeNoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasLogAccess(labels)

	assert.False(t, result)
}

func TestHasDeploymentAccess_Viewers_UnMatchedLabelNoAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
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

func TestHasDeploymentAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasDeploymentAccess(labels)

	assert.True(t, result)
}

func TestHasServiceAccess_Viewers_MissingMatchLabels(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasServiceAccess(labels)

	assert.False(t, result)
}

func TestHasServiceAccess_Viewers_UnMatchedLabel(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasServiceAccess(labels)

	assert.True(t, result)
}

func TestHasServiceAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasServiceAccess(labels)

	assert.True(t, result)
}

func TestHasDaemonSetAccess_Viewers_MissingMatchLabels(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasDaemonSetAccess(labels)

	assert.False(t, result)
}

func TestHasDaemonSetAccess_Viewers_UnMatchedLabel(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
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

func TestHasDaemonSetAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasDaemonSetAccess(labels)

	assert.True(t, result)
}

func TestHasJobAccess_Viewers_MissingMatchLabels(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasJobAccess(labels)

	assert.False(t, result)
}

func TestHasJobAccess_Viewers_UnMatchedLabel(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
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

func TestHasJobAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasJobAccess(labels)

	assert.True(t, result)
}

func TestHasReplicaSetAccess_Viewers_MissingMatchLabels(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test2"

	result := ra.HasReplicaSetAccess(labels)

	assert.False(t, result)
}

func TestHasReplicaSetAccess_Viewers_UnMatchedLabel(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
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

func TestHasReplicaSetAccess_Viewers_MatchedLabelHasAccess(t *testing.T) {
	config.C.EnableRBAC = true
	r := Role{
		Operators:   false,
		Viewers:     true,
		MatchLabels: []string{"app=test"},
		Exclusions:  []string{},
	}
	ra := RoleAssignment{r}

	labels := make(map[string]string, 0)
	labels["app"] = "test"

	result := ra.HasReplicaSetAccess(labels)

	assert.True(t, result)
}
