package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/alecthomas/kong"
	"github.com/cyclimse/scwtui/internal/config"
)

type RootState struct {
	config.Config
	Logger  *slog.Logger
	Version string
}

type CLI struct {
	Config config.Config `embed:""`

	Tui TuiCmd `cmd:"" default:"withargs"`
}

func main() {
	var cli CLI

	ctx := kong.Parse(&cli)
	cfg := cli.Config

	logger, err := initLogger(cfg)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	rs := &RootState{
		Config: cfg,
		Logger: logger,
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		rs.Version = info.Main.Version
		logger.Debug("main: using scwtui version",
			slog.String("version", rs.Version),
			slog.String("go_version", info.GoVersion),
		)
	}

	err = ctx.Run(rs)
	ctx.FatalIfErrorf(err)
}

func initLogger(config config.Config) (*slog.Logger, error) {
	w, err := os.OpenFile(config.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	logLevel := config.Logging.Level
	if config.Debug {
		logLevel = slog.LevelDebug
	}

	loggerHandler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: config.Debug,
	})
	return slog.New(loggerHandler), nil
}
