package reporter

import (
	"fmt"
	"time"

	"github.com/ducknightii/cakes/design-model/c04/metric/aggregator"
	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
	cron "github.com/robfig/cron/v3"
)

type EmailReporter struct {
	storage     storage.Storage
	emailSender []string
}

func (e EmailReporter) StartReport() {
	cronIns := cron.New()

	_, err := cronIns.AddFunc("0 */1 * * *", func() {
		now := time.Now()
		startAt := now.Add(-1 * time.Hour)

		requestInfos := e.storage.List("test", startAt.Unix(), now.Unix())

		stat := aggregator.Aggregator(requestInfos)

		fmt.Printf("[%s-%s] email to [%v] stat: %+v\n", startAt, now, e.emailSender, stat)

	})
	if err != nil {
		panic(err)
	}

	cronIns.Start()
}

func NewEmailReporter(storage storage.Storage, emailSender []string) EmailReporter {
	return EmailReporter{
		storage:     storage,
		emailSender: emailSender,
	}
}
