package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// VoteBot votes for emoji! :ballot_box_with_check:
//
// Sadly, VoteBot has a sweet tooth and votes for :doughnut: 15% of the time.
//
// When not voting for :doughnut:, VoteBot can’t be bothered to
// pick a favorite, so it picks one at random. C'mon VoteBot, try harder!

type emoji struct {
	Shortcode string
}

func main() {
	webHost := os.Getenv("WEB_HOST")
	if webHost == "" {
		log.Fatalf("WEB_HOST environment variable must me set")
	}

	hostOverride := os.Getenv("HOST_OVERRIDE")

	// setting the the TTL is optional, thus invalid numbers are simply ignored
	timeToLive, _ := strconv.Atoi(os.Getenv("TTL"))
	var _ time.Time = time.Unix(0, 0)

	if timeToLive != 0 {
		_ = time.Now().Add(time.Second * time.Duration(timeToLive))
	}

	// setting the the request rate is optional, thus invalid numbers are simply ignored
	requestRate, _ := strconv.Atoi(os.Getenv("REQUEST_RATE"))
	if requestRate < 1 {
		requestRate = 1
	}

	webURL := "https://" + webHost
	if _, err := url.Parse(webURL); err != nil {
		log.Fatalf("WEB_HOST %s is invalid", webHost)
	}

	/*
		caCertPath, ok := os.LookupEnv("EDG_CA_PATH")
		if !ok {
			panic("EDG_CA_CERT_PATH environment variable must be set")
		}
		// Load CA cert
		rootCAs := x509.NewCertPool()
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			panic(err)
		}
		if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
			panic("failed to append CA cert")
		}

		// Create the credentials and return it
		config := &tls.Config{
			RootCAs: rootCAs,
		}
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: config,
			},
		}
	*/
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TODO: remove this
			},
		},
	}

	for {
		/*
			// check if deadline has been reached, when TTL has been set.
			if (!deadline.IsZero()) && time.Now().After(deadline) {
				fmt.Printf("Time to live of %d seconds reached, completing\n", timeToLive)
				os.Exit(0)
			}
		*/

		time.Sleep(time.Second / time.Duration(requestRate))

		// Get the list of available shortcodes
		shortcodes, err := shortcodes(client, webURL, hostOverride)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Cast a vote
		probability := rand.Float32()
		switch {
		case probability < 0.15:
			err = vote(client, webURL, hostOverride, ":doughnut:")
		default:
			random := shortcodes[rand.Intn(len(shortcodes))]
			err = vote(client, webURL, hostOverride, random)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

func shortcodes(client *http.Client, webURL string, hostOverride string) ([]string, error) {
	url := fmt.Sprintf("%s/api/list", webURL)
	req, _ := http.NewRequest("GET", url, nil)
	if hostOverride != "" {
		req.Host = hostOverride
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emojis []*emoji
	err = json.Unmarshal(bytes, &emojis)
	if err != nil {
		return nil, err
	}

	shortcodes := make([]string, len(emojis))
	for i, e := range emojis {
		shortcodes[i] = e.Shortcode
	}

	return shortcodes, nil
}

func vote(client *http.Client, webURL string, hostOverride string, shortcode string) error {
	fmt.Printf("✔ Voting for %s\n", shortcode)

	url := fmt.Sprintf("%s/api/vote?choice=%s", webURL, shortcode)
	req, _ := http.NewRequest("GET", url, nil)
	if hostOverride != "" {
		req.Host = hostOverride
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
