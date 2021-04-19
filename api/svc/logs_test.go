package svc

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	klog "github.com/kubelens/kubelens/api/log"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestPodLogs(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", "/logs/test?namespace=default&tail=1000", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Logs(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b k8sv1.Log
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, len(b.Output) > 0)
}
