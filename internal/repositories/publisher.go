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
	if len(w.fields) > 0 {
		w.union = "and"
	}

	res, err := getOne(publisherTable, &fields{"id", "name"}, &w, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if id, ok := res["id"].value.(int); ok {
		result.id = id
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
	var l limit = 1
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
	_, err := delete(publisherTable, &w, &l)
	return err
}

func UpdatePublisher(id int, newName string) error {
	var w = where{
		[]string{fmt.Sprintf("id = %d", id)},
		"",
	}
	var f = fields{
		fmt.Sprintf("name = '%s'", newName),
	}
	_, err := update(publisherTable, &f, &w)
	return err
}

func AddPublisher(name string) (int, error) {
	return insert(publisherTable, &fields{"name"}, &values{name})
}
