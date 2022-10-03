package main

import (
	"numberserver/internal"
	"numberserver/internal/app"
	"numberserver/internal/config"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	var conf config.App
	err := envconfig.Process("NumberServer", &conf)
	if err != nil {
		log.Panic(err.Error())
	}

	l := parseLogLevel(conf.LogLevel)
	log.SetLevel(l)

	reporter := internal.NewReport()
	logger := internal.NewLog(conf.LogName)
	validator := internal.NewNumberValidator(conf.TerminationWord, conf.MessageLenRequired)
	srv := internal.New(config.Server{
		Host:    conf.Server.Host,
		Port:    conf.Server.Port,
		MaxConn: conf.Server.MaxConn,
	}, validator)

	numberServer := app.Init(conf.ReportPeriodicity, reporter, logger, srv)
	numberServer.Run()
}

func parseLogLevel(level string) log.Level {
	l, err := log.ParseLevel(level)
	if err != nil {
		l = log.InfoLevel
		log.Errorf("setting log info %v", err)
	}
	return l
}
