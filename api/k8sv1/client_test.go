package k8sv1

import "github.com/kubelens/kubelens/api/config"

func setupClient(ns, n string, fail, innerFail bool) Clienter {
	config.Set("../testdata/mock_config.json")
	config.C.EnableAuth = false
	w := &mockWrapper{
		namespace: ns,
		appname:   n,
		fail:      fail,
		innerFail: innerFail,
	}
	return New(w)
}
