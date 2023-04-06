package main

import (
	"github.com/ducknightii/cakes/design-model/c04/metric/viewer"
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

	consoleViewer := viewer.NewConsoleViewer()

	consoleReporterIns := reporter.NewConsoleReporter(storageIns, consoleViewer)
	consoleReporterIns.StartReport()

	emailViewer := viewer.NewEmailViewer([]string{"ducknightii@gmail.com"})
	emailReporterIns := reporter.NewEmailReporter(storageIns, emailViewer)
	emailReporterIns.StartReport()

	time.Sleep(time.Minute * 10)
}
