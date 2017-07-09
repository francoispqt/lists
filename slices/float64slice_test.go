package slices

import (
	"fmt"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/stretchr/testify/assert"
)

func TestFloat64Slice(t *testing.T) {
	var test Float64Slice
	test = []float64{1.0, 2.0}

	indexes := test.Indexes()
	assert.Equal(t, []int{0, 1}, indexes, "indexes should contain the indexes")

	assert.True(t, test.Contains(1.0), "test should contain hello")
	assert.False(t, test.Contains(3.0), "test should contain hello")

	forEachT := 0
	test.ForEach(func(k int, v float64) {
		forEachT++
	})

	assert.Equal(t, 2, forEachT, "foreach should have updated forEachT to be 2")

	var test2 Float64Slice
	test2 = test.Map(func(k int, v float64) float64 {
		v += 1.0
		return v
	})

	assert.Equal(t, 2.0, test2[0], "should be the same")
	assert.Equal(t, 3.0, test2[1], "should be the same")

	testIntf := test2.MapInterface(func(k int, v float64) interface{} {
		v += 1.0
		return v
	})

	assert.Equal(t, 3.0, testIntf[0].(float64), "should be the same")
	assert.Equal(t, 4.0, testIntf[1].(float64), "should be the same")

	var reduce Float64Slice
	reduce = test2.Reduce(func(k int, v float64, agg interface{}) interface{} {
		result := agg.(Float64Slice)
		if v == 3.0 {
			result = append(result, v)
			vv := v + 1.0
			result = append(result, vv)
		}
		return result
	}, Float64Slice{}).(Float64Slice)

	assert.Len(t, reduce, 2, "reduce len should be 2")
	assert.Equal(t, 3.0, reduce[0], "val should be 3.0")
	assert.Equal(t, 4.0, reduce[1], "val should be 4.0")

	// test mapAsync
	// mapAsync might not return the result in the same order
	var ret Float64Slice
	ret = reduce.MapAsync(func(k int, v float64, done chan []interface{}) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- []interface{}{k, v + 1.0}
	})

	assert.Len(t, ret, 2, "test your be of len 1")
	assert.Equal(t, 5.0, ret[1], "Index 1 should be 5.0")

	var retIntf InterfaceSlice
	retIntf = reduce.MapAsyncInterface(func(k int, v float64, done chan []interface{}) {
		fmt.Println("async map", k, v)
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- []interface{}{k, Float64Slice{v}}
	})

	assert.Len(t, retIntf, len(reduce), "len of retIntf should be same as reduce")
	assert.IsType(t, Float64Slice{}, retIntf[0], "Should of type float64Slice")

	reduceAsync := ret.ReduceAsync(func(k int, v float64, agg *lists.AsyncAggregator) {
		if v == 4.0 {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- []float64{2.1}
			return
		}
		agg.Done <- <-agg.Agg
	}).([]float64)

	assert.Len(t, reduceAsync, 1, "should be of len 1")
	assert.Equal(t, 2.1, reduceAsync[0], "should be equal to 2.1")

	assert.IsType(t, []float64{}, reduce.Cast(), "cast should give original type")

}
