package utils

import (
	"expat-news/queue-manager/pkg/logger"
	"expat-news/queue-manager/pkg/utils/httpServer"
	"fmt"
	"net/http"
)

func Send(w http.ResponseWriter, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Message(fmt.Sprintf("%d %s", response.Code, response.Msg))
}

func SendError(w http.ResponseWriter, response httpServer.Response) {
	httpServer.SendResponse(w, response)
	logger.Error(fmt.Sprintf("%d %s", response.Code, response.Msg))
}
