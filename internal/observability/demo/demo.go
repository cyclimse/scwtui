package demo

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cyclimse/scwtui/internal/resource"
)

const (
	// maxNumLogs is the maximum number of logs that can be returned by the demo monitor.
	maxNumLogs = 100
)

func NewDemo() *Demo {
	return &Demo{}
}

func (d *Demo) Logs(_ context.Context, _ resource.Resource) ([]resource.Log, error) {
	numLogs := gofakeit.Number(1, maxNumLogs)

	logs := make([]resource.Log, 0, numLogs)
	begin := gofakeit.DateRange(time.Now().Add(-time.Hour*24*7), time.Now())

	for i := 0; i < numLogs; i++ {
		logs = append(logs, resource.Log{
			Timestamp: begin,
			Line:      gofakeit.Sentence(gofakeit.Number(1, 10)),
		})
		begin = begin.Add(time.Duration(gofakeit.Number(1, 10)) * time.Minute)
	}

	return logs, nil
}

type Demo struct{}
