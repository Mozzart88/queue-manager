package httpServer

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func BadRequest(msg string) Response {
	return Response{
		Msg:  "Bad Request: " + msg,
		Code: http.StatusBadRequest,
	}
}

func InternalServerError(msg string) Response {
	return Response{
		Msg:  "Internal Server Error: " + msg,
		Code: http.StatusInternalServerError,
	}
}

func NotFound(msg string) Response {
	return Response{
		Msg:  "Not Found: " + msg,
		Code: http.StatusNotFound,
	}
}

func Found(msg string) Response {
	return Response{
		Msg:  "Found: " + msg,
		Code: http.StatusFound,
	}
}

func NoContent(msg string) Response {
	return Response{
		Msg:  "Found: " + msg,
		Code: http.StatusNoContent,
	}
}

func MethodNotAllowed(method string) Response {
	return Response{
		Msg:  "Method not Allowed: " + method,
		Code: http.StatusMethodNotAllowed,
	}
}

func Created(msg string) Response {
	return Response{
		Msg:  msg,
		Code: http.StatusCreated,
	}
}

func OK(msg string) Response {
	return Response{
		Msg:  msg,
		Code: http.StatusOK,
	}
}

func SendResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}
