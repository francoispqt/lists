package maps

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/francoispqt/lists/slices"
	"github.com/stretchr/testify/assert"
)

func makeMapStringInt() MapStringInt {
	myMap := make(map[string]int, 500)
	for i := 0; i <= 499; i++ {
		iAk := strconv.Itoa(i + 1)
		myMap[iAk] = i + 1
	}
	return myMap
}

func TestHeavyLiftingMapStringInt(t *testing.T) {

	myMap := makeMapStringInt()
	// max concurrency is set to 20
	// test is relying on external api, we don't need to stress it too much
	result := myMap.MapAsync(func(k string, v int, done chan [2]interface{}) {
		// do some async
		go func() {
			// write response to channel
			// index must be first element
			done <- [2]interface{}{k, v}
		}()
	}, 100)
	assert.Len(t, result, 500, "len should be 500")
	for k, v := range result {
		assert.True(t, (k != "" && v != 0), "None of the walue should be zero val")
	}

	resultIntf := myMap.MapAsyncInterface(func(k string, v int, done chan [2]interface{}) {
		// do some async
		go func() {
			// write response to channel
			// index must be first element
			done <- [2]interface{}{k, float32(v)}
		}()
	}, 100)
	assert.Len(t, resultIntf, 500, "len should be 500")
	assert.IsType(t, float32(0), resultIntf["1"], "type of values in resultIntf should be float32")

	filtered := result.Filter(func(k string, v int) bool {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		return kInt <= 100
	})

	assert.Len(t, filtered, 100, "len after filter should be 100")

	ctForEach := 0
	filtered.ForEach(func(k string, v int) {
		ctForEach++
	})

	assert.Equal(t, ctForEach, 100, "forEach counter should be 100")

	fmt.Println("Done testing heavy lifting")
}

func TestMapStringInt(t *testing.T) {

	var test MapStringInt
	test = map[string]int{
		"hello": 1,
		"foo":   2,
	}

	var indexes slices.StringSlice
	indexes = test.Indexes()
	assert.True(t, indexes.Contains("hello"), "indexes should be equal")
	assert.True(t, indexes.Contains("foo"), "indexes should be equal")

	assert.True(t, test.Contains(1), "should contain world")
	assert.False(t, test.Contains(3), "should contain world")

	var test2 MapStringInt
	test2 = test.Map(func(k string, v int) int {
		v++
		return v
	})

	assert.Equal(t, 2, test2["hello"], "should be the same")
	assert.Equal(t, 3, test2["foo"], "should be the same")

	testIntf := test2.MapInterface(func(k string, v int) interface{} {
		return fmt.Sprintf("test%d", v)
	})

	assert.Equal(t, "test2", testIntf["hello"].(string), "should be the same")
	assert.Equal(t, "test3", testIntf["foo"].(string), "should be the same")

	var reduce int
	reduce = test2.Reduce(
		func(k string, v int, agg interface{}) interface{} {
			result := agg.(int)
			result += v
			return result
		},
		0,
	).(int)

	assert.Equal(t, reduce, 5, "reduction should be equal to 5")

	mapAsync := test2.MapAsync(func(k string, v int, done chan [2]interface{}) {
		if k == "hello" {
			time.Sleep(time.Second * 1)
			v++
			done <- [2]interface{}{k, v}
		} else {
			v++
			done <- [2]interface{}{k, v}
		}
	})

	assert.Len(t, mapAsync, 2, "mapAsync should be of len 1")
	assert.Equal(t, mapAsync["hello"], 3, "should be equal to 3")

	mapAsyncIntf := test2.MapAsyncInterface(func(k string, v int, done chan [2]interface{}) {
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
	redAsyncIntf := test2.ReduceAsync(func(k string, v int, agg *lists.AsyncAggregator) {
		if v == 2 {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- MapStringInterface{"foo": "bar"}
			return
		}
		agg.Done <- <-agg.Agg
	}).(MapStringInterface)

	fmt.Println(redAsyncIntf, "redAsyncIntf")

	assert.IsType(t, MapStringInterface{}, redAsyncIntf, "should be the same")
	assert.IsType(t, map[string]int{}, test2.Cast(), "should be the same")

}
