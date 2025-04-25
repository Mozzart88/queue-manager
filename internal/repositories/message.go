package repos

import (
	"expat-news/queue-manager/internal/repositories/crud"
	"fmt"
)

const messageTable = "message"

func DeleteMessage(id int) error {
	w := crud.Where.New(crud.Where{})

	w.Equals("id", id)
	affected, err := crud.Delete(messageTable, &w)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no message with id: %d", id)
	}

	return err
}

func UpdateStateMessage(id int, newState State_t) error {
	NewState, err := GetState(nil, &newState)
	if err != nil {
		return err
	}
	if NewState == nil {
		return fmt.Errorf("invalid state: %s", newState)
	}

	var w = crud.Where.New(crud.Where{})
	w.Equals("id", id)

	var f = crud.Fields{
		fmt.Sprintf("status_id = %d", NewState.id),
	}
	affected, err := crud.Update(messageTable, &f, &w)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no message with id: %d", id)
	}
	return err
}

func AddMessage(content string, publisherId int) (int, error) {
	return crud.Insert(messageTable, &crud.Fields{"content", "publisher_id"}, &crud.Values{content, publisherId})
}

func AddMessages(publisherId int, msgs *[]string) (int, error) {
	var v []crud.Values
	for _, msg := range *msgs {
		v = append(v, crud.Values{msg, publisherId})
	}

	return crud.InsertMany(messageTable, &crud.Fields{"content", "publisher_id"}, &v)
}
