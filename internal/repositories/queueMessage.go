package repos

import (
	"fmt"
)

const queueTable = "queue"

type QueueMessage struct {
	id            int
	content       string
	publisherId   int
	publisherName string
	stateId       int
	stateName     string
}

func (m *QueueMessage) ID() int {
	return m.id
}

func (m *QueueMessage) Content() string {
	return m.content
}

func (m *QueueMessage) PublisherId() int {
	return m.publisherId
}

func (m *QueueMessage) PublisherName() string {
	return m.publisherName
}

func (m *QueueMessage) StateId() int {
	return m.stateId
}

func (m *QueueMessage) StateName() string {
	return m.stateName
}

func (q *QueueMessage) setId(value any) error {
	if id, ok := value.(int64); ok {
		q.id = int(id)
	} else {
		return fmt.Errorf("fail to assert type %v", value)
	}
	return nil
}

func (q *QueueMessage) setContent(value any) error {
	if v, ok := value.(string); ok {
		q.content = v
	} else {
		return fmt.Errorf("fail to assert type %v", value)
	}
	return nil
}

func (q *QueueMessage) setPublisherId(value any) error {
	if v, ok := value.(int64); ok {
		q.publisherId = int(v)
	} else {
		return fmt.Errorf("fail to assert type %v", value)
	}
	return nil
}

func (q *QueueMessage) setPublisherName(value any) error {
	if v, ok := value.(string); ok {
		q.publisherName = v
	} else {
		return fmt.Errorf("fail to assert type %v", value)
	}
	return nil
}

func (q *QueueMessage) setStateId(value any) error {
	if v, ok := value.(int64); ok {
		q.stateId = int(v)
	} else {
		return fmt.Errorf("fail to assert type %v", value)
	}
	return nil
}

func (q *QueueMessage) setStateName(value any) error {
	if v, ok := value.(string); ok {
		q.stateName = v
	} else {
		return fmt.Errorf("fail to assert type %v", value)
	}
	return nil
}

func (q *QueueMessage) fillValues(row QueryRow) error {
	if err := q.setId(row["id"].value); err != nil {
		return err
	}
	if err := q.setContent(row["msg"].value); err != nil {
		return err
	}
	if err := q.setPublisherId(row["publisher_id"].value); err != nil {
		return err
	}
	if err := q.setStateId(row["status_id"].value); err != nil {
		return err
	}
	if err := q.setPublisherName(row["publisher"].value); err != nil {
		return err
	}
	if err := q.setStateName(row["status"].value); err != nil {
		return err
	}
	return nil
}

type oldest_t *bool

func queueFields() fields {
	return fields{
		"id",
		"msg",
		"publisher_id",
		"publisher",
		"status_id",
		"status",
	}
}

func NewQueueMessage(id int, content string, publisherId int, publisherName string, stateId int, stateName string) *QueueMessage {
	return &QueueMessage{
		id,
		content,
		publisherId,
		publisherName,
		stateId,
		stateName,
	}
}

func GetQueueMessage(publisher string, s State_t, o oldest_t) (*QueueMessage, error) {
	var result *QueueMessage = &QueueMessage{}
	var w where = where{[]string{}, "AND"}
	var f fields = queueFields()
	var order = order{[]string{"id"}, "ASC"}

	if o != nil && !*o {
		order.order = "DESC"
	}

	w.fields = append(w.fields, fmt.Sprintf("publisher = '%s'", publisher))
	w.fields = append(w.fields, fmt.Sprintf("status = '%s'", s))

	res, err := getOne(queueTable, &f, &w, &order)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if err := result.fillValues(res); err != nil {
		return nil, err
	}
	return result, nil
}

func GetUniqQueueMessage(id int) (*QueueMessage, error) {
	var result = &QueueMessage{}
	var w where = where{[]string{}, ""}
	var f fields = queueFields()

	w.fields = append(w.fields, fmt.Sprintf("id = %d", id))

	res, err := getOne(queueTable, &f, &w, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if err := result.fillValues(res); err != nil {
		return nil, err
	}
	return result, nil
}

func GetMessages(publisher string, s State_t, o oldest_t) ([]QueueMessage, error) {
	var result []QueueMessage
	var w where = where{[]string{}, "AND"}
	var f fields = queueFields()
	var order = order{[]string{"id"}, "ASC"}

	if o != nil && !*o {
		order.fields = []string{"id"}
		order.order = "DESC"
	}

	w.fields = append(w.fields, fmt.Sprintf("publisher = '%s'", publisher))
	w.fields = append(w.fields, fmt.Sprintf("status = '%s'", s))

	res, err := get(queueTable, &f, &w, &order, nil)
	if err != nil {
		return nil, err
	}
	for _, msg := range res {
		var m = &QueueMessage{}
		if err := m.fillValues(msg); err != nil {
			return nil, err
		}
		result = append(result, *m)
	}
	return result, nil
}

func AddQueueMessage(content string, publisherName string) (*QueueMessage, error) {
	var publisherId int
	publisher, err := GetPublisher(nil, &publisherName)
	if err != nil {
		return nil, err
	}
	if publisher == nil {
		return nil, fmt.Errorf("unregistered publisher: %s", publisherName)
	}
	publisherId = publisher.ID()
	id, err := AddMessage(content, publisherId)
	if err != nil {
		return nil, err
	}
	return GetUniqQueueMessage(id)
}
