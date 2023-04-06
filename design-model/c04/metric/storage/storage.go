package storage

type Storage interface {
	Save(info RequestInfo) error
	List(apiName string, startTmp, endTmp int64) []RequestInfo
	AllList(startTmp, endTmp int64) map[string][]RequestInfo
}

type RequestInfo struct {
	ApiName    string
	Timestamp  int64
	ResponseMs int32
}
