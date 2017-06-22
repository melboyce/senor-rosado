package cmdcurr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/weirdtales/senor-rosado/slack"
)

// localeplanet currency symbol map
type lpCurrencyMap map[string]struct {
	Symbol        string  `json:"symbol"`
	SymbolNative  string  `json:"symbol_native"`
	DecimalDigits int     `json:"decimal_digits"`
	Rounding      float64 `json:"rounding"`
	Code          string  `json:"code"`
}

type currencyMap struct {
	Results map[string]struct {
		Name   string `json:"currencyName"`
		Symbol string `json:"currencySymbol"`
	} `json:"results"`
}

type ccaResult map[string]struct {
	Value float64 `json:"val"`
}

// Register ...
func Register() (r string, h string) {
	r = `^(fixer|fx) (\w\w\w) ([\d.]+) (\w\w\w)`
	h = "`fixer eur 99.95 aud` use fixer.io to convert currency"
	return
}

// Respond ...
func Respond(r slack.Reply, c chan slack.Reply) {
	u := "https://free.currencyconverterapi.com/api/v3/convert?q=%s_%s&compact=y"
	base := strings.ToUpper(r.Matches[0][2])
	tgt := strings.ToUpper(r.Matches[0][4])

	if len(base) != 3 {
		r.Text = fmt.Sprintf("%s :point_left: ¿Qué es esto?", base)
		c <- r
		return
	}

	if len(tgt) != 3 {
		r.Text = fmt.Sprintf("%s :point_left: ¿Qué es esto?", tgt)
		c <- r
		return
	}

	u = fmt.Sprintf(u, base, tgt)

	var resp ccaResult
	if err := slack.GetJSON(u, &resp); err != nil {
		r.Text = "problem hitting api"
		c <- r
		return
	}

	u = "http://www.localeplanet.com/api/auto/currencymap.json"
	u = "http://free.currencyconverterapi.com/api/v3/currencies"
	var cmap currencyMap
	if err := slack.GetJSON(u, &cmap); err != nil {
		r.Text = "problem hitting api (currency map)"
		c <- r
		return
	}

	amt, err := strconv.ParseFloat(r.Matches[0][3], 64)
	if err != nil {
		r.Text = fmt.Sprintf("problem parsing amount %s", r.Matches[0][3])
		c <- r
		return
	}

	key := fmt.Sprintf("%s_%s", base, tgt)
	result := amt * resp[key].Value
	bsym, ok := cmap.Results[base]
	if !ok {
		r.Text = fmt.Sprintf("problem getting symbols for %s", base)
		c <- r
		return
	}
	tsym, ok := cmap.Results[tgt]
	if !ok {
		r.Text = fmt.Sprintf("problem getting symbols for %s", tgt)
		c <- r
		return
	}

	r.Text = fmt.Sprintf("%s%.2f :point_right: %s%.2f", bsym.Symbol, amt, tsym.Symbol, result)
	c <- r
}
