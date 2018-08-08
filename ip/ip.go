package ip

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
)

// ExtranetIP return external IP address.
func ExtranetIP() (string, error) {
	resp, err := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	idx := bytes.Index(bs, []byte(`"cip": "`))
	bs = bs[idx+len(`"cip": "`):]
	idx = bytes.Index(bs, []byte(`"`))
	bs = bs[:idx]
	return string(bs), nil
}

// IntranetIPs return internal IP address.
func IntranetIPs() ([]string, error) {
	out := make([]string, 0)
	is, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range is {
		// interface down
		if i.Flags&net.FlagUp == 0 {
			continue
		}

		// loopback interface
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			out = append(out, ip.String())
		}
	}

	return out, nil
}
