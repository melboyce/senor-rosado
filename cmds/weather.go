// Package cmds provides commands to a bot
package cmds

import (
    "fmt"
    "os"
    "time"

    "encoding/json"
    "net/http"

    "github.com/weirdtales/senor-rosado/slack"
)


// partial struct for a google map api hit - so gross
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

// partial struct for a darksky api hit - so pretty
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


// Weather reports the current forecast for a location
// The location is queried against the Google maps API
// The location's long and lat are sent to Darksky to get the weather
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

    // location query. not sure what to do about m.Tail...
    q := m.Subcommand
    loc, err := getLocation(q)
    if err != nil {
        return err
    }

    // weather query
    w, err := getWeather(token, loc)
    if err != nil {
        return err
    }

    if w.Timezone == "" {
        r.Text = "hmm, something went wrong... (no data)"
        return nil
    }

    report := "*%s* :point_right: %s %s, %.1f°C. %s %s"
    t := time.Unix(w.Currently.Time, 0).Format("15:04")
    r.Text = fmt.Sprintf(report, w.GoogleName, t, w.Currently.Summary, w.Currently.Temperature, w.Hourly.Summary, w.Daily.Summary)

    return nil

}

// returns a weatherLocation struct from google maps api
func getLocation(q string) (loc weatherLocation, err error) {
    url := "http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false"
    url = fmt.Sprintf(url, q) // TODO does `q` need to be sanitized?
    err = getJson(url, &loc)
    return
}

// returns a weatherInfo struct from darksky
func getWeather(token string, loc weatherLocation) (w weatherInfo, err error) {
    w.GoogleName = loc.Results[0].FormattedAddress
    url := "https://api.darksky.net/forecast/%s/%f,%f?units=si"
    url = fmt.Sprintf(url, token, loc.Results[0].Geometry.Location.Lat, loc.Results[0].Geometry.Location.Lng)
    err = getJson(url, &w)
    return
}


// pointless scaffolding
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
