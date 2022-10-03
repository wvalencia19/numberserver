package internal

import (
	"bufio"
	"context"
	"os"

	"github.com/stretchr/testify/assert"

	"path"
	"strconv"
	"sync"
	"testing"
)

const logFile = "numbers.log"

func TestNumbersLog_WriteLog(t *testing.T) {
	tests := []struct {
		name  string
		ctx   context.Context
		input []int
		want  []int
	}{
		{
			name:  "valid creation",
			ctx:   context.Background(),
			input: []int{0, 1, 2},
			want:  []int{0, 1, 2},
		},
		{
			name:  "with canceled context",
			ctx:   canceledContext(),
			input: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			name := path.Join(tmpdir, logFile)
			logger := NewLog(name)
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				logger.WriteLog(tt.ctx, populateChannel(tt.input...))
			}()
			wg.Wait()
			file, err := os.Open(name)
			assert.Nil(t, err)
			scanner := bufio.NewScanner(file)

			var rows []int

			for scanner.Scan() {
				currentLineText := scanner.Text()
				val, _ := strconv.Atoi(currentLineText)
				rows = append(rows, val)
			}
			assert.Nil(t, err)
			assert.Equal(t, rows, tt.want)
		})
	}
}
