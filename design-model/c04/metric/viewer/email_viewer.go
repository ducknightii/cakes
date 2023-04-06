package viewer

import (
	"fmt"
	"github.com/ducknightii/cakes/design-model/c04/metric/aggregator"
)

type EmailViewer struct {
	emailSender []string
}

func (e EmailViewer) Output(stats map[string]aggregator.RequestStat, startTmp, endTmp int64) {
	fmt.Printf("[%d-%d] email to [%v] stat: %+v\n", startTmp, endTmp, e.emailSender, stats)
}

func NewEmailViewer(emailSender []string) EmailViewer {
	return EmailViewer{
		emailSender: emailSender,
	}
}
