package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"expat-news/queue-manager/pkg/utils"
	"fmt"
)

type Publisher struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func (p *Publisher) setId(v int) {
	if p.Id == nil {
		p.Id = utils.Ptr(v)
	} else {
		*p.Id = v
	}
}

func (p *Publisher) setName(v string) {
	if p.Name == nil {
		p.Name = utils.Ptr(v)
	} else {
		*p.Name = v
	}
}

func (p Publisher) String() string {
	var id, name string

	if p.Id == nil {
		id = "nil"
	} else {
		id = fmt.Sprintf("%d", *p.Id)
	}
	if p.Name == nil {
		name = "nil"
	} else {
		name = *p.Name
	}
	return fmt.Sprintf("Publisher{%s %s}", id, name)
}

func (p *Publisher) Get() error {
	mu.Lock()
	defer mu.Unlock()
	data, err := repos.GetPublisher(p.Id, p.Name)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("unregistered publisher: %v", *p)
	}
	p.setId(data.ID())
	p.setName(data.Name())
	return nil
}

func (p *Publisher) Update(newName string) error {
	if p.Id == nil {
		return fmt.Errorf("id is undefined")
	}
	mu.Lock()
	defer mu.Unlock()
	if err := repos.UpdatePublisher(*p.Id, newName); err != nil {
		return err
	}
	p.setName(newName)
	return nil
}

func (p *Publisher) Delete() error {
	mu.Lock()
	defer mu.Unlock()
	if err := repos.DeletePublisher(p.Id, p.Name); err != nil {
		return err
	}
	p.Name = nil
	p.Id = nil
	return nil
}

func (p *Publisher) Register() error {
	if p.Name == nil {
		return fmt.Errorf("name is undefined")
	}
	mu.Lock()
	defer mu.Unlock()
	id, err := repos.AddPublisher(*p.Name)
	if err != nil {
		return err
	}
	p.setId(id)
	return nil
}
