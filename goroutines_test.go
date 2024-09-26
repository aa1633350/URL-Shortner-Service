package main

import (
	"fmt"
	"net/http"
)

func main() {

	websitelist := []string{
		"https://lco.dev",
		"https://google.com",
		"https://go.dev",
	}

	for _, website := range websitelist {
		getStatusCode(website)
	}

}
func getStatusCode(website string) {
	res, err := http.Get(website)
	if err != nil {
		fmt.Println("Oops in endpoint")

	}
	fmt.Printf("%d status code for website %s\n", res.StatusCode, website)
}
