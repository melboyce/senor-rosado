package cmdweather

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/weirdtales/senor-rosado/slack"
)

// partial struct for a google map api hit - so gross.
type weatherLocation struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
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
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		Time                int64   `json:"time"`
		Summary             string  `json:"summary"`
		Icon                string  `json:"icon"`
		Temperature         float64 `json:"temperature"`
		ApparentTemperature float64 `json:"apparentTemperature"`
		DewPoint            float64 `json:"dewPoint"`
		Humidity            float64 `json:"humidity"`
		WindSpeed           float64 `json:"windSpeed"`
		Visibility          float64 `json:"visibility"`
		CloudCover          float64 `json:"cloudCover"`
		Pressure            float64 `json:"pressure"`
		Ozone               float64 `json:"ozone"`
	} `json:"currently"`
	Daily struct {
		Summary string `json:"summary"`
	} `json:"daily"`
	Hourly struct {
		Summary string `json:"summary"`
	} `json:"hourly"`

	GoogleName string // use the google maps FormattedAddress value - it's good
}

// Register ...
func Register() (r string, h string) {
	r = `^(weather) (.+)`
	h = "`weather stalins tomb` get a weather report for a place"
	return
}

// Respond reports the current forecast for a location.
func Respond(r slack.Reply, c chan slack.Reply) {
	// TODO needs to be cached for a minute or so
	var err error
	token := os.Getenv("DARKSKY_TOKEN")
	if token == "" {
		r.Text = "missing Darksky token"
		c <- r
		return
	}

	loc, err := getLocation(strings.Join(r.Args, ","))
	if err != nil {
		r.Text = "unable to convert location"
		c <- r
		return
	}

	w, err := getWeather(token, loc)
	if err != nil {
		r.Text = "unable to get weather"
		c <- r
		return
	}

	tz, err := time.LoadLocation(w.Timezone)
	if err != nil {
		panic(err)
	}
	t := time.Unix(w.Currently.Time, 0)
	ts := t.In(tz).Format("15:04")

	report := "*%s* @ %s :point_right: %s, %.1f°C. %s %s"
	r.Text = fmt.Sprintf(report, w.GoogleName, ts, w.Currently.Summary, w.Currently.Temperature, w.Hourly.Summary, w.Daily.Summary)
	c <- r
}

// returns a weatherLocation struct from google maps api.
func getLocation(q string) (loc weatherLocation, err error) {
	u := "http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false"
	u = fmt.Sprintf(u, url.QueryEscape(q))

	if err = slack.GetJSON(u, &loc); err != nil {
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
	u := "https://api.darksky.net/forecast/%s/%f,%f?units=si"
	u = fmt.Sprintf(u, token, loc.Results[0].Geometry.Location.Lat, loc.Results[0].Geometry.Location.Lng)

	if err = slack.GetJSON(u, &w); err != nil {
		return
	}

	// TODO brittle?
	if w.Timezone == "" {
		err = fmt.Errorf("darksky doesn't have weather data for the location")
	}
	return
}
