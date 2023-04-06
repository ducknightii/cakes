package viewer

import (
	"fmt"
	"github.com/ducknightii/cakes/design-model/c04/metric/aggregator"
)

type ConsoleViewer struct{}

func (c ConsoleViewer) Output(stats map[string]aggregator.RequestStat, startTmp, endTmp int64) {
	fmt.Printf("[%d-%d] stat: %+v\n", startTmp, endTmp, stats)
}

func NewConsoleViewer() ConsoleViewer {
	return ConsoleViewer{}
}
