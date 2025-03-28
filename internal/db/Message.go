package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"fmt"
)

type Message struct {
	ID        *int    `json:"id,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
	Msg       *string `json:"msg,omitempty"`
	State     *string `json:"state,omitempty"`
}

func (m *Message) fillValues(data *repos.QueueMessage) {
	*m.ID = data.ID()
	*m.Publisher = data.PublisherName()
	*m.Msg = data.Content()
	*m.State = data.StateName()
}

func (m *Message) Add() error {
	if m.Msg == nil || m.Publisher == nil {
		return fmt.Errorf("required fields id and/or publisher are not defined")
	}
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
	if m.ID == nil || m.State == nil {
		return fmt.Errorf("missing id and/or state")
	}
	if err := repos.UpdateStateMessage(*m.ID, repos.State_t(*m.State)); err != nil {
		return err
	}
	return nil
}

func (m *Message) SetActive() error {
	if m.ID == nil {
		return fmt.Errorf("can't update state of unsaved message")
	}
	if err := repos.UpdateStateMessage(*m.ID, repos.STATE_ACTIVE); err != nil {
		return err
	}
	*m.State = (string)(repos.STATE_ACTIVE)
	return nil
}

func (m *Message) SetNew() error {
	if m.ID == nil {
		return fmt.Errorf("can't update state of unsaved message")
	}
	if err := repos.UpdateStateMessage(*m.ID, repos.STATE_NEW); err != nil {
		return err
	}
	*m.State = (string)(repos.STATE_NEW)
	return nil
}

func (m *Message) Rollback() error {
	return m.SetNew()
}

func (m *Message) SetDone() error {
	if m.ID == nil {
		return fmt.Errorf("can't update state of unsaved message")
	}
	if err := repos.DeleteMessage(*m.ID); err != nil {
		return err
	}
	*m.State = (string)(repos.STATE_DONE)
	m.ID = nil
	return nil
}

func (m *Message) Get() error {
	if m.ID == nil {
		return fmt.Errorf("can't get message: id undefined")
	}
	msg, err := repos.GetUniqQueueMessage(*m.ID)
	if err != nil {
		return err
	}
	if msg == nil {
		return fmt.Errorf("no message in queue with specified id: %d", *m.ID)
	}
	m.fillValues(msg)
	return nil
}

func (m *Message) Delete() error {
	if m.ID == nil {
		return fmt.Errorf("can't delete unsaved message")
	}
	if err := repos.DeleteMessage(*m.ID); err != nil {
		return err
	}
	m.State = nil
	m.ID = nil
	return nil
}
