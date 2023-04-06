package aggregator

import (
	"github.com/ducknightii/cakes/design-model/c04/metric/storage"
)

type RequestStat struct {
	Count         int
	AvgResponseMs int32
	MaxMs         int32
}

func Aggregator(allInfos map[string][]storage.RequestInfo) map[string]RequestStat {
	var res = make(map[string]RequestStat)

	for apiName, infos := range allInfos {
		var stat RequestStat
		stat.Count = count(infos)
		stat.AvgResponseMs = avg(infos)
		stat.MaxMs = max(infos)

		res[apiName] = stat
	}
	return res
}

////// 功能拆分 //////

func count(infos []storage.RequestInfo) int {
	return len(infos)
}

func avg(infos []storage.RequestInfo) int32 {
	var avgResponseMs int32
	for _, info := range infos {
		avgResponseMs += info.ResponseMs
		avgResponseMs /= 2
	}

	return avgResponseMs
}

func max(infos []storage.RequestInfo) int32 {
	var maxResponseMs int32
	for _, info := range infos {
		if info.ResponseMs > maxResponseMs {
			maxResponseMs = info.ResponseMs
		}
	}

	return maxResponseMs
}
