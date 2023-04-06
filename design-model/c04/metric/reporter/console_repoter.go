package reporter

import (
	"fmt"
	"time"

	"github.com/ducknightii/cakes/design-model/c04/metric/aggregator"
	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
	cron "github.com/robfig/cron/v3"
)

type ConsoleReporter struct {
	storage storage.Storage
}

func (c ConsoleReporter) StartReport() {
	cronIns := cron.New()

	_, err := cronIns.AddFunc("*/1 * * * *", func() {
		now := time.Now()
		startAt := now.Add(-1 * time.Hour)

		requestInfos := c.storage.List("test", startAt.Unix(), now.Unix())

		stat := aggregator.Aggregator(requestInfos)

		fmt.Printf("[%s-%s] stat: %+v\n", startAt, now, stat)

	})
	if err != nil {
		panic(err)
	}

	cronIns.Start()
}

func NewConsoleReporter(storage storage.Storage) ConsoleReporter {
	return ConsoleReporter{storage: storage}
}
