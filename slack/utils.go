package slack

import (
	"encoding/json"
	"fmt"
	"log"
)

// GetJSON unmarshals an API call into a struct
func GetJSON(url string, t interface{}) (err error) {
	log.Printf("-i- HTTP .GET: %s", url)
	r, err := httpClient.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		return
	}

	if r.StatusCode != 200 {
		err = fmt.Errorf("ERR HTTP S%d: %s", r.StatusCode, url)
	}
	return
}
