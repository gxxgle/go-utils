package main

import (
	"log"

	"github.com/gxxgle/go-utils/ip"
)

func main() {
	log.Println(ip.ExtranetIP())
	log.Println(ip.IntranetIPs())
}
