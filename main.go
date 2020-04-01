package main

import (
	"flag"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/mbrostami/gcron-server/configs"
	"github.com/mbrostami/gcron-server/db"
	"github.com/mbrostami/gcron-server/grpc"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Override config file values
	flag.Bool("log.enable", false, "Enable file log")
	flag.String("log.path", "/var/log/gcron/gcron-server.log", "Log file path")
	flag.String("log.level", "warning", "Log level")
	flag.String("server.host", "localhost", "Server host")
	flag.String("server.port", "1400", "Server port")
	cfg := configs.GetConfig(flag.CommandLine)

	log.SetLevel(cfg.GetLogLevel())
	// Setup log
	log.SetFormatter(&nested.Formatter{
		NoColors: false,
	})
	log.SetOutput(os.Stdout)

	var dbAdapter db.DB

	dbAdapter = db.NewLedis()

	//taskCollection := dbAdapter.Get(1446109160, 0, 5)

	grpc.Run(cfg.Server.Host, cfg.Server.Port, dbAdapter)
}
