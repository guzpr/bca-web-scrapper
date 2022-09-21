package scrapper

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"webscrapping/client"

	"github.com/gocolly/colly/v2"
)

func Login() {
	c := client.GetCollyClient()
	s := c.Clone()

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	data := map[string]string{
		"value(user_id)": os.Getenv("USER_ID"),
		"value(pswd)":    os.Getenv("PASSWORD"),
		"value(actions)": "login",
		"value(Submit)":  "LOGIN",
		"value(user_ip)": ip,
	}

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	s.OnRequest(func(r *colly.Request) {
		log.Println("login", r.URL.String())
	})

	s.OnResponse(func(r *colly.Response) {
		html := string(r.Body)
		log.Println("response received", r.StatusCode)
		if strings.Contains(html, "Anda dapat melakukan login kembali setelah 5 menit") {
			log.Println("Error: Login colldown 5 minutes")
			return
		}

		if strings.Contains(html, "Mohon masukkan User ID / Password Anda yg benar") ||
			strings.Contains(html, "User ID harus Alpha Numerik/User ID must be Alpha Numeric") {
			log.Println("Error: Invalid credential")
			return
		}

		GoToMenu()
	})

	s.Request("POST", "http://m.klikbca.com/authentication.do",
		strings.NewReader(form.Encode()),
		nil,
		http.Header{"Referer": []string{"https://m.klikbca.com/login.jsp"}},
	)

	s.Wait()
}
