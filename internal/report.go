package internal

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Reporter interface {
	PrintReport(ctx context.Context, t *time.Ticker)
	WriteReport(ctx context.Context, reportChan <-chan int) <-chan int
	Total() int
	ReceivedUnique() int
	Duplicates() int
}

type report struct {
	total          int
	receivedUnique int
	duplicates     int
	sync.Mutex
}

func NewReport() Reporter {
	return &report{}
}

func (r *report) PrintReport(ctx context.Context, t *time.Ticker) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				logrus.Debugf("cancel printing report: %v \n", ctx.Err())
				return
			case <-t.C:
				r.Lock()
				logrus.Infof("Received %v unique numbers, %v duplicates. Unique total: %v\n", r.receivedUnique, r.duplicates, r.total)
				r.receivedUnique = 0
				r.duplicates = 0
				r.Unlock()
			}
		}
	}()
}

func (r *report) WriteReport(ctx context.Context, reportChan <-chan int) <-chan int {
	logChan := make(chan int)
	numberIndex := numberSet{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range reportChan {
			select {
			case <-ctx.Done():
				logrus.Debugf("cancel writing report: %v \n", ctx.Err())
				return
			default:
				r.Lock()
				if !numberIndex.has(data) {
					r.total++
					r.receivedUnique++
					numberIndex.add(data)
					logChan <- data
				} else {
					r.duplicates++
				}
				r.Unlock()
			}
		}
	}()
	go func() {
		wg.Wait()
		close(logChan)
	}()
	return logChan
}

func (r *report) Total() int {
	r.Lock()
	total := r.total
	r.Unlock()
	return total
}

func (r *report) ReceivedUnique() int {
	r.Lock()
	receivedUnique := r.receivedUnique
	r.Unlock()
	return receivedUnique
}

func (r *report) Duplicates() int {
	r.Lock()
	duplicates := r.duplicates
	r.Unlock()
	return duplicates
}
