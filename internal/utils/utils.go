package utils

import (
	"expat-news/queue-manager/pkg/logger"
	"expat-news/queue-manager/pkg/utils/httpServer"
	"fmt"
	"net/http"
)

func Send(w http.ResponseWriter, r *http.Request, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Message(fmt.Sprintf("%s %s %d", r.Method, r.URL, response.Code))
}

func SendError(w http.ResponseWriter, r *http.Request, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Error(fmt.Sprintf("%s %s %d %s", r.Method, r.URL, response.Code, response.Msg))
}
