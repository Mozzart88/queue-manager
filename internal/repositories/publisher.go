package repos

import "fmt"

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
	var w where = where{
		[]string{},
		"",
	}
	if id == nil && name == nil {
		return nil, fmt.Errorf("empty id and name")
	}

	if id != nil {
		predicate := fmt.Sprintf("id = %d", *id)
		w.fields = append(w.fields, predicate)
	}
	if name != nil {
		predicate := fmt.Sprintf("name = '%s'", *name)
		w.fields = append(w.fields, predicate)
	}
	if len(w.fields) > 1 {
		w.union = "and"
	}

	res, err := getOne(publisherTable, &fields{"id", "name"}, &w, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if id, ok := res["id"].value.(int64); ok {
		result.id = int(id)
	} else {
		return nil, fmt.Errorf("fail to assert type %v", res["id"].value)
	}
	if name, ok := res["name"].value.(string); ok {
		result.name = name
	} else {
		return nil, fmt.Errorf("fail to assert type %v", res["name"].value)
	}
	return &result, nil
}

func DeletePublisher(id *int, name *string) error {
	var w where = where{
		[]string{},
		"",
	}
	if id == nil && name == nil {
		return fmt.Errorf("empty id and name")
	}

	if id != nil {
		predicate := fmt.Sprintf("id = %d", *id)
		w.fields = append(w.fields, predicate)
	}
	if name != nil {
		predicate := fmt.Sprintf("name = '%s'", *name)
		w.fields = append(w.fields, predicate)
	}
	if len(w.fields) > 0 {
		w.union = "and"
	}
	affected, err := delete(publisherTable, &w)
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
	var w = where{
		[]string{fmt.Sprintf("id = %d", id)},
		"",
	}
	var f = fields{
		fmt.Sprintf("name = '%s'", newName),
	}
	affected, err := update(publisherTable, &f, &w)
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
	return insert(publisherTable, &fields{"name"}, &values{name})
}
