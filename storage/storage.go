package storage

import (
	"encoding/json"
	"sync"
)

func New() *Storage {
	s := &Storage{}
	s.Data = make(map[string]interface{})
	return s
}

type Storage struct {
	Data map[string]interface{}
	sync.RWMutex
}

func (s *Storage) Set(key string, value interface{}) error {
	s.Lock()
	s.Data[key] = value
	s.Unlock()

	return nil
}

func (s *Storage) Get(key string) interface{} {
	return s.Data[key]
}

func (s *Storage) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
