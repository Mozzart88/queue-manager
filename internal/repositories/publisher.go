package repos

import (
	"expat-news/queue-manager/internal/repositories/crud"
	"fmt"
)

const publisherTable = "publisher"

type Publisher struct {
	id   int
	name string
}

func (p *Publisher) ID() int {
	return p.id
}

func (p *Publisher) Name() string {
	return p.name
}

func NewPublisher(id int, name string) *Publisher {
	return &Publisher{
		id,
		name,
	}
}

func GetPublisher(id *int, name *string) (*Publisher, error) {
	var result Publisher
	w := crud.Where.New(crud.Where{})
	if id == nil && name == nil {
		return nil, fmt.Errorf("empty id and name")
	}

	if id != nil {
		w.Equals("id", *id)
	}
	if name != nil {
		w.Equals("name", *name)
	}
	if len(w.Statements) > 1 {
		w.Union = crud.U_And
	}

	res, err := crud.GetOne(publisherTable, &crud.Fields{"id", "name"}, &w, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if id, ok := res["id"].Get().(int64); ok {
		result.id = int(id)
	} else {
		return nil, fmt.Errorf("fail to assert type %v", res["id"].Get())
	}
	if name, ok := res["name"].Get().(string); ok {
		result.name = name
	} else {
		return nil, fmt.Errorf("fail to assert type %v", res["name"].Get())
	}
	return &result, nil
}

func DeletePublisher(id *int, name *string) error {
	w := crud.Where.New(crud.Where{})
	if id == nil && name == nil {
		return fmt.Errorf("empty id and name")
	}

	if id != nil {
		w.Statements["id"] = crud.Statement{Value: *id, Comparator: crud.Equals}
	}
	if name != nil {
		w.Statements["name"] = crud.Statement{Value: *name, Comparator: crud.Equals}
	}
	if len(w.Statements) > 0 {
		w.Union = crud.U_And
	}
	affected, err := crud.Delete(publisherTable, &w)
	if err != nil {
		return err
	}
	if affected == 0 {
		var idStr string
		var nameStr string
		if id == nil {
			idStr = "nil"
		} else {
			idStr = fmt.Sprintf("%d", *id)
		}
		if name == nil {
			nameStr = "nil"
		} else {
			nameStr = *name
		}
		return fmt.Errorf("unregistered publisher with id: %s and name: %s", idStr, nameStr)
	}
	return nil
}

func UpdatePublisher(id int, newName string) error {
	w := crud.Where.New(crud.Where{})
	w.Statements["id"] = crud.Statement{Value: id, Comparator: crud.Equals}
	var f = crud.Fields{
		fmt.Sprintf("name = '%s'", newName),
	}
	affected, err := crud.Update(publisherTable, &f, &w)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("unregistered publisher with id: %d and name: %s", id, newName)
	}
	return err
}

func AddPublisher(name string) (int, error) {
	if res, err := GetPublisher(nil, &name); err != nil {
		return -1, err
	} else if res != nil {
		return res.id, fmt.Errorf("already exists")
	}
	return crud.Insert(publisherTable, &crud.Fields{"name"}, &crud.Values{name})
}
