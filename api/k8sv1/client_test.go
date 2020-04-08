package k8sv1

func setupClient(ns, n string, fail, innerFail bool) Clienter {
	w := &mockWrapper{
		namespace: ns,
		appname:   n,
		fail:      fail,
		innerFail: innerFail,
	}
	return New(w)
}
