package db

import (
	"sync"
)

var mu sync.Mutex

type Message struct {
	ID        *int    `json:"id,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
	Msg       *string `json:"msg,omitempty"`
	State     *string `json:"state,omitempty"`
}

func (m *Message) Insert() error {
	mu.Lock()
	defer mu.Unlock()

	return nil
}

func (m *Message) SetState() error {
	if *m.State == "done" {
		return m.Delete()
	}
	mu.Lock()
	defer mu.Unlock()
	return nil
}

func (m *Message) Get() error {
	mu.Lock()
	defer mu.Unlock()

	return nil
}

func (m *Message) Delete() error {
	mu.Lock()
	defer mu.Unlock()
	return nil
}
