package db

import (
	repos "expat-news/queue-manager/internal/repositories"
	"fmt"
)

type Queue struct {
	Publisher string        `json:"publisher"`
	State     repos.State_t `json:"state"`
}

func (q *Queue) GetMessages(oldest *bool) ([]Message, error) {
	var result []Message

	mu.Lock()
	defer mu.Unlock()
	res, err := repos.GetMessages(q.Publisher, q.State, oldest)
	if err != nil {
		return nil, err
	}
	for _, msg := range res {
		var m Message
		m.fillValues(&msg)
		result = append(result, m)
	}
	return result, nil
}

func (q *Queue) GetMessage(oldest *bool) (*Message, error) {
	var result Message

	mu.Lock()
	defer mu.Unlock()
	msg, err := repos.GetQueueMessage(q.Publisher, q.State, oldest)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, nil
	}
	result.fillValues(msg)
	return &result, nil
}

func (q *Queue) AddMessages(msgs *[]string) (int, error) {
	var publisherId int

	mu.Lock()
	defer mu.Unlock()
	publisher, err := repos.GetPublisher(nil, &q.Publisher)
	if err != nil {
		return 0, err
	}
	if publisher == nil {
		return 0, fmt.Errorf("unregistered publisher: %s", q.Publisher)
	}
	publisherId = publisher.ID()
	return repos.AddMessages(publisherId, msgs)
}
