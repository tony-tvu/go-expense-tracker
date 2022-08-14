package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/playwright-community/playwright-go"
	"github.com/tony-tvu/goexpense/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExpenseHandler(cfg *config.AppConfig, client *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			testScrape()

			fmt.Fprint(w, "Scrape succeeded")
			return
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

func testScrape() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://news.ycombinator.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	entries, err := page.QuerySelectorAll(".athing")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}
	for i, entry := range entries {
		titleElement, err := entry.QuerySelector("td.title > a")
		if err != nil {
			log.Fatalf("could not get title element: %v", err)
		}
		title, err := titleElement.TextContent()
		if err != nil {
			log.Fatalf("could not get text content: %v", err)
		}
		fmt.Printf("%d: %s\n", i+1, title)
	}
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
