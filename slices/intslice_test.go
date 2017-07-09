package slices

import (
	"fmt"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/stretchr/testify/assert"
)

func TestIntSlice(t *testing.T) {
	var test IntSlice
	test = []int{0, 1}

	indexes := test.Indexes()
	assert.Equal(t, []int{0, 1}, indexes, "indexes should contain the indexes")

	assert.True(t, test.Contains(1), "test should contain hello")
	assert.False(t, test.Contains(2), "test should contain hello")

	forEachT := 0
	test.ForEach(func(k int, v int) {
		forEachT++
	})

	assert.Equal(t, 2, forEachT, "foreach should have updated forEachT to be 2")

	var test2 IntSlice
	test2 = test.Map(func(k int, v int) int {
		v++
		return v
	})

	assert.Equal(t, 1, test2[0], "should be the same")
	assert.Equal(t, 2, test2[1], "should be the same")

	testIntf := test2.MapInterface(func(k int, v int) interface{} {
		v++
		return v
	})

	assert.Equal(t, 2, testIntf[0].(int), "should be the same")
	assert.Equal(t, 3, testIntf[1].(int), "should be the same")

	reduce := test2.Reduce(func(k int, v int, agg interface{}) interface{} {
		result := agg.(IntSlice)
		if v == 2 {
			result = append(result, v)
			vv := v + 1
			result = append(result, vv)
		}
		return result
	}, IntSlice{}).(IntSlice)

	assert.Len(t, reduce, 2, "test your be of len 1")
	assert.Equal(t, 2, reduce[0], "test your be of len 1")
	assert.Equal(t, 3, reduce[1], "test your be of len 1")

	// test mapAsync
	// mapAsync might not return the result in the same order
	var ret IntSlice
	ret = reduce.MapAsync(func(k int, v int, done chan []interface{}) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- []interface{}{k, v}
	})

	assert.Len(t, ret, 2, "test your be of len 1")
	assert.Equal(t, 3, ret[1], "index 1 should be '3'")

	var retIntf InterfaceSlice
	retIntf = reduce.MapAsyncInterface(func(k int, v int, done chan []interface{}) {
		fmt.Println("async map", k, v)
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- []interface{}{k, IntSlice{v}}
	})

	assert.Len(t, retIntf, len(reduce), "len of retIntf should be same as reduce")
	assert.IsType(t, IntSlice{}, retIntf[0], "Should of type stringSlice")

	reduceAsync := reduce.ReduceAsync(func(k int, v int, agg *lists.AsyncAggregator) {
		if v == 3 {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- []int{v}
			return
		}
		agg.Done <- <-agg.Agg
	}).([]int)

	fmt.Println(reduceAsync)
	assert.Len(t, reduceAsync, 1, "should be of len 1")
	assert.Equal(t, 3, reduceAsync[0], "should be equal to 1")

	assert.IsType(t, []int{}, reduce.Cast(), "cast should give original type")
}
