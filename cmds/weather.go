package cmds

import "fmt"
import "os"

import "encoding/json"
import "io/ioutil"
import "net/http"

import "github.com/jmespath/go-jmespath"
import "github.com/weirdtales/senor-rosado/slack"


// Weather reports the current forecast for an airport code
func Weather(m slack.Message, r *slack.Reply) error {
    token := os.Getenv("DARKSKY_API")
    if token == "" {
        r.Text = "I got no DARKSKY_API token in my env. Sort it."
        return nil
    }

    if m.Subcommand == "" {
        r.Text = "Specify a location, fool, then I can help."
        return nil
    }

    lat, lng, err := getLongLat(m.Subcommand)
    if err != nil {
        return err
    }
    fmt.Printf("lat=%f lng=%f\n\n", lat, lng)
    return nil
}

func getLongLat(q string) (lat float64, lng float64, err error) {
    url := fmt.Sprintf("http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false", q)
    r, err := http.Get(url)
    if err != nil {
        return
    }
    if r.StatusCode != 200 {
        err = fmt.Errorf("Error: status code: %d", r.StatusCode)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        return
    }

    var data interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        return
    }

    v, err := jmespath.Search("results[0].geometry.location", data)
    if err != nil {
        return
    }
    vals := v.(map[string]interface{})

    return vals["lat"].(float64), vals["lng"].(float64), nil
}
