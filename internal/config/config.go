package config

import "log/slog"

type Config struct {
	Logging struct {
		Level slog.Level `default:"INFO"            enum:"DEBUG,INFO,WARN,ERROR" help:"Log level"`
		File  string     `default:"/tmp/scwtui.log" help:"File to write logs to"`
	} `embed:"" prefix:"log-"`
	Scaleway `embed:""`
	Tui      `embed:"" prefix:"ui-"`
}

type Scaleway struct {
	Profile string `default:"" help:"Scaleway profile name"`
}

type Tui struct {
	Theme string `default:"monokai" help:"The theme to use for syntax highlighting."`
}
