package internal

import (
	"context"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"sync"
)

type Logger interface {
	WriteLog(ctx context.Context, logChan <-chan int)
}

type numbersLog struct {
	name string
}

func NewLog(name string) Logger {
	return numbersLog{name: name}
}

func (nl numbersLog) WriteLog(ctx context.Context, logChan <-chan int) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, _ := os.Create(nl.name)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logrus.Error(err)
				return
			}
		}(file)

		logUtil := log.New(file, "", 0)

		for data := range logChan {
			select {
			case <-ctx.Done():
				logrus.Debugf("cancel writing log: %v \n", ctx.Err())
				return
			default:
				logUtil.Print(data)
			}
		}
	}()
	wg.Wait()
}
