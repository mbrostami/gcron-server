package configs

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// Config keep all config file values
type Config struct {
	Auth struct {
		Username string
		Password string
	}
	Log struct {
		Enable bool
		Level  string
		Path   string
	}
	Server struct {
		Host string
		Port string
	}
}

// GetLogLevel finds the integer map to the log level string in config
func (cfg *Config) GetLogLevel() log.Level {
	if cfg.Log.Level == "trace" {
		return log.TraceLevel
	} else if cfg.Log.Level == "debug" {
		return log.DebugLevel
	} else if cfg.Log.Level == "info" {
		return log.InfoLevel
	} else if cfg.Log.Level == "warning" {
		return log.WarnLevel
	} else if cfg.Log.Level == "error" {
		return log.ErrorLevel
	} else if cfg.Log.Level == "fatal" {
		return log.FatalLevel
	} else {
		return log.PanicLevel
	}
}

// GetConfig returns the configuration map
func GetConfig(flagset *flag.FlagSet) Config {
	var cfg Config
	lviper := readFile(flagset)
	lviper.Unmarshal(&cfg)
	return cfg
}

func readFile(flagset *flag.FlagSet) *viper.Viper {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/gcron/")
	viper.AddConfigPath(".") // FIXME path should be absolute
	pflag.CommandLine.AddGoFlagSet(flagset)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(err)
	}
	return viper.GetViper()
}
