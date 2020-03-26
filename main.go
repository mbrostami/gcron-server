package main

import (
	"flag"

	"github.com/mbrostami/gcron-server/grpc"
	"github.com/mbrostami/gcron-server/server"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host")
	port := flag.String("port", "1400", "Port number")
	protocolType := flag.String("prot", "tcp", "Protocol (tcp/udp/unix)")
	//configpath := flag.String("config", ".", "Config path")

	// Override config file values
	flag.Bool("out.notime", false, "Clean output")
	flag.Bool("out.clean", false, "Clean output")
	flag.String("server.tcp.port", "", "TCP Server port")
	flag.String("server.tcp.host", "", "TCP Server host")
	flag.String("server.udp.port", "", "UDP Server port")
	flag.String("server.udp.host", "", "UDP Server host")
	flag.String("server.unix", "/tmp/gcron-server.sock", "UNIX socket path")

	flag.String("server.protocol", "unix", "Protocol (tcp/udp/unix)")

	//cfg := configs.GetConfig(*configpath, flag.CommandLine)
	flag.Parse()
	if *protocolType == "unix" {
		path := "/tmp/gcron-server.sock"
		server.ListenUNIX(path)
	} else if *protocolType == "tcp" {
		server.ListenTCP(*host, *port)
	} else if *protocolType == "udp" {
		server.ListenUDP(*host, *port)
	} else if *protocolType == "grpc" {
		grpc.Run()
	}
}
