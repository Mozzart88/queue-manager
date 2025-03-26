package db

type Message struct {
	ID        *int    `json:"id,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
	Msg       *string `json:"msg,omitempty"`
	State     *string `json:"state,omitempty"`
}
