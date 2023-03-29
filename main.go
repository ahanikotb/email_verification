package main

import (
	"net/http"
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
	Status    string
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

func cleanFName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
func cleanLName(s string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(s)), " ", "")
}
func parseDomain(d string) string {
	d = strings.Split(strings.Replace(strings.Replace(strings.Replace(d, "https://", "", -1), "http://", "", -1), "www.", "", -1), "/")[0]
	return d
}
func makeRoutes(r *gin.Engine) {
	// r.GET("/verify_emails", func(ctx *gin.Context) {})
	// r.GET("/verify_email", func(ctx *gin.Context) {})
	r.GET("/find_emails", func(c *gin.Context) {

		var requestBody EmailVerifierRequest
		c.BindJSON(&requestBody)

		var results []EmailResult
		for _, req := range requestBody.Requests {
			options := makeOptions(cleanFName(req.FirstName), cleanLName(req.LastName))
			for i := range options {
				domain := parseDomain(req.Domain)
				username := options[i]
				ret, err := verifier.CheckSMTP(domain, username)
				if err != nil {
					continue
				}

				if ret.Deliverable {
					results = append(results, EmailResult{
						FirstName: req.FirstName,
						LastName:  req.LastName,
						Domain:    domain,
						Result:    username + "@" + domain,
						Status:    "ok",
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
		string(f_name[0]) + l_name,
		string(f_name[0]) + "_" + l_name,
		string(f_name[0]) + "." + l_name,
		f_name + "." + l_name,
		f_name + "." + string(l_name[0]),
		f_name + "_" + string(l_name[0]),
		f_name + l_name,
		f_name + "_" + l_name,
		f_name,
		l_name,
	}
	return options
}

//to send a request send this json body to xxxxx/find_email
// {"Requests": [

//     {
//     "FirstName":"John",
//     "LastName":"smith",
//     "Domain":"apple.com"}
// ]}
