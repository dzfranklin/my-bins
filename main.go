package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"my-bins/parse"
	"net/http"
	"sync"
	"time"
)

//go:embed index.gohtml
var indexTmplSource string

//go:embed normalize.css
var normalizeCSS string

//go:embed style.css
var styleCSS string

const userAgent = "github.com/dzfranklin/my-bins (daniel@danielzfranklin.org)"

var collectionsMu sync.Mutex
var collections []parse.Collection

func main() {
	go func() {
		for {
			err := updateCollections()
			if err != nil {
				log.Printf("error updating: %s", err)
			}
			time.Sleep(24*time.Hour*7 + time.Duration(rand.Intn(int(time.Hour))))
		}
	}()

	indexTmpl := template.Must(template.New("index").Parse(indexTmplSource))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		collectionsMu.Lock()
		defer collectionsMu.Unlock()
		err := indexTmpl.Execute(w, collections)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		_, _ = io.WriteString(w, normalizeCSS+"\n\n\n"+styleCSS)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateCollections() error {
	log.Println("starting update")

	c := http.Client{}

	req1, err := http.NewRequest("GET", "https://www.fife.gov.uk/api/citizen?preview=false&locale=en", nil)
	if err != nil {
		return err
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("User-Agent", userAgent)
	resp1, err := c.Do(req1)
	if err != nil {
		return fmt.Errorf("step1: %w", err)
	}
	_ = resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		return fmt.Errorf("step1: expected status 200, got %d", resp1.StatusCode)
	}
	authz := resp1.Header.Get("Authorization")
	if authz == "" {
		return fmt.Errorf("step1: expected Authorization header in response")
	}

	req2body := []byte(`{"name":"bin_calendar","data":{"uprn":"320073043"},"email":"","caseid":"","xref":"","xref1":"","xref2":""}`)
	req2, err := http.NewRequest("POST", "https://www.fife.gov.uk/api/custom?action=powersuite_bin_calendar_collections&actionedby=bin_calendar&loadform=true&access=citizen&locale=en", bytes.NewReader(req2body))
	if err != nil {
		return err
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("User-Agent", userAgent)
	req2.Header.Set("Authorization", authz)
	resp2, err := c.Do(req2)
	if err != nil {
		return fmt.Errorf("step2: %w", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf("step2: expected status 200, got %d", resp2.StatusCode)
	}
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return fmt.Errorf("step2: %w", err)
	}

	value, err := parse.Parse(body2)
	if err != nil {
		return fmt.Errorf("step2: parse: %w", err)
	}

	collectionsMu.Lock()
	collections = value
	collectionsMu.Unlock()

	log.Println("completed update")

	return nil
}
