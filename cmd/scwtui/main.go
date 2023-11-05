package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/cyclimse/scwtui/internal/config"
)

type CmdContext struct {
	config.Config

	Logger *slog.Logger
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

	err = ctx.Run(&CmdContext{Config: cfg, Logger: logger})
	ctx.FatalIfErrorf(err)
}

func initLogger(config config.Config) (*slog.Logger, error) {
	w, err := os.OpenFile(config.Logging.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	loggerHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: config.Logging.Level,
	})
	return slog.New(loggerHandler), nil
}
