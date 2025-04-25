package repos

import (
	"expat-news/queue-manager/internal/repositories/crud"
	"fmt"
)

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
	if id == nil && name == nil {
		return nil, fmt.Errorf("empty id and name")
	}
	w := crud.NewWhere()

	if id != nil {
		w.Equals("id", *id)
	}
	if name != nil {
		w.Equals("name", *name)
	}
	if w.Len() > 0 {
		w.Union = crud.U_And
	}

	res, err := crud.GetOne(stateTable, &crud.Fields{"id", "name"}, w, nil)
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
		result.name = State_t(name)
	} else {
		return nil, fmt.Errorf("fail to assert type %v", res["name"].Get())
	}
	return &result, nil
}
