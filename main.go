package main

import (
	"flag"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/mbrostami/gcron-server/configs"
	"github.com/mbrostami/gcron-server/db"
	"github.com/mbrostami/gcron-server/grpc"
	"github.com/mbrostami/gcron-server/web"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Override config file values
	flag.Bool("log.enable", false, "Enable file log")
	flag.String("log.path", "/var/log/gcron/gcron-server.log", "Log file path")
	flag.String("log.level", "", "Log level")
	flag.String("server.host", "", "Server host")
	flag.String("server.port", "", "Server port")
	cfg := configs.GetConfig(flag.CommandLine)
	log.SetLevel(cfg.GetLogLevel())

	// Setup log
	log.SetFormatter(&nested.Formatter{
		NoColors: false,
	})
	log.SetOutput(os.Stdout)

	var dbAdapter db.DB
	dbAdapter = db.NewLedis()

	// Run in second thread
	go web.Listen(dbAdapter, cfg)

	// Run in main thread
	//taskCollection := dbAdapter.Get(1446109160, 0, 5)
	grpc.Run(cfg.Server.Host, cfg.Server.Port, dbAdapter)
}
