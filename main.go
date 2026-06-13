package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Basics-gif/go_crawler/internal/browser"
	"github.com/Basics-gif/go_crawler/internal/storage"
	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
)

const maxIterations = 5

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error load .env file")

	}

	site := os.Getenv("SITE_URL")
	session := os.Getenv("SESSION")

	log.Printf("============VARIAVEIS=================")
	log.Printf("DEBUG -> SESSION: %q", session)
	log.Printf("DEBUG -> SITE_URL: %s", site)
	log.Printf("======================================")

	bm, err := browser.New(session)
	if err != nil {
		log.Fatalf("could not creates browser manager: %v", err)
	}
	defer bm.Close()

	os.MkdirAll("debug", 0755)

	st, err := storage.New("data/videos.db")
	if err != nil {
		log.Fatalf("could not open storage: %v", err)
	}
	defer st.Close()

	currentURL := site

	for i := 0; i < maxIterations; i++ {
		log.Printf("[%d] navegando para %s", i, currentURL)

		if _, err = bm.Page.Goto(currentURL); err != nil {
			log.Fatalf("could not get entries: %v", err)
		}

		time.Sleep(2 * time.Second)

		raw, err := bm.ExtractInitials()
		if err != nil {
			debugPath := fmt.Sprintf("debug/error_%d.png", i)
			if _, screenErr := bm.Page.Screenshot(playwright.PageScreenshotOptions{
				Path:     playwright.String(debugPath),
				FullPage: playwright.Bool(true),
			}); screenErr != nil {
				log.Printf("could not take debug screenshot: %v", screenErr)
			} else {
				log.Printf("screenshot de erro salvo em %s", debugPath)
			}
			log.Fatalf("could not extract initials: %v", err)
		}

		var list *browser.VideoList
		if i == 0 {
			list, err = browser.ParseInitialPage(raw)
		} else {
			list, err = browser.ParseVideoPage(raw)
		}
		if err != nil {
			log.Fatalf("could not parse video list: %v", err)
		}

		if len(list.Videos) == 0 {
			log.Println("nenhum video encontrado, parando")
			break
		}

		if err = st.SaveAll(list); err != nil {
			log.Fatalf("could not save videos: %v", err)
		}
		log.Printf("Salvo %d videos no banco", len(list.Videos))

		out, _ := json.MarshalIndent(list, "", " ")
		path := fmt.Sprintf("debug/video_%d.json", i)
		if err = os.WriteFile(path, out, 0644); err != nil {
			log.Fatalf("could not write %s: %v", path, err)
		}
		log.Printf("salvo %s (%d videos)", path, len(list.Videos))

		currentURL = list.Videos[0].PageURL
		log.Printf("próximo: %s - %s", list.Videos[0].Title, currentURL)
	}
}
