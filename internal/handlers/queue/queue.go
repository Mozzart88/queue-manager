package queue

import (
	"encoding/json"
	"errors"
	"expat-news/queue-manager/internal/db"
	repos "expat-news/queue-manager/internal/repositories"
	"expat-news/queue-manager/internal/utils"
	"expat-news/queue-manager/pkg/utils/httpServer"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func insert(q *db.Queue, msgs *[]string) httpServer.Response {
	added, err := q.AddMessages(msgs)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unregistered publisher") {
			return httpServer.NotFound(err.Error())
		}
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.Created(fmt.Sprintf("%d", added))
}

func get(q *db.Queue) httpServer.Response {
	msgs, err := q.GetMessages(nil)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	if len(msgs) == 0 {
		msgs = []db.Message{}
	}
	result, err := json.Marshal(msgs)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK(string(result))
}

func parseBody(queue *db.Queue, msgs *[]string, body io.ReadCloser) error {
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		return errors.New("mising required fields: publisher and/or msgs")
	} else {
		var data struct {
			Msgs      []string `json:"msgs"`
			Publisher string   `json:"publisher"`
		}
		if err := json.Unmarshal(bytes, &data); err != nil {
			return err
		}
		if data.Publisher == "" || len(data.Msgs) == 0 {
			return errors.New("mising required fields: publisher and/or msgs")
		}
		queue.Publisher = data.Publisher
		queue.State = repos.STATE_NEW
		*msgs = data.Msgs
	}
	return nil
}

func parseQuery(queue *db.Queue, query url.Values) error {
	state := query.Get("state")
	publisher := query.Get("publisher")
	if state == "" || publisher == "" {
		return errors.New("state and publisher parameter are mandatory")
	}
	queue.State = repos.State_t(state)
	queue.Publisher = publisher
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var queue db.Queue
	var response httpServer.Response

	if r.Method == http.MethodGet {
		if err := parseQuery(&queue, r.URL.Query()); err != nil {
			utils.SendError(w, r, httpServer.BadRequest(err.Error()))
			return
		}
		response = get(&queue)
	} else if r.Method == http.MethodPost {
		var msgs []string
		if err := parseBody(&queue, &msgs, r.Body); err != nil {
			utils.SendError(w, r, httpServer.BadRequest(err.Error()))
			return
		}
		response = insert(&queue, &msgs)
	} else {
		utils.SendError(w, r, httpServer.MethodNotAllowed(r.Method))
		return
	}
	if response.Code >= 400 {
		utils.SendError(w, r, response)
	} else {
		utils.Send(w, r, response)
	}
}
