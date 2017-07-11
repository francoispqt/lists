package maps

import (
	"fmt"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/francoispqt/lists/slices"
	"github.com/stretchr/testify/assert"
)

func TestMapStringFloat32(t *testing.T) {

	var test MapStringFloat32
	test = map[string]float32{
		"hello": float32(1.0),
		"foo":   float32(2.0),
	}

	var indexes slices.StringSlice
	indexes = test.Indexes()
	assert.True(t, indexes.Contains("hello"), "indexes should be equal")
	assert.True(t, indexes.Contains("foo"), "indexes should be equal")

	assert.True(t, test.Contains(float32(1.0)), "should contain world")
	assert.False(t, test.Contains(float32(3.0)), "should contain world")

	var test2 MapStringFloat32
	test2 = test.Map(func(k string, v float32) float32 {
		v++
		return v
	})

	assert.Equal(t, float32(2.0), test2["hello"], "should be the same")
	assert.Equal(t, float32(3.0), test2["foo"], "should be the same")

	testIntf := test2.MapInterface(func(k string, v float32) interface{} {
		return fmt.Sprintf("test%.1f", v)
	})

	assert.Equal(t, "test2.0", testIntf["hello"].(string), "should be the same")
	assert.Equal(t, "test3.0", testIntf["foo"].(string), "should be the same")

	var reduce float32
	reduce = test2.Reduce(
		func(k string, v float32, agg interface{}) interface{} {
			result := agg.(float32)
			result += v
			return result
		},
		float32(0),
	).(float32)

	assert.Equal(t, reduce, float32(5.0), "reduction should be equal to 5")

	mapAsync := test2.MapAsync(func(k string, v float32, done chan [2]interface{}) {
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
	assert.Equal(t, mapAsync["hello"], float32(3.0), "should be equal to 3")

	mapAsyncIntf := test2.MapAsyncInterface(func(k string, v float32, done chan [2]interface{}) {
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
	redAsyncIntf := test2.ReduceAsync(func(k string, v float32, agg *lists.AsyncAggregator) {
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
	assert.IsType(t, map[string]float32{}, test2.Cast(), "should be the same")

}
