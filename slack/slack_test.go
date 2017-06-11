package slack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeJSON struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// TestConnect ...
func TestConnect(t *testing.T) {
	token := "test"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"ok": true,
			"url": "wss://ms9.slack-msgs.com/websocket/2I5yBpcvk",
			"team": {
				"id": "T654321",
				"name": "Librarian Society of Soledad",
				"domain": "libsocos",
				"enterprise_id": "E234567",
				"enterprise_name": "Intercontinental Librarian Society"
			},
			"self": {
				"id": "W123456",
				"name": "brautigan"
			}
		}`
		fmt.Fprintln(w, resp)
	}))
	defer ts.Close()
	slackAPIURL = ts.URL + "?token=%s"
	conn, err := Connect(token)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "websocket.Dial") {
			t.Errorf("ERR %s", err)
		}
	}
	if !conn.Ok {
		t.Errorf("ERR conn.Ok != true")
	}
}

// TestGetJSON ...
func TestGetJSON(t *testing.T) {
	url := "https://jsonplaceholder.typicode.com/posts/1"
	target := fakeJSON{}
	if err := GetJSON(url, &target); err != nil {
		t.Errorf("ERR %s", err)
	}
	if target.ID != 1 {
		t.Errorf("ERR target not populated correctly")
	}
	if target.UserID != 1 {
		t.Errorf("ERR target not populated correctly")
	}
	if target.Title != "sunt aut facere repellat provident occaecati excepturi optio reprehenderit" {
		t.Errorf("ERR target not populated correctly")
	}
}
