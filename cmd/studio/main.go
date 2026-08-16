package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"studiopeople/internal/people"
)

func main() {
	address := flag.String("addr", ":8080", "HTTP listen address")
	fixturePath := flag.String("fixture", "", "people fixture path")
	flag.Parse()

	file, err := os.Open(resolveFixturePath(*fixturePath))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	clock := people.FixedClock(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	initial, err := people.LoadFixture(file, clock)
	if err != nil {
		log.Fatal(err)
	}
	service := people.NewService(people.NewMemoryRepository(initial), clock)
	log.Printf("studio people service listening on %s", *address)
	if err := http.ListenAndServe(*address, people.NewHandler(service)); err != nil {
		log.Fatal(err)
	}
}

func resolveFixturePath(configured string) string {
	if configured != "" {
		return configured
	}
	for _, candidate := range []string{"fixtures/people.yaml", "../../fixtures/people.yaml"} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return "fixtures/people.yaml"
}
