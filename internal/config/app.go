package config

import "time"

type App struct {
	LogLevel           string        `default:"debug"`
	TerminationWord    string        `default:"terminate"`
	MessageLenRequired int           `default:"9"`
	ReportPeriodicity  time.Duration `default:"10s"`
	LogName            string        `default:"./numbers.log"`
	Server             Server
}

type Server struct {
	Host    string `default:"0.0.0.0"`
	Port    string `default:"4000"`
	MaxConn int    `default:"5"`
}
