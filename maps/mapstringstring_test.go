package maps

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/francoispqt/lists/slices"
	"github.com/stretchr/testify/assert"
)

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
