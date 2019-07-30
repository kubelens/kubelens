package fakes

import (
	"net/http"

	"github.com/kubelens/kubelens/api/k8v1"
)

type SocketFactory struct{}

func (m *SocketFactory) Run() {}

func (m *SocketFactory) Register(k8Client k8v1.Clienter, w http.ResponseWriter, r *http.Request) {}
