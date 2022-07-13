package transport

import "sync"

type SyncMapSocketMeta struct {
	store sync.Map
}

func (m *SyncMapSocketMeta) MetaLoad(key string) (interface{}, bool) {
	return m.store.Load(key)
}

func (m *SyncMapSocketMeta) MetaStore(key string, value interface{}) {
	m.store.Store(key, value)
}

func (m *SyncMapSocketMeta) MetaDelete(key string) {
	m.store.Delete(key)
}
