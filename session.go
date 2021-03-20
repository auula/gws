package sessionx

import (
	"fmt"
	"sync"
	"time"
)

var (
	mux sync.Mutex
	mgr *manager
)

type Session struct {
	ID      string
	Expires time.Time
	Data    map[string]interface{}
}

func (s *Session) Get(key string) (interface{}, error) {
	if len(key) <= 0 {
		return nil, fmt.Errorf("key '%s' invalid", key)
	}
	bytes, err := mgr.store.Reader(s)
	if err != nil {
		return nil, err
	}
	_ = decoder(bytes, s)
	if ele, ok := s.Data[key]; ok {
		return ele, nil
	}
	return nil, fmt.Errorf("key '%s' does not exist", key)
}

func (s *Session) Set(key string, v interface{}) error {
	if len(key) <= 0 {
		return fmt.Errorf("key '%s' invalid", key)
	}
	mux.Lock()
	if s.Data == nil{
		s.Data = make(map[string]interface{}, 8)
	}
	s.Data[key] = v
	mux.Unlock()
	if _, err := mgr.store.Update(s); err != nil {
		return err
	}
	return nil
}

func (s *Session) Remove(key string) error {
	if len(key) <= 0 {
		return fmt.Errorf("key '%s' invalid", key)
	}
	if err := mgr.store.Remove(s, key); err != nil {
		return err
	}
	return nil
}

func (s *Session) Clean() error {
	if err := mgr.store.Delete(s); err != nil {
		return err
	}
	return nil
}

// async code

// type Session struct {
// 	ID      string
// 	Expires time.Time
// 	Data    map[string]interface{}
// 	// 错误一次性收集
// 	tasks []task
// 	collgroup.Group
// }

// func (s *Session) Get(key string) interface{} {
// 	s.tasks = append(s.tasks, func() error {
// 		bytes, err := mgr.store.Reader(s)
// 		if err != nil {
// 			return err
// 		}
// 		Decoder(bytes, s)
// 		return nil
// 	})
// 	if ele, ok := s.Data[key]; ok {
// 		return ele
// 	}
// 	return nil
// }
