package main

import (
	"fmt"
	"net/http"

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

func makeRoutes(r *gin.Engine) {
	r.POST("/find_emails", func(c *gin.Context) {

		var requestBody EmailVerifierRequest
		c.BindJSON(&requestBody)

		var results FoundResults
		for _, req := range requestBody.Requests {
			options := makeOptions(req.FirstName, req.LastName, req.Domain)
			for i, _ := range options {

				domain := req.Domain
				username := options[i]
				ret, err := verifier.CheckSMTP(domain, username)
				if err != nil {
					fmt.Println("check smtp failed: ", err)
					return
				}

				fmt.Println("smtp validation result: ", ret)

			}

			// results = append(results, options...)
		}

		c.JSON(http.StatusOK, results)

	})
}
func main() {
	r := gin.Default()

	makeRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func makeOptions(f_name string, l_name string, domain string) []string {
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
