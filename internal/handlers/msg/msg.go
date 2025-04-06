package msg

import (
	"encoding/json"
	"expat-news/queue-manager/internal/db"
	"expat-news/queue-manager/internal/utils"
	httpServer "expat-news/queue-manager/pkg/utils/httpServer"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func insert(msg *db.Message) httpServer.Response {
	if msg.Publisher == nil || msg.Msg == nil {
		return httpServer.BadRequest("missing requiered fields - publisher and/or msg")
	}
	err := msg.Add()
	if err != nil {
		if strings.HasPrefix(err.Error(), "unregistered publisher") {
			return httpServer.BadRequest(err.Error())
		}
		return httpServer.InternalServerError(err.Error())
	}
	result, err := json.Marshal(*msg)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.Created(string(result))
}

func updateState(msg *db.Message) httpServer.Response {
	if msg.Id == nil || msg.State == nil {
		return httpServer.BadRequest("missing requiered fields - id and/or state")
	}
	if err := msg.SetState(); err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK("ok")
}

func delete(msg *db.Message) httpServer.Response {
	if msg.Id == nil {
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
		if strings.HasPrefix(err.Error(), "no message in queue with specified id") {
			return httpServer.NotFound(err.Error())
		}
		return httpServer.InternalServerError(err.Error())
	}
	if msg.Id != nil && msg.Msg == nil {
		return httpServer.NotFound(fmt.Sprintf("message with id %d in queue", *msg.Id))
	}
	result, err := json.Marshal(*msg)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK(string(result))
}

func parseRequest(msg *db.Message, body io.ReadCloser, query url.Values) error {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		id := query.Get("id")
		if id == "" {
			return fmt.Errorf("id parameter is mandatory")
		}
		value, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		msg.Id = &value
	} else {
		if err := json.Unmarshal(data, msg); err != nil {
			return err
		}
	}
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var message db.Message
	var response httpServer.Response
	if err := parseRequest(&message, r.Body, r.URL.Query()); err != nil {
		utils.SendError(w, httpServer.BadRequest(err.Error()))
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
		utils.SendError(w, httpServer.MethodNotAllowed(r.Method))
		return
	}
	if response.Code >= 400 {
		utils.SendError(w, response)
	} else {
		utils.Send(w, response)
	}
}
