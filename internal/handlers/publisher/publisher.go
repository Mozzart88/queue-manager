package publisher

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

func register(publisher *db.Publisher) httpServer.Response {
	if publisher.Name == nil {
		return httpServer.BadRequest("missing required field: name")
	}
	if err := publisher.Register(); err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	result, err := json.Marshal(publisher)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.Created(string(result))
}

func rename(publisher *db.Publisher) httpServer.Response {
	if publisher.Id == nil || publisher.Name == nil {
		return httpServer.BadRequest("missing requiered fields: id and/or name")
	}
	if err := publisher.Update(*publisher.Name); err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK("ok")
}

func delete(publisher *db.Publisher) httpServer.Response {
	if publisher.Id == nil {
		return httpServer.BadRequest("missing requiered field: id")
	}
	if err := publisher.Delete(); err != nil {
		if ok := strings.HasPrefix(err.Error(), "unregistered publisher"); ok {
			return httpServer.NotFound(err.Error())
		} else {
			return httpServer.InternalServerError(err.Error())
		}
	}
	return httpServer.OK("ok")
}

func get(publisher *db.Publisher) httpServer.Response {
	if publisher.Id == nil && publisher.Name == nil {
		return httpServer.BadRequest("missing required fields: id and name")
	}
	if err := publisher.Get(); err != nil {
		if ok := strings.HasPrefix(err.Error(), "unregistered publisher"); ok {
			return httpServer.NotFound(err.Error())
		}
		return httpServer.InternalServerError(err.Error())
	}
	if publisher.Id == nil || publisher.Name == nil {
		return httpServer.NotFound(fmt.Sprintf("publisher not found: %v", publisher))
	}
	result, err := json.Marshal(publisher)
	if err != nil {
		return httpServer.InternalServerError(err.Error())
	}
	return httpServer.OK(string(result))
}

func parseRequest(publisher *db.Publisher, body io.ReadCloser, query url.Values) error {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		id := query.Get("id")
		name := query.Get("name")
		if id == "" && name == "" {
			return fmt.Errorf("missing parameter id and name")
		}
		if id != "" {
			value, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			publisher.Id = &value
		}
		if name != "" {
			publisher.Name = &name
		}
	} else {
		if err := json.Unmarshal(data, publisher); err != nil {
			return err
		}
	}
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var publisher db.Publisher
	var response httpServer.Response
	if err := parseRequest(&publisher, r.Body, r.URL.Query()); err != nil {
		utils.SendError(w, httpServer.BadRequest(err.Error()))
		return
	}
	if r.Method == http.MethodGet {
		response = get(&publisher)
	} else if r.Method == http.MethodPost {
		response = register(&publisher)
	} else if r.Method == http.MethodPatch {
		response = rename(&publisher)
	} else if r.Method == http.MethodDelete {
		response = delete(&publisher)
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
