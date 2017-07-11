package maps

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/francoispqt/lists/slices"
	"github.com/stretchr/testify/assert"
)

func makeMap() MapStringString {
	myMap := make(map[string]string, 500)
	for i := 0; i <= 499; i++ {
		iAk := strconv.Itoa(i + 1)
		iAv := strconv.Itoa(i)
		myMap[iAk] = iAv
	}
	return myMap
}

func TestHeavyLifting(t *testing.T) {

	myMap := makeMap()
	// max concurrency is set to 20
	// test is relying on external api, we don't need to stress it too much
	result := myMap.MapAsync(func(k, v string, done chan [2]string) {
		// do some async
		go func() {
			// build uri
			uri := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments/%s", k)
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
			log.Printf("got response :", uri)
			// write response to channel
			// index must be first element
			done <- [2]string{k, bodyString}
		}()
	}, 20)
	assert.Len(t, result, 500, "len should be 500")
	for k, v := range result {
		fmt.Println(k, v)
		assert.True(t, (k != "" && v != ""), "None of the walue should be zero val")
	}

	resultIntf := myMap.MapAsyncInterface(func(k string, v string, done chan [2]interface{}) {
		// do some async
		go func() {
			// write response to channel
			// index must be first element
			vInt, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			done <- [2]interface{}{k, vInt}
		}()
	}, 100)
	assert.Len(t, resultIntf, 500, "len should be 500")
	assert.IsType(t, 0, resultIntf["1"], "type of values in resultIntf should be int")

	filtered := result.Filter(func(k, v string) bool {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		return kInt <= 100
	})

	assert.Len(t, filtered, 100, "len after filter should be 100")

	ctForEach := 0
	filtered.ForEach(func(k, v string) {
		ctForEach++
	})

	assert.Equal(t, ctForEach, 100, "forEach counter should be 100")

	fmt.Println("Done testing heavy lifting")
}

func TestMapStringString(t *testing.T) {

	var test MapStringString
	test = map[string]string{
		"hello": "world",
		"foo":   "bar",
	}

	var indexes slices.StringSlice
	indexes = test.Indexes()
	assert.True(t, indexes.Contains("hello"), "indexes should be equal")
	assert.True(t, indexes.Contains("foo"), "indexes should be equal")

	assert.True(t, test.Contains("world"), "should contain world")
	assert.False(t, test.Contains("coffee"), "should contain world")

	var test2 MapStringString
	test2 = test.Map(func(k string, v string) string {
		v += " world"
		return v
	})

	assert.Equal(t, "world world", test2["hello"], "should be the same")
	assert.Equal(t, "bar world", test2["foo"], "should be the same")

	testIntf := test2.MapInterface(func(k string, v string) interface{} {
		v += " world"
		return v
	})

	assert.Equal(t, "world world world", testIntf["hello"].(string), "should be the same")
	assert.Equal(t, "bar world world", testIntf["foo"].(string), "should be the same")

	var reduce map[string]string
	reduce = test2.Reduce(
		func(k string, v string, agg interface{}) interface{} {
			result := agg.(map[string]string)
			if k == "hello" {
				result[k] = v
				vv := v + " !"
				result[k+"world"] = vv
			}
			return result
		},
		map[string]string{},
	).(map[string]string)

	assert.Equal(t, "world world", reduce["hello"], "should be the same")
	assert.Equal(t, "world world !", reduce["helloworld"], "should be the same")

	mapAsync := test2.MapAsync(func(k string, v string, done chan [2]string) {
		if k == "hello" {
			time.Sleep(time.Second * 1)
			done <- [2]string{k, "foobar"}
		} else {
			done <- [2]string{k, "hello world"}
		}
	})

	assert.Len(t, mapAsync, 2, "mapAsync should be of len 1")
	assert.Equal(t, mapAsync["hello"], "foobar", "should be the same")

	mapAsyncIntf := test2.MapAsyncInterface(func(k string, v string, done chan [2]interface{}) {
		if k == "hello" {
			time.Sleep(time.Second * 1)
			done <- [2]interface{}{k, MapStringInterface{}}
		} else {
			done <- [2]interface{}{k, MapStringInterface{}}
		}
	})

	assert.Len(t, mapAsyncIntf, 2, "mapAsync should be of len 1")
	assert.IsType(t, MapStringInterface{}, mapAsyncIntf["hello"], "should be the same")

	fmt.Println(test2)
	redAsyncIntf := test2.ReduceAsync(func(k, v string, agg *lists.AsyncAggregator) {
		if strings.Contains(v, "world world") {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- MapStringInterface{"foo": "bar"}
			return
		}
		agg.Done <- <-agg.Agg
	}).(MapStringInterface)

	fmt.Println(redAsyncIntf, "redAsyncIntf")

	assert.IsType(t, MapStringInterface{}, redAsyncIntf, "should be the same")
	assert.IsType(t, map[string]string{}, test2.Cast(), "should be the same")

}
