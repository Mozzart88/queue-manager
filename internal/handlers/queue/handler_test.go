//go:build integration

package queue_test

import (
	"bytes"
	"encoding/json"
	"expat-news/queue-manager/internal/handlers/queue"
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
	const target = "/queue"
	db_test_utils.SetupDB(t)
	onLogging := test_utils.SuppressLogging()
	defer onLogging()

	tests := []struct {
		req        *http.Request
		expected   *expected
		whantError *whantError
	}{
		{
			httptest.NewRequest(http.MethodGet, target+`?publisher=pagina12&state=new`, nil),
			&expected{
				http.StatusOK,
				httpServer.Response{
					Msg:  `[{"id":3,"publisher":"pagina12","msg":"some post from Pagina 12","state":"new"},{"id":4,"publisher":"pagina12","msg":"some other post from Pagina 12","state":"new"}]`,
					Code: http.StatusOK,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodGet, target+`?publisher=pagina12&state=active`, nil),
			&expected{
				http.StatusNoContent,
				httpServer.Response{
					Msg:  "[]",
					Code: http.StatusNoContent,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodGet, target+"?publisher=some", nil),
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
			httptest.NewRequest(http.MethodGet, target+"?state=new", nil),
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
			httptest.NewRequest(http.MethodGet, target, nil),
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
					Msg:  `Bad Request: mising required fields: publisher and/or msgs`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{"publisher":"perfil"}`)),
			&expected{
				http.StatusBadRequest,
				httpServer.Response{
					Msg:  `Bad Request: mising required fields: publisher and/or msgs`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{"msgs":["some","invalid","data"]}`)),
			&expected{
				http.StatusBadRequest,
				httpServer.Response{
					Msg:  `Bad Request: mising required fields: publisher and/or msgs`,
					Code: http.StatusBadRequest,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{"publisher":"perfil","msgs":["some new message"]}`)),
			&expected{
				http.StatusCreated,
				httpServer.Response{
					Msg:  `1`,
					Code: http.StatusCreated,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(http.MethodPost, target, bytes.NewBufferString(`{"publisher":"perfil","msgs":["some other new message","and the third new message"]}`)),
			&expected{
				http.StatusCreated,
				httpServer.Response{
					Msg:  `2`,
					Code: http.StatusCreated,
				},
			},
			nil,
		},
	}

	for i, test := range tests {
		test.req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		queue.Handler(w, test.req)
		res := w.Result()
		if res.StatusCode != test.expected.status {
			body, _ := io.ReadAll(res.Body)
			t.Errorf("test %d: status codes mismatch - expected %d, got %d\nbody:\n%s", i, test.expected.status, res.StatusCode, string(body))
			continue
		}
		var actual httpServer.Response
		if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
			t.Errorf("test %d: unexpected error occured while decoding json: %v", i, err)
			continue
		}
		if !test_handler_body_eq(actual, test.expected.body) {
			t.Errorf("test %d: body mismatch - expected\n%v\ngot\n%v", i, test.expected.body, actual)
			continue
		}
	}
}
