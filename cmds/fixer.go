package cmds

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

// FixerRegister ...
func FixerRegister() (r string, h string) {
	r = `^(fixer|fx) (\w\w\w) ([\d.]+) (\w\w\w)`
	h = "`fixer eur 99.95 aud` use fixer.io to convert currency"
	return
}

// FixerRespond ...
func FixerRespond(c *slack.Conn, m *slack.Message) {
	fmt.Printf("%+v\n", m.Matches)
	u := "https://api.fixer.io/latest?base=%s"
	u = fmt.Sprintf(u, url.QueryEscape(m.Matches[0][2]))

	var resp fixerData
	if err := slack.GetJSON(u, &resp); err != nil {
		c.ReplyWithErr(m, "problem hitting api.fixer.io", err)
		return
	}

	u = "http://www.localeplanet.com/api/auto/currencymap.json"
	var cmap fixerCurrencyMap
	if err := slack.GetJSON(u, &cmap); err != nil {
		c.ReplyWithErr(m, "problem hitting localeplanet.com", err)
		return
	}

	var msg string
	amt, err := strconv.ParseFloat(m.Matches[0][3], 64)
	if err != nil {
		msg = fmt.Sprintf("problem parsing amount %s", m.Matches[0][3])
		c.ReplyWithErr(m, msg, err)
		return
	}

	t := strings.ToUpper(m.Matches[0][4])
	rate, ok := resp.Rates[t]
	if !ok {
		msg = "%s :point_left: ¿Qué es esto?"
		c.ReplyWithErr(m, fmt.Sprintf(msg, t), err)
		return
	}

	b := strings.ToUpper(m.Matches[0][2])
	result := amt * rate
	bsym, ok := cmap[b]
	if !ok {
		msg = fmt.Sprintf("problem getting symbols for %s", b)
		c.ReplyWithErr(m, msg, err)
		return
	}
	tsym, ok := cmap[t]
	if !ok {
		msg = fmt.Sprintf("problem getting symbols for %s", t)
		c.ReplyWithErr(m, msg, err)
		return
	}

	r := slack.Reply{}
	r.Text = fmt.Sprintf("%s%.2f :point_right: %s%.2f", bsym.SymbolNative, amt, tsym.SymbolNative, result)
	c.ReplyTo(m, r)
}
