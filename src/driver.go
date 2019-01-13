package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type CoinEvent struct {
	Coin      string
	Event     string
	EventDate string
}

func main() {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains
		colly.AllowedDomains("cryptocalendar.pro"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./page_cache"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1 * time.Second,
		Parallelism: 1,
	})

	coinEvents := make([]CoinEvent, 0, 200)
	stringifyEvents := "Source: https://cryptocalendar.pro \n"

	c.OnHTML(`div[class="col-md-6"]`, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.DOM.Find(`h4[class="bold"]`).Text())
		if len(title) > 0 {
			formattedTitle := strings.Join(strings.Fields(title), " ")
			if strings.Contains(strings.ToLower(formattedTitle), "upcoming events") {
				log.Printf("%v \n", formattedTitle)
				stringifyEvents += fmt.Sprintf("%v \n", formattedTitle)
				e.DOM.Find(`li`).Each(func(i int, sel *goquery.Selection) {
					coinEvent := new(CoinEvent)
					eventDesc := strings.Join(strings.Fields(strings.TrimSpace(sel.Text())), " ")
					if strings.Contains(strings.ToLower(eventDesc), "no upcoming events") {
						stringifyEvents += fmt.Sprintf("  %v\n", "No events")
					} else {
						log.Printf("%v\n", eventDesc)
						stringifyEvents += fmt.Sprintf("  %v\n", eventDesc)
						splittedEventDesc := strings.Split(eventDesc, "—")
						coinEvent.Coin = strings.Replace(e.Request.URL.Path, "/events/", "", -1)

						coinEvent.EventDate = splittedEventDesc[0]
						coinEvent.Event = splittedEventDesc[1]
						coinEvents = append(coinEvents, *coinEvent)
					}
				})
			}

		}
	})

	sites := [...]string{
		"https://cryptocalendar.pro/events/bitcoin",
		"https://cryptocalendar.pro/events/bitcoin-cash",
		"https://cryptocalendar.pro/events/bitcoin-cash-abc",
		"https://cryptocalendar.pro/events/litecoin",
		"https://cryptocalendar.pro/events/ethereum",
		"https://cryptocalendar.pro/events/oax",
		"https://cryptocalendar.pro/events/xrp",
	}
	for i := 0; i < len(sites); i++ {
		c.Visit(sites[i])
		fmt.Println("")
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	//enc.Encode(coinEvents)

	sendEmail(stringifyEvents)
}

func sendEmail(content string) {
	sender := NewSender("brian.yap.test@gmail.com", "<YOUR EMAIL PASSWORD>")

	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{"brianyap@bc.holdings", "brian.yap.sand@gmail.com"}

	Subject := "Cryptocurrency Events"
	message := content
	bodyMessage := sender.WritePlainEmail(Receiver, Subject, message)

	sender.SendMail(Receiver, Subject, bodyMessage)
}
