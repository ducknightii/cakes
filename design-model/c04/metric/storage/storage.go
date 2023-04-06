package storage

type Storage interface {
	Save(info RequestInfo) error
	List(apiName string, startTmp, endTmp int64) []RequestInfo
}

type RequestInfo struct {
	ApiName    string
	Timestamp  int64
	ResponseMs int32
}
