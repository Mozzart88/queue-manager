package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"expat-news/queue-manager/pkg/utils"
	"fmt"
)

type Message struct {
	Id        *int    `json:"id,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
	Msg       *string `json:"msg,omitempty"`
	State     *string `json:"state,omitempty"`
}

func (m *Message) setId(v int) {
	if m.Id == nil {
		m.Id = utils.Ptr(v)
	} else {
		*m.Id = v
	}
}

func (m *Message) setPublisher(v string) {
	if m.Publisher == nil {
		m.Publisher = utils.Ptr(v)
	} else {
		*m.Publisher = v
	}
}

func (m *Message) setMsg(v string) {
	if m.Msg == nil {
		m.Msg = utils.Ptr(v)
	} else {
		*m.Msg = v
	}
}

func (m *Message) setState(v string) {
	if m.State == nil {
		m.State = utils.Ptr(v)
	} else {
		*m.State = v
	}
}

func (m Message) String() string {
	var id, msg, publisher, state string

	if m.Id == nil {
		id = "nil"
	} else {
		id = fmt.Sprintf("%d", *m.Id)
	}
	if m.Msg == nil {
		msg = "nil"
	} else {
		msg = *m.Msg
	}
	if m.Publisher == nil {
		publisher = "nil"
	} else {
		publisher = *m.Publisher
	}
	if m.State == nil {
		state = "nil"
	} else {
		state = *m.State
	}
	return fmt.Sprintf("Publisher{%s %s %s %s}", id, msg, publisher, state)
}

func (m *Message) fillValues(data *repos.QueueMessage) {
	m.setId(data.ID())
	m.setPublisher(data.PublisherName())
	m.setMsg(data.Content())
	m.setState(data.StateName())
}

func (m *Message) Add() error {
	if m.Msg == nil || m.Publisher == nil {
		return fmt.Errorf("required fields id and/or publisher are not defined")
	}
	mu.Lock()
	defer mu.Unlock()
	msg, err := repos.AddQueueMessage(*m.Msg, *m.Publisher)
	if err != nil {
		return err
	}
	if msg == nil {
		return fmt.Errorf("receive nil messge 🤨 thats should not happen at all")
	}
	m.fillValues(msg)
	return nil
}

func (m *Message) SetState() error {
	if m.Id == nil || m.State == nil {
		return fmt.Errorf("missing id and/or state")
	}
	mu.Lock()
	defer mu.Unlock()
	if err := repos.UpdateStateMessage(*m.Id, repos.State_t(*m.State)); err != nil {
		return err
	}
	return nil
}

func (m *Message) SetActive() error {
	m.setState(string(repos.STATE_ACTIVE))
	if err := m.SetState(); err != nil {
		return err
	}
	return nil
}

func (m *Message) SetNew() error {
	m.setState(string(repos.STATE_NEW))
	if err := m.SetState(); err != nil {
		return err
	}
	return nil
}

func (m *Message) Rollback() error {
	return m.SetNew()
}

func (m *Message) SetDone() error {
	if m.Id == nil {
		return fmt.Errorf("can't update state of unsaved message: id is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if err := repos.DeleteMessage(*m.Id); err != nil {
		return err
	}
	m.setState((string)(repos.STATE_DONE))
	m.Id = nil
	return nil
}

func (m *Message) Get() error {
	if m.Id == nil {
		return fmt.Errorf("can't get message: id undefined")
	}
	mu.Lock()
	defer mu.Unlock()
	msg, err := repos.GetUniqQueueMessage(*m.Id)
	if err != nil {
		return err
	}
	if msg == nil {
		return fmt.Errorf("no message in queue with specified id: %d", *m.Id)
	}
	m.fillValues(msg)
	return nil
}

func (m *Message) Delete() error {
	if m.Id == nil {
		return fmt.Errorf("can't delete unsaved message")
	}
	mu.Lock()
	defer mu.Unlock()
	if err := repos.DeleteMessage(*m.Id); err != nil {
		return err
	}
	m.State = nil
	m.Publisher = nil
	m.Msg = nil
	m.Id = nil
	return nil
}
