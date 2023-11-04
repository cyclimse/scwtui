package main

import (
	"github.com/alecthomas/kong"
	"github.com/cyclimse/scaleway-dangling/internal/config"
)

type CmdContext struct {
	config.Config
}

type CLI struct {
	Config config.Config `embed:""`
	Tui    TuiCmd        `cmd:"" default:"withargs"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run(&CmdContext{Config: cli.Config})
	ctx.FatalIfErrorf(err)
}
