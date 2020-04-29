package env

import (
	"time"

	"github.com/imroc/req"
)

var (
	Local *time.Location
)

func init() {
	Local, _ = time.LoadLocation("Asia/Chongqing")
	time.Local = Local
	req.EnableInsecureTLS(true)
}
