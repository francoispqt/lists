package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/francoispqt/lists/slices"
)

//It calls a test api and retrieves all the 500 comments in the API, it keeps a max concurrency at 100 to avoid maxing file handlers limit
func main() {

	start := time.Now()

	result := slices.StringSlice(make([]string, 500)).MapAsync(func(k int, v string, done chan [2]interface{}) {

		// do some async
		go func() {
			// build uri
			uri := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments/%d", k+1)
			log.Printf("calling :", "GET/"+uri)

			// make get request
			rs, err := http.Get(uri)

			if err != nil {
				panic(err) // More idiomatic way would be to print the error and die unless it's a serious error
			}
			defer rs.Body.Close()

			bodyBytes, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				panic(err)
			}

			bodyString := string(bodyBytes)
			log.Printf("got response :", bodyString)
			// write response to channel
			// index must be first element
			done <- [2]interface{}{k, bodyString}
		}()

	}, 100).Cast()

	log.Printf(fmt.Sprintf("Result length : %d", len(result)))
	log.Printf("Map async took %s", time.Since(start))
}
