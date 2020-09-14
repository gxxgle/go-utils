package http

import (
	"testing"
)

func TestGet(t *testing.T) {
	resp, err := C.R().Get("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		t.Fatal("http get failed, err:", err)
	}

	t.Log("http response:", resp.String())
}
