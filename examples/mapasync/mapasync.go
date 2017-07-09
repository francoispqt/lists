package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/francoispqt/lists/slices"
)

/*
GET/https://jsonplaceholder.typicode.com/comments/1
{
    "postId": 1,
    "id": 1,
    "name": "id labore ex et quam laborum",
    "email": "Eliseo@gardner.biz",
    "body": "laudantium enim quasi est quidem magnam voluptate ipsam eos\ntempora quo necessitatibus\ndolor quam autem quasi\nreiciendis et nam sapiente accusantium"
}

The api has 500 comments
*/

func makeURISlice() []string {
	URISlice := make([]string, 101)
	for i := 0; i <= 100; i++ {
		URISlice[i] = fmt.Sprintf("https://jsonplaceholder.typicode.com/comments/%d", i+1)
	}
	return URISlice
}

func main() {
	slice := slices.StringSlice(makeURISlice())
	result := slice.MapAsync(func(k int, v string, done chan []interface{}) {
		// Make a get request
		rs, err := http.Get(v)
		// Process response
		if err != nil {
			panic(err) // More idiomatic way would be to print the error and die unless it's a serious error
		}
		defer rs.Body.Close()

		bodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			panic(err)
		}

		bodyString := string(bodyBytes)
		done <- []interface{}{k, bodyString}
	}).Cast()
	fmt.Println(result)
}
