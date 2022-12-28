package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	v3 "github.com/tony2001/go-vulners/api/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger

	config Config
)

type CommandRunner interface {
	Init([]string) error
	Run() error
	Name() string
}

func newLogger(level string, format string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	_ = cfg.Level.UnmarshalText([]byte(level))
	if strings.EqualFold(format, "json") {
		cfg.Encoding = "json"
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return cfg.Build()
}

func initFlags() {
	fset := flag.NewFlagSet("vulners", flag.ContinueOnError)
	fset.Usage = cleanenv.FUsage(fset.Output(), &config, nil, fset.Usage)
	fset.Parse(os.Args[1:])
}

func initLogger() *zap.SugaredLogger {
	logger, err := newLogger("debug", "console")
	if err != nil {
		log.Fatalf("failed to initialize logger")
	}

	return logger.Sugar()
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		logger.Fatalw("failed to initialize config", "error", err)
	}
}

func main() {
	initFlags()
	logger = initLogger()
	initConfig()

	fmt.Printf("%+v\n", config)

	client, err := v3.NewClientWithResponses(config.Server)
	if err != nil {
		logger.Fatalw("failed ot initialize Client", "error", err)
	}

	commands := []CommandRunner{
		NewSearchCommand(config, client),
	}

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	subcommand := os.Args[1]
	for _, cmd := range commands {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			err := cmd.Run()
			if err != nil {
				logger.Errorw("subcommand failed", "subcommand", subcommand, "error", err)
				os.Exit(1)
			}
			return
		}
	}

	logger.Errorw("unknown subcommand", "subcommand", subcommand)
	os.Exit(1)
}
