package database

type StorageContainer struct {
	data map[string]interface{}
}

func (s *StorageContainer) Get(key string) interface{} {
	return s.data[key]
}

func (s *StorageContainer) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *StorageContainer) Delete(key string) {
	delete(s.data, key)
}

func (s *StorageContainer) Clear() {
	s.data = make(map[string]interface{})
}
