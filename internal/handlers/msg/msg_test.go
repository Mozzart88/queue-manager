package msg

import (
	"expat-news/queue-manager/internal/db"
	"io"
	"net/url"
	"strings"
	"testing"
)

func TestParseRequest_invalidJson(t *testing.T) {
	var message db.Message
	body := io.NopCloser(strings.NewReader(`invalid json`))
	values := url.Values{}
	err := parseRequest(&message, body, values)
	if err == nil {
		t.Errorf("expected JSON decoding error, got nil")
	}
}

func TestParseRequest_nil_body_and_query(t *testing.T) {
	var message db.Message
	body := io.ReadCloser(nil)
	values := url.Values{}
	err := parseRequest(&message, body, values)
	if err == nil || err.Error() != "id parameter is mandatory" {
		t.Errorf("expected empty request error, got %v", err)
	}
}

func TestParseRequest_query(t *testing.T) {
	var message db.Message
	body := io.ReadCloser(nil)
	values := url.Values{}
	values.Add("id", "100")
	err := parseRequest(&message, body, values)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	const expected = 100
	if *message.ID != expected {
		t.Errorf("expected: %v, got: %v", expected, *message.ID)
	}
}

func TestParseRequest_body(t *testing.T) {
	var message db.Message
	body := io.NopCloser(strings.NewReader("{\"id\": 100}"))
	values := url.Values{}
	err := parseRequest(&message, body, values)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	const expected = 100
	if *message.ID != expected {
		t.Errorf("expected: %v, got: %v", expected, *message.ID)
	}
}

func TestParseRequest_body_and_values(t *testing.T) {
	var message db.Message
	body := io.NopCloser(strings.NewReader("{\"id\": 100}"))
	values := url.Values{}
	values.Add("id", "101")
	err := parseRequest(&message, body, values)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	const expected = 100
	if *message.ID != expected {
		t.Errorf("expected: %v, got: %v", expected, *message.ID)
	}
}
