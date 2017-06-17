package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/weirdtales/senor-rosado/slack"
)

// Register for giphy cartridge
func Register() (r string, h string) {
	r = "^(fixer|fx) "
	h = "`fixer eur 99.95 aud` use fixer.io to convert currency"
	return
}

type fixerData struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type currencyMap map[string]struct {
	Symbol        string  `json:"symbol"`
	SymbolNative  string  `json:"symbol_native"`
	DecimalDigits int     `json:"decimal_digits"`
	Rounding      float64 `json:"rounding"`
	Code          string  `json:"code"`
}

// Respond finds a Gif on Giphy
func Respond(m slack.Message, c slack.Conn, matches []string) {
	// TODO dumplication of problema
	defer slack.PanicSuppress()
	reply := slack.Reply{}
	reply.Channel = m.Channel

	u := "https://api.fixer.io/latest?base=%s"
	u = fmt.Sprintf(u, url.QueryEscape(m.Subcommand))

	var resp fixerData
	if err := slack.GetJSON(u, &resp); err != nil {
		slack.SendError(c, m, err, "")
		return
	}

	u = "http://www.localeplanet.com/api/auto/currencymap.json"
	var cmap currencyMap
	if err := slack.GetJSON(u, &cmap); err != nil {
		slack.SendError(c, m, err, "")
		return
	}

	reply.Text = "hoy no"
	amt, err := strconv.ParseFloat(string(m.Args[1]), 64)
	if err != nil {
		slack.SendError(c, m, err, "")
		return
	}

	t := strings.ToUpper(string(m.Args[2]))
	rate, ok := resp.Rates[t]
	if !ok {
		slack.SendError(c, m, err, fmt.Sprintf("%s :point_left: ¿Qué es esto?", t))
		c.Send(m, reply)
		return
	}

	b := strings.ToUpper(m.Subcommand)
	result := amt * rate
	bsym, ok := cmap[b]
	if !ok {
		slack.SendError(c, m, err, "")
		return
	}
	tsym, ok := cmap[t]
	if !ok {
		slack.SendError(c, m, err, "")
		return
	}
	reply.Text = fmt.Sprintf("%s%.2f :point_right: %s%.2f", bsym.SymbolNative, amt, tsym.SymbolNative, result)
	c.Send(m, reply)
}
