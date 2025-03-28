package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"fmt"
)

type State struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func (s *State) Get() error {
	state, err := repos.GetState(s.Id, (*repos.State_t)(s.Name))
	if err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("unkown state: %s", *s.Name)
	}
	*s.Id = state.ID()
	*s.Name = string(state.Name())
	return nil
}
