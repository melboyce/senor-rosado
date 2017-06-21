package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetJSON ...
func GetJSON(u string, t interface{}) (err error) {
	log.Printf("-i- HTTP .GET: %s", u)
	r, err := httpClient.Get(u)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(t); err != nil {
		return
	}

	if r.StatusCode != 200 {
		err = fmt.Errorf("!!! HTTP S%d: %s", r.StatusCode, u)
	}
	return
}
