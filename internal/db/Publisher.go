package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"fmt"
)

type Publisher struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func (p *Publisher) Get() error {
	data, err := repos.GetPublisher(p.Id, p.Name)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("unregistered publisher: %s", *p.Name)
	}
	*p.Id = data.ID()
	*p.Name = data.Name()
	return nil
}

func (p *Publisher) Update(newName string) error {
	if err := repos.UpdatePublisher(*p.Id, newName); err != nil {
		return err
	}
	*p.Name = newName
	return nil
}

func (p *Publisher) Delete() error {
	if err := repos.DeletePublisher(p.Id, p.Name); err != nil {
		return err
	}
	p.Name = nil
	p.Id = nil
	return nil
}

func (p *Publisher) Register() error {
	id, err := repos.AddPublisher(*p.Name)
	if err != nil {
		return err
	}
	*p.Id = id
	return nil
}
