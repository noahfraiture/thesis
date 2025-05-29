package main

type Host struct {
	ip   string
	port string
}

var (
	myHost Host
)

func main() {
	// starting
	myHost = Host{
		ip:   "0.0.0.0",
		port: "8080",
	}
	listen(myHost)
}
