package main

import (
	"flag"
	"gcron-server/server"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host")
	port := flag.String("port", "1400", "Port number")
	protocolType := flag.String("prot", "tcp", "Protocol (tcp/udp/unix)")
	flag.Parse()
	if *protocolType == "unix" {
		path := "/tmp/gcron.sock"
		server.ListenUNIX(path)
	} else if *protocolType == "tcp" {
		server.ListenTCP(*host, *port)
	} else if *protocolType == "udp" {
		server.ListenUDP(*host, *port)
	}
}