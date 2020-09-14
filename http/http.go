package http

import (
	"crypto/tls"

	"github.com/gxxgle/go-utils/time"

	"github.com/go-resty/resty/v2"
)

var (
	C = resty.New()
)

func init() {
	C.SetTimeout(time.Second * 30)
	C.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
