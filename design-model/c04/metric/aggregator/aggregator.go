package aggregator

import (
	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
)

type RequestStat struct {
	Count         int
	AvgResponseMs int32
}

func Aggregator(infos []storage.RequestInfo) RequestStat {
	res := RequestStat{
		Count:         len(infos),
		AvgResponseMs: 0,
	}

	for _, info := range infos {
		res.AvgResponseMs += info.ResponseMs
		res.AvgResponseMs /= 2
	}

	return res
}
