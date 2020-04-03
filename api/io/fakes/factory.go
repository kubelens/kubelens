package fakes

import (
	"net/http"

	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
)

type SocketFactory struct{}

func (m *SocketFactory) Run() {}

func (m *SocketFactory) Register(k8Client k8sv1.Clienter, w http.ResponseWriter, r *http.Request) {}
