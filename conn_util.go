package main

import (
	"net"
	"strings"
)

// Get the bytes of an IP address from a net.Conn
func GetIP(addr net.Addr) net.IP {
	straddr := addr.String()
	host, _, _ := net.SplitHostPort(straddr)
	idx := strings.Index(host, "%")
	if idx >= 0 {
		host = host[:idx]
	}

	ipaddr := net.ParseIP(host)

	ip := ipaddr.To4()
	if ip == nil {
		ip = ipaddr.To16()
	}
	if ip == nil {
		panic("Could not convert IP")
	}
	return ip
}

// Get the bytes of a MAC address from a net.Conn
func GetMAC(addr net.Addr) net.HardwareAddr {
	ifaces, err := net.Interfaces()
	ip := GetIP(addr)
	if err != nil {
		panic("Cannot access interfaces")
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			if strings.Index(addr.String(), ip.String()) >= 0 {
				return iface.HardwareAddr
			}
		}
	}
	return nil
}
