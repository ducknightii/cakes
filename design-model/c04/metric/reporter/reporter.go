package reporter

import (
	"github.com/ducknightii/cakes/design-model/c04/metric/aggregator"
	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
	"github.com/ducknightii/cakes/design-model/c04/metric/viewer"
)

type reporter struct {
	storage storage.Storage
	viewer  viewer.Viewer
}

// notes 统一流程抽象
func (r reporter) doReport(startTmp, endTmp int64) {
	requestInfos := r.storage.AllList(startTmp, endTmp)

	stats := aggregator.Aggregator(requestInfos)

	r.viewer.Output(stats, startTmp, endTmp)
}
