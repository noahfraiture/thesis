package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Host struct {
	ip   string
	port string
}

func sleep(t time.Duration) {
	sleeping = true
	time.Sleep(t * time.Second)
	sleeping = false
}

var (
	myHost     Host
	parentHost Host
	childsHost []Host
	subnet     [3]string
	sleeping   bool
	c2Addr     string
)

func main() {
	// starting
	myHost = Host{ip: os.Getenv("MY_IP"), port: os.Getenv("MY_PORT")}
	subnet = [3]string(strings.Split(myHost.ip, "."))
	parentHost = Host{ip: os.Getenv("PARENT_IP"), port: os.Getenv("PARENT_PORT")}
	c2Addr = os.Getenv("C2_ADDR")
	fmt.Println(myHost)
	fmt.Println(parentHost)
	fmt.Println(c2Addr)
	// Listen for incoming packet from child
	go poll(net.JoinHostPort(parentHost.ip, parentHost.port), net.JoinHostPort(myHost.ip, myHost.port))
	listen(myHost)
}
