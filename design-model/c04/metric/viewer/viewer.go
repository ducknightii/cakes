package viewer

import (
	"github.com/ducknightii/cakes/design-model/c04/metric/aggregator"
)

type Viewer interface {
	Output(stats map[string]aggregator.RequestStat, startTmp, endTmp int64)
}
