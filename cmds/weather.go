// Package cmds provides commands to a bot.
package cmds

import (
    "fmt"
    "os"
    "time"

    "encoding/json"
    "net/http"

    "github.com/weirdtales/senor-rosado/slack"
)


// partial struct for a google map api hit - so gross.
type weatherLocation struct {
    Results []struct {
        FormattedAddress string `json:"formatted_address"`
        Geometry struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
        } `json:"geometry"`
        Status string `json:"status"`
    } `json:"results"`
}

// partial struct for a darksky api hit - so pretty.
type weatherInfo struct {
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Timezone string `json:"timezone"`
    Currently struct {
        Time int64 `json:"time"`
        Summary string `json:"summary"`
        Icon string `json:"icon"`
        Temperature float64 `json:"temperature"`
        ApparentTemperature float64 `json:"apparentTemperature"`
        DewPoint float64 `json:"dewPoint"`
        Humidity float64 `json:"humidity"`
        WindSpeed float64 `json:"windSpeed"`
        Visibility float64 `json:"visibility"`
        CloudCover float64 `json:"cloudCover"`
        Pressure float64 `json:"pressure"`
        Ozone float64 `json:"ozone"`
    } `json:"currently"`
    Daily struct {
        Summary string `json:"summary"`
    } `json:"daily"`
    Hourly struct {
        Summary string `json:"summary"`
    } `json:"hourly"`

    GoogleName string // use the google maps FormattedAddress value - it's good
}

var httpClient = &http.Client{Timeout: 10 * time.Second}


// Weather reports the current forecast for a location.
func Weather(m slack.Message, c slack.Conn) {
    // TODO needs to be cached for a minute or so
    reply := getReply(m)
    reply.Channel = m.Channel
    c.Send(m, reply)
}

// eats a slack.Message and returns a slack.Reply with some text set.
func getReply(m slack.Message) (r slack.Reply) {
    token := os.Getenv("DARKSKY_TOKEN")
    if token == "" {
        r.Text = "I got no DARKSKY_TOKEN token in my env. Sort it."
        return
    }

    if m.Subcommand == "" {
        r.Text = "Specify a location, then I can help."
        return
    }

    loc, err := getLocation(m.Subcommand)
    if err != nil {
        r.Text = fmt.Sprintf("Problem getting location: %s", err)
        return
    }

    w, err := getWeather(token, loc)
    if err != nil {
        r.Text = fmt.Sprintf("Problem getting weather: %s", err)
        return
    }

    report := "*%s* :point_right: %s, %.1f°C. %s %s"
    r.Text = fmt.Sprintf(report, w.GoogleName, w.Currently.Summary, w.Currently.Temperature, w.Hourly.Summary, w.Daily.Summary)
    return
}

// returns a weatherLocation struct from google maps api.
func getLocation(q string) (loc weatherLocation, err error) {
    url := "http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false"
    url = fmt.Sprintf(url, q) // TODO does `q` need to be sanitized?

    err = getJson(url, &loc)
    if err != nil {
        return
    }

    if len(loc.Results) < 1 {
        err = fmt.Errorf("google maps failed to convert the query into a location")
    }
    return
}

// returns a weatherInfo struct from darksky.
func getWeather(token string, loc weatherLocation) (w weatherInfo, err error) {
    w.GoogleName = loc.Results[0].FormattedAddress
    url := "https://api.darksky.net/forecast/%s/%f,%f?units=si"
    url = fmt.Sprintf(url, token, loc.Results[0].Geometry.Location.Lat, loc.Results[0].Geometry.Location.Lng)

    err = getJson(url, &w)
    if err != nil {
        return
    }

    if w.Timezone == "" {
        err = fmt.Errorf("darksky doesn't have weather data for the location")
    }
    return
}


// pointless scaffolding.
func getJson(url string, target interface{}) error {
    r, err := httpClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()
    if r.StatusCode != 200 {
        return err
    }

    return json.NewDecoder(r.Body).Decode(target)
}
