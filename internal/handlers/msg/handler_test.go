//go:build integration

package msg_test

import (
	"bytes"
	"encoding/json"
	"expat-news/queue-manager/internal/handlers/msg"
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

func test_msg_handler_body_eq(a, b httpServer.Response) bool {
	if a.Code != b.Code {
		return false
	}
	if a.Msg != b.Msg && !strings.HasPrefix(a.Msg, b.Msg) {
		return false
	}
	return true
}

func TestHandler(t *testing.T) {
	const target = "/msg"
	db_test_utils.SetupDB(t)
	releaseSuppress := test_utils.SuppressLogging()
	defer releaseSuppress()
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
					Msg:  `{"id":1,"publisher":"pagina12","msg":"some post from Pagina 12 that already Done","state":"done"}`,
					Code: http.StatusOK,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPost,
				target,
				bytes.NewBufferString(`{"publisher":"pagina12","msg":"some new msg"}`),
			),
			&expected{
				http.StatusCreated,
				httpServer.Response{
					Msg:  `{"id":7,"publisher":"pagina12","msg":"some new msg","state":"new"}`,
					Code: http.StatusCreated,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPost,
				target,
				bytes.NewBufferString(`{"msg": "{\"content\":\"**Latest news for today**\\n\\nA tragic accident occurred in Santiago del Estero, where a 2-year-old boy died after choking on a candy while at his grandmother's house with his parents. Despite medical efforts to save him, the child did not survive.\\n\\nIn other news, Argentine President Mauricio Macri visited Mar del Plata, where he criticized negotiations between the PRO and La Libertad Avanza parties.\\n\\nAdditionally, former President Roberto García Moritán was questioned on a TV program about his alleged infidelity to Pampita, leaving him emotional.\\n\\nThe US-China trade war continues to affect global markets, with the International Monetary Fund warning that the global debt public could surpass pandemic levels due to tariffs.\",\"recipient\":\"post\"}","publisher": "pagina12"}`),
			),
			&expected{
				http.StatusCreated,
				httpServer.Response{
					Msg:  `{"id":8,"publisher":"pagina12","msg":"{\"content\":\"**Latest news for today**\\n\\nA tragic accident occurred in Santiago del Estero, where a 2-year-old boy died after choking on a candy while at his grandmother's house with his parents. Despite medical efforts to save him, the child did not survive.\\n\\nIn other news, Argentine President Mauricio Macri visited Mar del Plata, where he criticized negotiations between the PRO and La Libertad Avanza parties.\\n\\nAdditionally, former President Roberto García Moritán was questioned on a TV program about his alleged infidelity to Pampita, leaving him emotional.\\n\\nThe US-China trade war continues to affect global markets, with the International Monetary Fund warning that the global debt public could surpass pandemic levels due to tariffs.\",\"recipient\":\"post\"}","state":"new"}`,
					Code: http.StatusCreated,
				},
			},
			nil,
		},
		{
			httptest.NewRequest(
				http.MethodPatch,
				target,
				bytes.NewBufferString(`{"id":1,"state":"new"}`),
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
				bytes.NewBufferString(`{"id":1}`),
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
		msg.Handler(w, test.req)
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
		if !test_msg_handler_body_eq(actual, test.expected.body) {
			t.Errorf("test %d: body mismatch - expected\n%v\ngot\n%v", i, test.expected.body, actual)
			continue
		}
	}
}

func TestHandler_negative(t *testing.T) {
	const target = "/msg"
	db_test_utils.SetupDB(t)
	releaseSuppress := test_utils.SuppressLogging()
	defer releaseSuppress()
	tests := []struct {
		req        *http.Request
		expected   *expected
		whantError *whantError
	}{
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
			httptest.NewRequest(http.MethodGet, target+"?id=256", nil),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found:`,
					Code: http.StatusNotFound,
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
			httptest.NewRequest(
				http.MethodPost,
				target,
				bytes.NewBufferString(`{"publisher":"pagina1","msg":"some new msg"}`),
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
				http.MethodPatch,
				target,
				bytes.NewBufferString(`{"id":1,"state":"ne"}`),
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
				http.MethodPatch,
				target,
				bytes.NewBufferString(`{"id":256,"state":"new"}`),
			),
			&expected{
				http.StatusNotFound,
				httpServer.Response{
					Msg:  `Not Found:`,
					Code: http.StatusNotFound,
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
					Msg:  `Not Found:`,
					Code: http.StatusNotFound,
				},
			},
			nil,
		},
	}

	for i, test := range tests {
		test.req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		msg.Handler(w, test.req)
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
		if !test_msg_handler_body_eq(actual, test.expected.body) {
			t.Errorf("test %d: body mismatch - expected\n%v\ngot\n%v", i, test.expected.body, actual)
			continue
		}
	}
}
