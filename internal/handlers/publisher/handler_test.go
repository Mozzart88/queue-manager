//go:build integration

package publisher_test

import (
	"bytes"
	"encoding/json"
	"expat-news/queue-manager/internal/handlers/publisher"
	"expat-news/queue-manager/internal/repositories/db_test_utils"
	"expat-news/queue-manager/internal/test_utils"
	"expat-news/queue-manager/pkg/utils/httpServer"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type whantError struct {
}

type expected struct {
	status int
	body   httpServer.Response
}

func test_handler_body_eq(a, b httpServer.Response) bool {
	if a.Code != b.Code {
		return false
	}
	if a.Msg != b.Msg && !strings.HasPrefix(a.Msg, b.Msg) {
		return false
	}
	return true
}

func TestHandler(t *testing.T) {
	const target = "/publisher"
	db_test_utils.SetupDB(t)
	onLogging := test_utils.SuppressLogging()
	defer onLogging()

	tests := []struct {
		req        *http.Request
		expected   *expected
		whantError *whantError
	}{
		{
			httptest.NewRequest(http.MethodGet, target+"?id=1", nil),
			&expected{
				http.StatusOK,
				httpServer.Response{
					Msg:  `{"id":1,"name":"pagina12"}`,
					Code: http.StatusOK,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPost,
				target,
				bytes.NewBufferString(`{"name":"some"}`),
			),
			&expected{
				http.StatusCreated,
				httpServer.Response{
					Msg:  `{"id":4,"name":"some"}`,
					Code: http.StatusCreated,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPatch,
				target,
				bytes.NewBufferString(`{"id":4,"name":"other"}`),
			),
			&expected{
				http.StatusOK,
				httpServer.Response{
					Msg:  `ok`,
					Code: http.StatusOK,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodDelete,
				target,
				bytes.NewBufferString(`{"id":4}`),
			),
			&expected{
				http.StatusOK,
				httpServer.Response{
					Msg:  `ok`,
					Code: http.StatusOK,
				},
			},
			nil,
		},
	}

	for i, test := range tests {
		test.req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		publisher.Handler(w, test.req)
		res := w.Result()
		if res.StatusCode != test.expected.status {
			body, _ := io.ReadAll(res.Body)
			test_utils.Fail(t, i, "status codes mismatch - expected %d, got %d\nbody:\n%s", test.expected.status, res.StatusCode, string(body))
			continue
		}
		var actual httpServer.Response
		if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
			test_utils.Fail(t, i, "unexpected error occured while decoding json: %v", err)
			continue
		}
		if !test_handler_body_eq(actual, test.expected.body) {
			test_utils.Fail(t, i, "body mismatch - expected\n%v\ngot\n%v", test.expected.body, actual)
			continue
		}
	}
}

func TestHandler_negative(t *testing.T) {
	const target = "/publisher"
	db_test_utils.SetupDB(t)
	onLogging := test_utils.SuppressLogging()
	defer onLogging()

	tests := []struct {
		req        *http.Request
		expected   *expected
		whantError *whantError
	}{
		{
			httptest.NewRequest(http.MethodGet, target+"?id=256", nil),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found: unregistered publisher`,
					Code: http.StatusNotFound,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodGet, target+"?name=some", nil),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found: unregistered publisher`,
					Code: http.StatusNotFound,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodGet, target+"?id=1&name=some", nil),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found: unregistered publisher`,
					Code: http.StatusNotFound,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodGet, target+"?msg=1", nil),
			&expected{
				http.StatusBadRequest,
				httpServer.Response{
					Msg:  `Bad Request:`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{}`)),
			&expected{
				http.StatusBadRequest,
				httpServer.Response{
					Msg:  `Bad Request:`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{"name":""}`)),
			&expected{
				http.StatusBadRequest,
				httpServer.Response{
					Msg:  `Bad Request:`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{"name":"pagina12"}`)),
			&expected{
				http.StatusFound,
				httpServer.Response{
					Msg:  `Found:`,
					Code: http.StatusFound,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPatch,
				target,
				bytes.NewBufferString(`{"id":256,"name":"some"}`),
			),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found: unregistered publisher`,
					Code: http.StatusNotFound,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPatch,
				target,
				bytes.NewBufferString(`{"id":1}`),
			),
			&expected{
				http.StatusBadRequest,
				httpServer.Response{
					Msg:  `Bad Request:`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodDelete,
				target,
				bytes.NewBufferString(`{"id":256}`),
			),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found: unregistered publisher with id: 256 and name: nil`,
					Code: http.StatusNotFound,
				},
			},
			nil,
		},
	}

	for i, test := range tests {
		test.req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		publisher.Handler(w, test.req)
		res := w.Result()
		if res.StatusCode != test.expected.status {
			body, _ := io.ReadAll(res.Body)
			test_utils.Fail(t, i, "status codes mismatch - expected %d, got %d\nbody:\n%s", test.expected.status, res.StatusCode, string(body))
			continue
		}
		var actual httpServer.Response
		if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
			test_utils.Fail(t, i, "unexpected error occured while decoding json: %v", err)
			continue
		}
		if !test_handler_body_eq(actual, test.expected.body) {
			test_utils.Fail(t, i, "body mismatch - expected\n%v\ngot\n%v", test.expected.body, actual)
			continue
		}
	}
}
