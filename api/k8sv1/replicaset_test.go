package k8sv1

import (
	"context"
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestReplicaSetsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "rs1", false, false)

	d, err := c.ReplicaSets(ReplicaSetOptions{
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "rs1",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetReplicaSetsDefaultFail(t *testing.T) {
	c := setupClient("testns", "rs2", true, true)

	_, err := c.ReplicaSets(ReplicaSetOptions{
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "rs2",
		Context:    context.Background(),
	})

	assert.NotNil(t, err)
}

func TestReplicaSetDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "rs3", false, false)

	d, err := c.ReplicaSet(ReplicaSetOptions{
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		Name:       "rs3",
		LinkedName: "rs3link",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, d.Namespace, "testns")
}

func TestGetReplicaSetDefaultFail(t *testing.T) {
	c := setupClient("testns", "rs4", true, true)

	_, err := c.ReplicaSet(ReplicaSetOptions{
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		Name:       "rs4",
		LinkedName: "rs4link",
		Context:    context.Background(),
	})

	assert.NotNil(t, err)
}
