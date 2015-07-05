package storage

import (
	"encoding/json"
	"strings"
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

func (s *Storage) ToJson(prefix string) ([]byte, error) {
	newData := make(map[string]interface{})

	for key, value := range s.Data {
		if strings.HasPrefix(key, prefix) {
			newData[key] = value
		}
	}
	return json.Marshal(newData)
}
