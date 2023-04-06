package reporter

import (
	"github.com/ducknightii/cakes/design-model/c04/metric/viewer"
	"time"

	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
	cron "github.com/robfig/cron/v3"
)

type EmailReporter struct {
	reporter
}

func (e EmailReporter) StartReport() {
	cronIns := cron.New()

	_, err := cronIns.AddFunc("0 */1 * * *", func() {
		now := time.Now()
		startAt := now.Add(-1 * time.Hour)

		e.doReport(startAt.Unix(), now.Unix())
	})
	if err != nil {
		panic(err)
	}

	cronIns.Start()
}

func NewEmailReporter(storage storage.Storage, viewer viewer.Viewer) EmailReporter {
	return EmailReporter{
		reporter{
			storage: storage,
			viewer:  viewer,
		},
	}
}
