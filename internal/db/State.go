package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"expat-news/queue-manager/pkg/utils"
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
		return fmt.Errorf("unkown state: %v", *s)
	}
	if s.Id == nil {
		s.Id = utils.Ptr(state.ID())
	} else {
		*s.Id = state.ID()
	}
	if s.Name == nil {
		s.Name = utils.Ptr(string(state.Name()))
	} else {
		*s.Name = string(state.Name())
	}
	return nil
}
