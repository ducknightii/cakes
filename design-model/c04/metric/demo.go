package main

import (
	"time"

	"github.com/ducknightii/cakes/design-model/c04/metric/collector"
	"github.com/ducknightii/cakes/design-model/c04/metric/reporter"
	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
)

func main() {
	storageIns := storage.NewCacheStorage()

	collectorIns := collector.NewCollector(storageIns)
	_ = collectorIns.Report(storage.RequestInfo{
		ApiName:    "test",
		Timestamp:  time.Now().UnixMilli(),
		ResponseMs: 100,
	})
	_ = collectorIns.Report(storage.RequestInfo{
		ApiName:    "test",
		Timestamp:  time.Now().UnixMilli(),
		ResponseMs: 200,
	})

	consoleReporterIns := reporter.NewConsoleReporter(storageIns)
	consoleReporterIns.StartReport()

	emailReporterIns := reporter.NewEmailReporter(storageIns, []string{"ducknightii@gmail.com"})
	emailReporterIns.StartReport()

	time.Sleep(time.Minute * 10)
}
