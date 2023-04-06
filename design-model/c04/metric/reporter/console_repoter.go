package reporter

import (
	"time"

	"github.com/ducknightii/cakes/design-model/c04/metric/viewer"

	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
	cron "github.com/robfig/cron/v3"
)

// notes: Reporter 为组装类 负责串联 Storage->Aggregator->Viewer 流程

type ConsoleReporter struct {
	reporter
}

// StartReport notes 本部分组合逻辑足够简单 不必写单元测试
func (c ConsoleReporter) StartReport() {
	cronIns := cron.New()

	_, err := cronIns.AddFunc("*/1 * * * *", func() {
		now := time.Now()
		startAt := now.Add(-1 * time.Hour)

		c.doReport(startAt.Unix(), now.Unix())

	})
	if err != nil {
		panic(err)
	}

	cronIns.Start()
}

func NewConsoleReporter(storage storage.Storage, viewer viewer.Viewer) ConsoleReporter {
	return ConsoleReporter{reporter{
		storage: storage,
		viewer:  viewer,
	}}
}
