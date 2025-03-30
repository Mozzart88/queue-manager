package repos

import "fmt"

type State_t string

const (
	STATE_NEW    State_t = "new"
	STATE_ACTIVE State_t = "active"
	STATE_DONE   State_t = "done"
)

const stateTable = "message_status"

type State struct {
	id   int
	name State_t
}

func (s *State) ID() int {
	return s.id
}

func (s *State) Name() State_t {
	return s.name
}

func NewState(id int, name State_t) *State {
	return &State{id, name}
}

func GetState(id *int, name *State_t) (*State, error) {
	var result State
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

	res, err := getOne(stateTable, &fields{"id", "name"}, &w, nil)
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
		result.name = State_t(name)
	} else {
		return nil, fmt.Errorf("fail to assert type %v", res["name"].value)
	}
	return &result, nil
}
