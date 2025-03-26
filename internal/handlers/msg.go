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
)

func send(w http.ResponseWriter, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Message(fmt.Sprintf("%d %s", response.Code, response.Msg))
}

func sendError(w http.ResponseWriter, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Error(fmt.Sprintf("%d %s", response.Code, response.Msg))
}

func insert(msg *db.Message) httpServer.Response {
	if msg.Publisher == nil || msg.Msg == nil {
		return httpServer.BadRequest("missing requiered fields - publisher and/or msg")
	}
	err := msg.Insert()
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	result, err := json.Marshal(*msg)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.Created(string(result))
}

func updateState(msg *db.Message) httpServer.Response {
	if msg.ID == nil || msg.State == nil {
		return httpServer.BadRequest("missing requiered fields - id and/or state")
	}
	err := msg.SetState()
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK("ok")
}

func delete(msg *db.Message) httpServer.Response {
	if msg.ID == nil {
		return httpServer.BadRequest("missing requiered fields - id")
	}
	err := msg.Delete()
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK("ok")
}

func get(msg *db.Message) httpServer.Response {
	err := msg.Get()
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	if msg.ID != nil && msg.Msg == nil {
		return httpServer.NotFuond(fmt.Sprintf("message with id %d in queue", *msg.ID))
	}
	result, err := json.Marshal(*msg)
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

func Handler(w http.ResponseWriter, r *http.Request) {
	var message db.Message
	var response httpServer.Response
	if err := parseRequest(&message, r.Body, r.URL.Query()); err != nil {
		sendError(w, httpServer.BadRequest(err.Error()))
		return
	}
	if r.Method == http.MethodGet {
		response = get(&message)
	} else if r.Method == http.MethodPost {
		response = insert(&message)
	} else if r.Method == http.MethodPatch {
		response = updateState(&message)
	} else if r.Method == http.MethodDelete {
		response = delete(&message)
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
