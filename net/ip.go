package net

import (
	"net"
)

var (
	privateIPNets []*net.IPNet
)

func init() {
	AddPrivateIPNets(
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
		"fd00::/8",
	)
}

func AddPrivateIPNets(ips ...string) {
	for _, ip := range ips {
		if _, ipnet, err := net.ParseCIDR(ip); err == nil {
			privateIPNets = append(privateIPNets, ipnet)
		}
	}
}

func IsPrivateIP(ip net.IP) bool {
	for _, ipnet := range privateIPNets {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

func GetLocalIPs() []net.IP {
	nets, err := net.Interfaces()
	if err != nil {
		return nil
	}

	ips := []net.IP{}

	for _, n := range nets {
		if n.Flags&net.FlagUp != net.FlagUp {
			continue
		}

		inAddrs, _ := n.Addrs()
		for _, a := range inAddrs {
			ip := GetIPFromAddr(a)
			if ip == nil {
				continue
			}

			if ip.IsLoopback() {
				continue
			}

			if ip.IsUnspecified() {
				continue
			}

			ips = append(ips, ip)
		}
	}

	return ips
}

func GetIPFromAddr(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPAddr:
		return v.IP
	case *net.IPNet:
		return v.IP
	default:
		return nil
	}
}

// FormatIP format "0.0.0.0" to real private ip, like: "192.168.xx.xx"
func FormatIP(s string) string {
	ip := net.ParseIP(s)
	if ip == nil {
		ip = net.IPv4zero
	}

	if !ip.IsUnspecified() {
		return ip.String()
	}

	for _, ip := range GetLocalIPs() {
		if IsPrivateIP(ip) {
			return ip.String()
		}
	}

	return ip.String()
}
