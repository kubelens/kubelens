package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareLabelsDefault(t *testing.T) {
	lbl := make(map[string]string, 0)
	lbl["app"] = "test"

	r := Role{}

	result := r.CompareLabels(lbl, false)

	assert.True(t, result)
}

func TestCompareLabelsMatchLabelsNoLabel(t *testing.T) {
	r := Role{}
	r.MatchLabels = []string{"app=test"}

	result := r.CompareLabels(nil, false)

	assert.False(t, result)
}

func TestCompareLabelsMatchLabelsMissingSplitter(t *testing.T) {
	lbl := make(map[string]string, 0)
	lbl["app"] = "test"

	r := Role{}
	r.MatchLabels = []string{"app=test"}
	r.MatchPrefix = "app"

	result := r.CompareLabels(lbl, false)

	assert.False(t, result)
}

func TestCompareLabelsMatchLabelsNoMatch(t *testing.T) {
	lbl := make(map[string]string, 0)
	lbl["component"] = "test"

	r := Role{}
	r.MatchLabels = []string{"app=test"}
	r.MatchPrefix = "app"
	r.MatchSplitChar = "="

	result := r.CompareLabels(lbl, false)

	assert.False(t, result)
}

func TestCompareLabelsMatchLabelsSuccess(t *testing.T) {
	lbl := make(map[string]string, 0)
	lbl["app"] = "test"

	r := Role{}
	r.MatchLabels = []string{"app=test"}
	r.MatchPrefix = "app"
	r.MatchSplitChar = "="

	result := r.CompareLabels(lbl, false)

	assert.True(t, result)
}

func TestInExclusionsFalse(t *testing.T) {
	r := Role{}
	r.Exclusions = []string{"kube-system", "mongo"}

	result := r.InExclusions("test")

	assert.False(t, result)
}

func TestInExclusionsTrue(t *testing.T) {
	r := Role{}
	r.Exclusions = []string{"mongo"}

	result := r.InExclusions("dev-mongodb-0")

	assert.True(t, result)
}

func TestGetMatchLabels(t *testing.T) {
	r := Role{}
	r.MatchLabels = []string{"app=test"}

	result := r.MatchLabels

	assert.ObjectsAreEqual(r.MatchLabels, result)
}
