package config

type Config struct {
	Logging struct {
		Level   string `enum:"debug,info,warn,error" default:"info"`
		LogFile string `help:"Log file" default:"/tmp/scwcleaner.log"`
	} `embed:"" prefix:"log-"`
	Scaleway `embed:""`
}

type Scaleway struct {
	Profile string `help:"Scaleway profile name" default:""`
}
