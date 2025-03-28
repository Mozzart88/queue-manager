package repos

import "fmt"

const messageTable = "message"

func DeleteMessage(id int) error {
	var l limit = 1
	var w where = where{
		[]string{},
		"",
	}
	predicate := fmt.Sprintf("id = %d", id)
	w.fields = append(w.fields, predicate)
	_, err := delete(messageTable, &w, &l)
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

	var w = where{
		[]string{fmt.Sprintf("id = %d", id)},
		"",
	}
	var f = fields{
		fmt.Sprintf("status_id = %d", NewState.id),
	}
	_, err = update(messageTable, &f, &w)
	return err
}

func AddMessage(content string, publisherId int) (int, error) {
	return insert(messageTable, &fields{"content", "publisher_id"}, &values{content, publisherId})
}
