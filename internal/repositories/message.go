package repos

import "fmt"

const messageTable = "message"

func DeleteMessage(id int) error {
	var w where = where{
		[]string{},
		"",
	}
	predicate := fmt.Sprintf("id = %d", id)
	w.fields = append(w.fields, predicate)
	affected, err := delete(messageTable, &w)
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

	var w = where{
		[]string{fmt.Sprintf("id = %d", id)},
		"",
	}
	var f = fields{
		fmt.Sprintf("status_id = %d", NewState.id),
	}
	affected, err := update(messageTable, &f, &w)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no message with id: %d", id)
	}
	return err
}

// set args publisherId, content
func AddMessage(content string, publisherId int) (int, error) {
	return insert(messageTable, &fields{"content", "publisher_id"}, &values{content, publisherId})
}

func AddMessages(publisherId int, msgs *[]string) (int, error) {
	var v []values
	for _, msg := range *msgs {
		v = append(v, values{publisherId, msg})
	}

	return insertMany(messageTable, &fields{"content", "publisher_id"}, &v)
}
