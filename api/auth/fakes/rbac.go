package fakes

import "github.com/kubelens/kubelens/api/auth/rbac"

type RoleAssignment struct {
	Fail bool
}

func (ra RoleAssignment) Matches(labels map[string]string, appname *string) bool {
	if ra.Fail {
		return false
	}
	return true
}

func (ra RoleAssignment) CompareLabels(labels map[string]string, exact bool) bool {
	if ra.Fail {
		return false
	}
	return true
}

func (ra RoleAssignment) InExclusions(value string) bool {
	if ra.Fail {
		return false
	}
	return false
}

func (ra RoleAssignment) GetMatchLabels() []string {
	if ra.Fail {
		return []string{"app", "test", "component", "default"}
	}

	return []string{}
}

func (ra RoleAssignment) HasApplicationAccess() bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasNamespaceAccess returns whether or not a user has access to the current namespace
func (ra RoleAssignment) HasNamespaceAccess(namespace string) bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasEnvVarsAccess returns whether or not a user has access to view environment variables by label selectors
func (ra RoleAssignment) HasEnvVarsAccess(labels map[string]string) bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasConfigMapAccess returns whether or not a user has access to view configMap by label selectors
func (ra RoleAssignment) HasConfigMapAccess(labels map[string]string) bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasPodAccess returns whether or not a user has access to view pods by label selectors
func (ra RoleAssignment) HasPodAccess(labels map[string]string) bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasLogAccess returns whether or not a user has access to view/stream logs by label selectors
func (ra RoleAssignment) HasLogAccess(labels map[string]string) bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasDeploymentAccess returns whether or not a user has access to view deployment detail by label selectors
func (ra RoleAssignment) HasDeploymentAccess(labels map[string]string) bool {
	if ra.Fail {
		return false
	}
	return true
}

// HasServiceAccess returns whether or not a user has access to view service detail by label selectors
func (ra RoleAssignment) HasServiceAccess(labels map[string]string) bool {
	if ra.Fail {
		return false
	}
	return true
}

func (ra RoleAssignment) GetRole() rbac.Role {
	return rbac.Role{}
}
