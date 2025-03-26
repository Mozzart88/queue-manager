package msg

import (
	"encoding/json"
	"errors"
	"expat-news/queue-manager/internal/db"
	"expat-news/queue-manager/pkg/logger"
	httpServer "expat-news/queue-manager/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

var mu sync.Mutex

func send(w http.ResponseWriter, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Message(fmt.Sprintf("%d %s", response.Code, response.Msg))
}

func sendError(w http.ResponseWriter, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Error(fmt.Sprintf("%d %s", response.Code, response.Msg))
}

func insert(
	msg *db.Message,
	fn func(publisher string, msg string) (db.Message, error),
) httpServer.Response {
	if msg.Publisher == nil || msg.Msg == nil {
		return httpServer.BadRequest("missing requiered fields - publisher and/or msg")
	}
	data, err := fn(*msg.Publisher, *msg.Msg)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	result, err := json.Marshal(data)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.Created(string(result))
}

func updateState(
	msg *db.Message,
	fn func(id int, state string) error,
) httpServer.Response {
	if msg.ID == nil || msg.State == nil {
		return httpServer.BadRequest("missing requiered fields - id and/or state")
	}
	err := fn(*msg.ID, *msg.State)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK("ok")
}

func delete(
	msg *db.Message,
	fn func(id int) error,
) httpServer.Response {
	if msg.ID == nil {
		return httpServer.BadRequest("missing requiered fields - id")
	}
	err := fn(*msg.ID)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK("ok")
}

func get(
	msg *db.Message,
	fn func(id int) (*db.Message, error),
) httpServer.Response {
	if msg.ID == nil {
		return httpServer.BadRequest("missing requiered fields - id")
	}
	data, err := fn(*msg.ID)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	if data == nil {
		return httpServer.NotFuond(fmt.Sprintf("message with id %d in queue", *msg.ID))
	}
	result, err := json.Marshal(*data)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK(string(result))
}

func parseRequest(msg *db.Message, body io.ReadCloser, query url.Values) error {
	if body == nil {
		id := query.Get("id")
		if id == "" {
			return errors.New("id parameter is mandatory")
		}
		value, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		msg.ID = &value
	} else {
		defer body.Close()
		if err := json.NewDecoder(body).Decode(msg); err != nil {
			return err
		}
	}
	return nil
}

func Handler(
	insertHandler func(publisher string, msg string) (db.Message, error),
	updateStateHandler func(id int, state string) error,
	getHandler func(id int) (*db.Message, error),
	deleteHandler func(id int) error,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message db.Message
		var response httpServer.Response
		if err := parseRequest(&message, r.Body, r.URL.Query()); err != nil {
			sendError(w, httpServer.BadRequest(err.Error()))
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if r.Method == http.MethodGet {
			response = get(&message, getHandler)
		} else if r.Method == http.MethodPost {
			response = insert(&message, insertHandler)
		} else if r.Method == http.MethodPut {
			response = updateState(&message, updateStateHandler)
		} else if r.Method == http.MethodDelete {
			response = delete(&message, deleteHandler)
		} else {
			sendError(w, httpServer.MethodNotAllowed(r.Method))
			return
		}
		if response.Code >= 400 {
			sendError(w, response)
		} else {
			send(w, response)
		}
	}

}
