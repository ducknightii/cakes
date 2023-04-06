package collector

import "github.com/ducknightii/cakes/design-model/c04/metric/storage"

type MetricCollector struct {
	storage storage.Storage
}

func (m MetricCollector) Report(info storage.RequestInfo) error {
	return m.storage.Save(info)
}

func NewCollector(storage storage.Storage) MetricCollector {
	return MetricCollector{storage: storage}
}
