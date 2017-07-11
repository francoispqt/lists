package slices

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/stretchr/testify/assert"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func TestHeavyLifting(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test := []string{}
	for i := 0; i < 100; i++ {
		test = append(test, fmt.Sprintf("test%d", i))
	}
	ctResult := 0
	result := StringSlice(test).ReduceAsync(
		func(k int, v string, agg *lists.AsyncAggregator) {
			randNum := random(0, 200)
			if randNum < 100 {
				ctResult++
				res := <-agg.Agg
				time.Sleep(time.Duration(randNum) * time.Millisecond)
				resultMap := res.(map[string]string)
				resultMap[v] = strconv.Itoa(randNum)
				agg.Done <- resultMap
				return
			}
			agg.Done <- <-agg.Agg
		},
		make(map[string]string, 0),
	).(map[string]string)

	assert.Len(t, result, ctResult, fmt.Sprintf("Result len should be %d", ctResult))
	fmt.Println("Heavy lifting async reduce success")

	fmt.Println("Start heavy lifting async map")
	resultMapAsync := StringSlice(test).MapAsync(func(k int, v string, done chan [2]interface{}) {
		randNum := random(0, 200)
		time.Sleep(time.Duration(randNum) * time.Millisecond)
		done <- [2]interface{}{k, v}
	}).Cast()

	assert.Len(t, test, len(resultMapAsync), fmt.Sprintf("Result len should be %d", len(test)))
	for i := 0; i < len(test); i++ {
		assert.Equal(t, test[i], resultMapAsync[i], "Values in map async should be same as in original slice")
	}
	fmt.Println("Heavy lifting async map success")
}

func TestStringSlice(t *testing.T) {

	fmt.Println("Starting test string slice")

	var test StringSlice
	test = []string{"hello", "foo"}

	var test3 []string
	test3 = test
	assert.IsType(t, []string{}, test3, "test")

	indexes := test.Indexes()
	assert.Equal(t, []int{0, 1}, indexes, "indexes should contain the indexes")

	assert.True(t, test.Contains("hello"), "test should contain hello")
	assert.False(t, test.Contains("world"), "test should contain hello")

	forEachT := 0
	test.ForEach(func(k int, v string) {
		forEachT++
	})

	assert.Equal(t, 2, forEachT, "foreach should have updated forEachT to be 2")

	var test2 StringSlice
	test2 = test.Map(func(k int, v string) string {
		v += " world"
		return v
	})

	assert.Equal(t, "hello world", test2[0], "should be the same")
	assert.Equal(t, "foo world", test2[1], "should be the same")

	testIntf := test2.MapInterface(func(k int, v string) interface{} {
		v += " world"
		return v
	})

	assert.Equal(t, "hello world world", testIntf[0].(string), "should be the same")
	assert.Equal(t, "foo world world", testIntf[1].(string), "should be the same")

	var reduce StringSlice
	reduce = test2.Reduce(func(k int, v string, agg interface{}) interface{} {
		result := agg.([]string)
		if v == "hello world" {
			result = append(result, v)
			vv := v + " !"
			result = append(result, vv)
		}
		return result
	}, []string{}).([]string)

	assert.Len(t, reduce, 2, "test your be of len 1")
	assert.Equal(t, "hello world", reduce[0], "test your be of len 1")
	assert.Equal(t, "hello world !", reduce[1], "test your be of len 1")

	// test mapAsync
	// mapAsync might not return the result in the same order
	var ret StringSlice
	ret = reduce.MapAsync(func(k int, v string, done chan [2]interface{}) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- [2]interface{}{k, v}
	})

	assert.Len(t, ret, len(reduce), "len of retIntf should be same as reduce")
	assert.IsType(t, "", ret[0], "Should of type stringSlice")

	var retIntf InterfaceSlice
	retIntf = reduce.MapAsyncInterface(func(k int, v string, done chan [2]interface{}) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- [2]interface{}{k, StringSlice{v}}
	})

	assert.Len(t, retIntf, len(reduce), "len of retIntf should be same as reduce")
	assert.IsType(t, StringSlice{}, retIntf[0], "Should of type stringSlice")

	reduceAsync := reduce.ReduceAsync(func(k int, v string, agg *lists.AsyncAggregator) {
		if strings.Contains(v, "world !") {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- []string{"foobar"}
			return
		}
		agg.Done <- <-agg.Agg
	}).([]string)

	assert.Len(t, reduceAsync, 1, "should be of len 1")
	assert.Equal(t, "foobar", reduceAsync[0], "should be equal to foobar")

	assert.IsType(t, []string{}, reduce.Cast(), "cast should give original type")

	fmt.Println("Test string slice done")
}
