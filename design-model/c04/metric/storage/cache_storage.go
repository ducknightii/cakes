package storage

type CacheStorage struct {
	records map[string][]RequestInfo
}

func (c CacheStorage) Save(info RequestInfo) error {
	c.records[info.ApiName] = append(c.records[info.ApiName], info)

	return nil
}

func (c CacheStorage) List(apiName string, startTmp, endTmp int64) []RequestInfo {
	// todo filter

	return c.records[apiName]
}

func NewCacheStorage() CacheStorage {
	return CacheStorage{
		records: make(map[string][]RequestInfo),
	}
}
