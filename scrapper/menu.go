package scrapper

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"webscrapping/client"

	"github.com/gocolly/colly/v2"
)

func GoToMenu() {
	c := client.GetCollyClient()
	s := c.Clone()

	data := map[string]string{
		"value(actions)": "menu",
	}

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	s.OnRequest(func(r *colly.Request) {
		log.Println("visiting menu page", r.URL.String())
	})

	s.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)

		data := map[string]string{
			"value(actions)": "acct_stmt",
		}

		form := url.Values{}
		for k, v := range data {
			form.Add(k, v)
		}

		GoToAccountMutationPage()
	})

	s.Request("POST", "https://m.klikbca.com/accountstmt.do",
		strings.NewReader(form.Encode()),
		nil,
		http.Header{"Referer": []string{"https://m.klikbca.com/authentication.do'"}},
	)

	s.Wait()
}
