package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport_WriteReport(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		populateWith  []int
		totalReport   int
		receiveUnique int
		duplicates    int
		prepare       func(total []int) <-chan int
		want          []int
	}{
		{
			name:          "with duplicated data",
			ctx:           context.Background(),
			populateWith:  []int{1, 0, 0, 0},
			totalReport:   2,
			receiveUnique: 2,
			duplicates:    2,
			prepare: func(values []int) <-chan int {
				return populateChannel(values...)
			},
			want: []int{1, 0},
		},
		{
			name:          "without duplicated data",
			ctx:           context.Background(),
			populateWith:  []int{0, 1, 2},
			totalReport:   3,
			receiveUnique: 3,
			duplicates:    0,
			prepare: func(values []int) <-chan int {
				return populateChannel(values...)
			},
			want: []int{0, 1, 2},
		},

		{
			name:          "with canceled context",
			ctx:           canceledContext(),
			populateWith:  []int{0, 1, 2},
			totalReport:   0,
			receiveUnique: 0,
			duplicates:    0,
			prepare: func(values []int) <-chan int {
				return populateChannel(values...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportChan := tt.prepare(tt.populateWith)
			reporter := NewReport()
			resultChan := reporter.WriteReport(tt.ctx, reportChan)
			var result []int

			for data := range resultChan {
				result = append(result, data)
			}

			assert.Equal(t, reporter.Total(), tt.totalReport)
			assert.Equal(t, reporter.ReceivedUnique(), tt.receiveUnique)
			assert.Equal(t, reporter.Duplicates(), tt.duplicates)
			assert.Equal(t, result, tt.want)
		})
	}
}

func populateChannel(values ...int) <-chan int {
	messageChan := make(chan int, 10)
	defer close(messageChan)
	for _, value := range values {
		messageChan <- value
	}
	return messageChan
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
