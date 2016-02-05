package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/oleksandr/bonjour"
)

const (
	port  = 33333
	iface = "eth0"
	name  = "iTuned"
)

var (
	txt = []string{
		"txtvers=1",
		"pw=false",
		"tp=UDP",
		"sm=false",
		"ek=1",
		"cn=0,1",
		"ch=2",
		"ss=16",
		"sr=44100",
		"vn=3",
		"et=0,1",
	}
)

func main() {
	fmt.Println("hello")

	iface, err := net.InterfaceByName(iface)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	bName := hex.EncodeToString(iface.HardwareAddr) + "@" + name
	s, err := bonjour.Register(bName, "_raop._tcp", "", port, txt, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)
	for sig := range handler {
		if sig == os.Interrupt {
			s.Shutdown()
			time.Sleep(100 * time.Millisecond)
			break
		}
	}
}
