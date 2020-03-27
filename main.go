package main

import (
	"flag"

	"github.com/mbrostami/gcron-server/configs"
	"github.com/mbrostami/gcron-server/grpc"
	"github.com/mbrostami/gcron-server/server"
)

func main() {
	// Override config file values
	flag.Bool("out.notime", false, "Clean output")
	flag.Bool("out.clean", false, "Clean output")
	flag.String("server.host", "localhost", "Server host")
	flag.String("server.port", "1400", "Server port")
	flag.String("server.protocol", "grpc", "Protocol (tcp/udp/unix/grpc)")
	flag.String("server.unix.socket", "/tmp/gcron-server.sock", "UNIX socket path")
	cfg := configs.GetConfig(".", flag.CommandLine)
	flag.Parse()

	if cfg.Server.Protocol == "unix" {
		server.ListenUNIX(cfg.Server.Unix.Socket)
	} else if cfg.Server.Protocol == "tcp" {
		server.ListenTCP(cfg.Server.Host, cfg.Server.Port)
	} else if cfg.Server.Protocol == "udp" {
		server.ListenUDP(cfg.Server.Host, cfg.Server.Port)
	} else if cfg.Server.Protocol == "grpc" {
		grpc.Run(cfg.Server.Host, cfg.Server.Port)
	}
}
