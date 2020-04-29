package proxy

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
)

var (
	Client = &http.Client{}
)

func init() {
	initSocks5()
}

func initSocks5() {
	var (
		socks5 = os.Getenv("SOCKS5_PROXY")
		auth   *proxy.Auth
		url, _ = url.Parse(socks5)
	)

	if socks5 == "" {
		return
	}

	if url.User != nil {
		auth = &proxy.Auth{}
		auth.User = url.User.Username()
		auth.Password, _ = url.User.Password()
	}

	dialer, err := proxy.SOCKS5("tcp", url.Host, auth, proxy.Direct)
	if err != nil {
		log.Fatalf("go-utils socks5 proxy err: %v\n", err)
	}

	Client.Transport = &http.Transport{Dial: dialer.Dial}
}
