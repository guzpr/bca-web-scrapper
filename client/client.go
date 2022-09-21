package client

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
)

var collyClient *colly.Collector

func GetCollyClient() *colly.Collector {
	if collyClient != nil {
		return collyClient
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Jar: jar,
	}

	collyClient = colly.NewCollector(
		colly.AllowedDomains("m.klikbca.com"),
		colly.UserAgent("Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/GRK39F) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"),
		colly.CacheDir("./klickbca"),
	)

	collyClient.SetClient(&client)

	return collyClient
}
