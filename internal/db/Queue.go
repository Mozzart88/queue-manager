package db

import (
	repos "expat-news/queue-manager/internal/repositories"
)

type Queue struct {
	Publisher string        `json:"publisher"`
	State     repos.State_t `json:"state"`
}

func (q *Queue) GetMessages(oldest *bool) ([]Message, error) {
	var result []Message

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
