package app

import (
	"context"
	"numberserver/internal"
	"time"
)

type App struct {
	reportPeriodicity time.Duration
	reporter          internal.Reporter
	logger            internal.Logger
	srv               *internal.Server
}

func Init(reportPeriodicity time.Duration, reporter internal.Reporter, logger internal.Logger, srv *internal.Server) App {
	return App{
		reportPeriodicity: reportPeriodicity,
		reporter:          reporter,
		logger:            logger,
		srv:               srv,
	}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	reportChan := make(chan int)
	defer close(reportChan)

	t := time.NewTicker(app.reportPeriodicity)
	logChan := app.reporter.WriteReport(ctx, reportChan)
	app.reporter.PrintReport(ctx, t)
	go app.logger.WriteLog(ctx, logChan)
	app.srv.Run(ctx, cancel, reportChan)
}
