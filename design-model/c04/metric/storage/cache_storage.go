package storage

import "sync"

type CacheStorage struct {
	records sync.Map // apiName => []Info
}

func (c *CacheStorage) Save(info RequestInfo) error {
	var infos []RequestInfo
	_infos, ok := c.records.Load(info.ApiName)
	if ok {
		infos = _infos.([]RequestInfo)
	} else {
		infos = make([]RequestInfo, 0)
	}
	infos = append(infos, info)
	c.records.Store(info.ApiName, infos)

	return nil
}

func (c *CacheStorage) List(apiName string, startTmp, endTmp int64) []RequestInfo {
	// todo filter

	var infos []RequestInfo
	_infos, ok := c.records.Load(apiName)
	if ok {
		infos = _infos.([]RequestInfo)
	}

	return infos
}

func (c *CacheStorage) AllList(startTmp, endTmp int64) map[string][]RequestInfo {
	// todo filter

	var allInfos = make(map[string][]RequestInfo)
	c.records.Range(func(key, value any) bool {
		infos := value.([]RequestInfo)
		apiName := key.(string)

		allInfos[apiName] = infos

		return true
	})

	return allInfos
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{}
}
