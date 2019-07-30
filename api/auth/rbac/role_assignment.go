package rbac

import (
	"context"

	"github.com/kubelens/kubelens/api/config"
)

type contextKey string

var rbacKey = contextKey("rbac")

// RoleAssignmenter interfaces RoleAssignment
type RoleAssignmenter interface {
	Matches(labels map[string]string, appname *string) bool
	// CompareLabels - if MatchLabels has values, we will assume RBAC includes
	// scoping results to labels matched with the user allowed apps
	CompareLabels(labels map[string]string, exact bool) bool
	// InExclusions checks to see if exclusions list CONTAINS value
	InExclusions(value string) bool
	// GetMatchLabels returns role.MatchLabels
	GetMatchLabels() []string
	// GetRole returns Role
	GetRole() Role
	// HasApplicationAccess returns a boolean indicating an access flag
	HasApplicationAccess() bool
	// HasNamespaceAccess returns whether or not a user has access to the current namespace
	HasNamespaceAccess(namespace string) bool
	// HasEnvVarsAccess returns whether or not a user has access to view environment variables by label selectors
	HasEnvVarsAccess(labels map[string]string) bool
	// HasConfigMapAccess returns whether or not a user has access to view configMap by label selectors
	HasConfigMapAccess(labels map[string]string) bool
	// HasPodAccess returns whether or not a user has access to view pods by label selectors
	HasPodAccess(labels map[string]string) bool
	// HasLogAccess returns whether or not a user has access to view/stream logs by label selectors
	HasLogAccess(labels map[string]string) bool
	// HasDeploymentAccess returns whether or not a user has access to view deployment detail by label selectors
	HasDeploymentAccess(labels map[string]string) bool
	// HasServiceAccess returns whether or not a user has access to view service detail by label selectors
	HasServiceAccess(labels map[string]string) bool
}

/*RoleAssignment provides flags to indicate which features a user has access to.
Operators role objectives:
	- view apps, pods, logs, deployments, services in any namespace
	- view env/configMap becomes only viewable if user has access to that label
	** 	if matchLabels, access will be further restricted by label matches
	**	any skipLabel negates any matchLabel

Viewers role objectives:
	- view apps, deplymens, services in any namespace for apps in any namespace given a matching label
	- view env/configMap for apps in any namespace given a matching label and not a skip label for user
	- view/stream logs for apps in any namespace with their GH team/domain
	** any skipLabel negates any matchLabel
*/
type RoleAssignment struct {
	Role
}

// NewContext adds RoleAssignment to context
func NewContext(ctx context.Context, r RoleAssignmenter) context.Context {
	return context.WithValue(ctx, rbacKey, r)
}

// MustFromContext retrieves a logger instance from context and panics if not found
func MustFromContext(ctx context.Context) RoleAssignmenter {
	instance, ok := ctx.Value(rbacKey).(RoleAssignmenter)
	if !ok {
		panic("rbac was not found in context")
	}
	return instance
}

// HasApplicationAccess returns a boolean indicating an access flag
func (ra RoleAssignment) HasApplicationAccess() bool {
	if config.C.EnableRBAC {
		if ra.Operators || ra.Viewers {
			return true
		}
		return false
	}
	return true
}

// HasNamespaceAccess returns whether or not a user has access to the current namespace
func (ra RoleAssignment) HasNamespaceAccess(namespace string) bool {
	if config.C.EnableRBAC {
		if len(namespace) == 0 {
			return true
		}

		// if scoping access to namespace, check against "skip label" namespaces for the user
		if !ra.Role.InExclusions(namespace) {
			return true
		}
		return false
	}
	return true
}

// HasEnvVarsAccess returns whether or not a user has access to view environment variables by label selectors
func (ra RoleAssignment) HasEnvVarsAccess(labels map[string]string) bool {
	if config.C.EnableRBAC {
		if ra.Role.Viewers && !ra.Role.Operators {
			// if not an operator, labels are required
			if len(ra.Role.MatchLabels) == 0 {
				return false
			}
		}
	}
	return ra.CompareLabels(labels, true)
}

// HasConfigMapAccess returns whether or not a user has access to view configMap by label selectors
func (ra RoleAssignment) HasConfigMapAccess(labels map[string]string) bool {
	if config.C.EnableRBAC {
		if ra.Role.Viewers && !ra.Role.Operators {
			// if not an operator, labels are required
			if len(ra.Role.MatchLabels) == 0 {
				return false
			}
		}
	}
	return ra.CompareLabels(labels, true)
}

// HasPodAccess returns whether or not a user has access to view pods by label selectors
func (ra RoleAssignment) HasPodAccess(labels map[string]string) bool {
	if config.C.EnableRBAC {
		if ra.Role.Viewers && !ra.Role.Operators {
			// if not an operator, labels are required
			if len(ra.Role.MatchLabels) == 0 {
				return false
			}
		}
	}
	return ra.CompareLabels(labels, false)
}

// HasLogAccess returns whether or not a user has access to view/stream logs by label selectors
func (ra RoleAssignment) HasLogAccess(labels map[string]string) bool {
	if config.C.EnableRBAC {
		if ra.Role.Viewers && !ra.Role.Operators {
			// if not an operator, labels are required
			if len(ra.Role.MatchLabels) == 0 {
				return false
			}
		}
	}
	return ra.CompareLabels(labels, true)
}

// HasDeploymentAccess returns whether or not a user has access to view deployment detail by label selectors
func (ra RoleAssignment) HasDeploymentAccess(labels map[string]string) bool {
	if config.C.EnableRBAC {
		if ra.Role.Viewers && !ra.Role.Operators {
			// if not an operator, labels are required
			if len(ra.Role.MatchLabels) == 0 {
				return false
			}
		}
	}
	return ra.CompareLabels(labels, false)
}

// HasServiceAccess returns whether or not a user has access to view service detail by label selectors
func (ra RoleAssignment) HasServiceAccess(labels map[string]string) bool {
	if config.C.EnableRBAC {
		if ra.Role.Viewers && !ra.Role.Operators {
			// if not an operator, labels are required
			if len(ra.Role.MatchLabels) == 0 {
				return false
			}
		}
	}
	return ra.CompareLabels(labels, false)
}

// GetRole returns Role
func (ra RoleAssignment) GetRole() Role {
	return ra.Role
}
