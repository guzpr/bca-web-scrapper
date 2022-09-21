package scrapper

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"webscrapping/client"
	"webscrapping/helper"

	"github.com/gocolly/colly/v2"
)

func GetMutation() {
	c := client.GetCollyClient()

	var transactions []helper.Transaction

	s := c.Clone()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	day := now.Format("02")
	month := now.Format("01")
	year := now.Format("2006")

	data := map[string]string{
		"value(actions)": "acctstmtview",
		// Bank Account index
		"value(D1)": "0",
		// Don't know what is this for
		"value(r1)":      "1",
		"value(startDt)": day,
		"value(startMt)": month,
		"value(startYr)": year,
		"value(endDt)":   day,
		"value(endMt)":   month,
		"value(endYr)":   year,
	}

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	s.OnResponse(func(r *colly.Response) {
		log.Println("submitting mutation: response received", r.StatusCode)
		log.Println(string(r.Body))
	})

	s.OnHTML("font", func(e *colly.HTMLElement) {
		content, _ := e.DOM.Html()
		if strings.Contains(content, "TIDAK ADA TRANSAKSI") {
			log.Println("Tidak ada tranasksi")
		}
	})

	s.OnHTML("table", func(table *colly.HTMLElement) {
		// Get table in index 4 for transaction table
		if table.Index == 4 {
			table.ForEach("tbody > tr", func(i int, tr *colly.HTMLElement) {
				if i == 0 {
					return
				}

				row := make([]string, 3)

				tr.ForEach("td", func(i int, td *colly.HTMLElement) {
					text, err := td.DOM.Html()
					if err != nil {
						log.Printf("Error when parsing table row :%s", err)
						return
					}

					row[i] = text
				})

				content := strings.Split(row[1], "<br/>")

				trx := &helper.Transaction{
					Date:   row[0],
					Method: helper.TransactionMethod(content[0]),
					Amount: content[len(content)-1],
					Type:   row[2],
				}

				switch trx.Method {
				case helper.TransactionMethodTransferEBanking:
					metadata := content[1]

					if strings.Contains(metadata, string(helper.TransactionMetadataFTFVA)) {
						trx.DestinationAccountName = content[2]
					} else if strings.Contains(metadata, string(helper.TransactionMetadataFTSCY)) {
						trx.DestinationAccountName = content[3]
					}
				case helper.TransactionMethodTarikanATM:
					break
				case helper.TransactionMethodSwitchingCR:
					trx.DestinationAccountName = content[3]
				default:
					log.Printf("Unmapped transaction method : %s", trx.Method)
				}

				transactions = append(transactions, *trx)
			})
		}
	})

	s.Request("POST", "https://m.klikbca.com/accountstmt.do",
		strings.NewReader(form.Encode()),
		nil,
		http.Header{"Referer": []string{"https://m.klikbca.com/accountstmt.do?value(actions)=menu'"}},
	)

	s.Wait()

	if len(transactions) != 0 {
		j, err := json.MarshalIndent(transactions, "", "  ")
		if err != nil {
			log.Fatalf("error when marshaling transaction: %s", err)
		}

		log.Println(string(j))
	}

	s.Wait()
}
