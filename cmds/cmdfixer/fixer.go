package cmdfixer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/weirdtales/senor-rosado/slack"
)

type fixerData struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type fixerCurrencyMap map[string]struct {
	Symbol        string  `json:"symbol"`
	SymbolNative  string  `json:"symbol_native"`
	DecimalDigits int     `json:"decimal_digits"`
	Rounding      float64 `json:"rounding"`
	Code          string  `json:"code"`
}

// Register ...
func Register() (r string, h string) {
	r = `^(fixer|fx) (\w\w\w) ([\d.]+) (\w\w\w)`
	h = "`fixer eur 99.95 aud` use fixer.io to convert currency"
	return
}

// Respond ...
func Respond(r slack.Reply, c chan slack.Reply) {
	u := "https://api.fixer.io/latest?base=%s"
	u = fmt.Sprintf(u, url.QueryEscape(r.Matches[0][2]))

	var resp fixerData
	if err := slack.GetJSON(u, &resp); err != nil {
		r.Text = "problem hitting fixer api"
		c <- r
		return
	}

	u = "http://www.localeplanet.com/api/auto/currencymap.json"
	var cmap fixerCurrencyMap
	if err := slack.GetJSON(u, &cmap); err != nil {
		r.Text = "problem hitting localeplanet.com"
		c <- r
		return
	}

	amt, err := strconv.ParseFloat(r.Matches[0][3], 64)
	if err != nil {
		r.Text = fmt.Sprintf("problem parsing amount %s", r.Matches[0][3])
		c <- r
		return
	}

	t := strings.ToUpper(r.Matches[0][4])
	rate, ok := resp.Rates[t]
	if !ok {
		r.Text = fmt.Sprintf("%s :point_left: ¿Qué es esto?", t)
		c <- r
		return
	}

	b := strings.ToUpper(r.Matches[0][2])
	result := amt * rate
	bsym, ok := cmap[b]
	if !ok {
		r.Text = fmt.Sprintf("problem getting symbols for %s", b)
		c <- r
		return
	}
	tsym, ok := cmap[t]
	if !ok {
		r.Text = fmt.Sprintf("problem getting symbols for %s", t)
		c <- r
		return
	}

	r.Text = fmt.Sprintf("%s%.2f :point_right: %s%.2f", bsym.SymbolNative, amt, tsym.SymbolNative, result)
	c <- r
}
