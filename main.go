package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
)

type FoundResults struct {
	Results []EmailResult
}

type EmailResult struct {
	FirstName string
	LastName  string
	Domain    string
	Result    string
}
type EmailVerifierRequest struct {
	Requests []EmailFindRequest
}

type EmailFindRequest struct {
	FirstName string
	LastName  string
	Domain    string
}

var (
	verifier = emailverifier.
		NewVerifier().
		EnableSMTPCheck()
)

func cleanName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func parseDomain(d string) string {
	domain := d
	u, err := url.Parse("http://" + domain)
	if err != nil {
		panic(err)
	}

	host := u.Hostname()
	dot := len(host) - 1
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == '.' {
			dot = i
			break
		}
	}

	name := host[:dot]
	tld := host[dot+1:]
	return name + "." + tld

}

func makeRoutes(r *gin.Engine) {
	r.GET("/find_emails", func(c *gin.Context) {

		var requestBody EmailVerifierRequest
		c.BindJSON(&requestBody)

		var results []EmailResult
		for _, req := range requestBody.Requests {
			options := makeOptions(cleanName(req.FirstName), cleanName(req.LastName))
			for i, _ := range options {

				domain := parseDomain(req.Domain)
				username := options[i]
				ret, err := verifier.CheckSMTP(domain, username)
				if err != nil {
					fmt.Println("check smtp failed: ", err)
					continue
				}

				if ret.Deliverable {
					results = append(results, EmailResult{
						FirstName: req.FirstName,
						LastName:  req.LastName,
						Domain:    domain,
						Result:    username + "@" + domain,
					})
				}

			}
		}

		c.JSON(http.StatusOK, results)

	})
}

func main() {
	r := gin.Default()
	makeRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func makeOptions(f_name string, l_name string) []string {
	options := []string{
		string(f_name[0]) + l_name,       //+ "@" + domain,
		string(f_name[0]) + "_" + l_name, //+ "@" + domain,

		string(f_name[0]) + "." + l_name, //+ "@" + domain,
		f_name + "." + l_name,            // + "@" + domain,
		f_name + "." + string(l_name[0]), //+ "@" + domain,
		f_name + "_" + string(l_name[0]), // + "@" + domain,
		f_name + l_name,
		f_name + "_" + l_name, // + "@" + domain,
		f_name,                // + "@" + domain,
		l_name,                // + "@" + domain,
	}
	return options
}
