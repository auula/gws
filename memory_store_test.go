package sessionx

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	t.Log("session = ", s)
}

func TestReader(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	bytes, err := m.Reader(s)
	if err != nil {
		t.Error(err.Error())
	}
	decoder(bytes, s)
	t.Log("session = ", s)
}

func TestDelete(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	t.Log("session = ", s)
	m.Delete(s)
	_, err := m.Reader(s)
	if err != nil {
		t.Error(err.Error())
	}

}

func TestUpdated(t *testing.T) {
	m := new(memoryStore)
	s := new(Session)
	s.ID = "20210320"
	m.Create(s)
	t.Log("session = ", s)
	v := make(map[string]interface{})
	v["v"] = "test"
	m.Update(&Session{ID: "20210320", Data: v, Expires: time.Now()})
	bytes, err := m.Reader(s)
	if err != nil {
		t.Error(err.Error())
	}
	decoder(bytes, s)
	t.Log("session = ", s)

}
