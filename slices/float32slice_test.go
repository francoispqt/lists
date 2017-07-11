package slices

import (
	"fmt"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/stretchr/testify/assert"
)

func TestFloat32Slice(t *testing.T) {
	var test Float32Slice
	test = []float32{1.0, 2.0}

	indexes := test.Indexes()
	assert.Equal(t, []int{0, 1}, indexes, "indexes should contain the indexes")

	assert.True(t, test.Contains(float32(1.0)), "test should contain hello")
	assert.False(t, test.Contains(float32(3.0)), "test should contain hello")

	forEachT := 0
	test.ForEach(func(k int, v float32) {
		forEachT++
	})

	assert.Equal(t, 2, forEachT, "foreach should have updated forEachT to be 2")

	var test2 Float32Slice
	test2 = test.Map(func(k int, v float32) float32 {
		v += float32(1.0)
		return v
	})

	assert.Equal(t, float32(2.0), test2[0], "should be the same")
	assert.Equal(t, float32(3.0), test2[1], "should be the same")

	testIntf := test.MapInterface(func(k int, v float32) interface{} {
		v += 1.0
		return v
	})

	assert.Equal(t, float32(2.0), testIntf[0].(float32), "should be the same")
	assert.Equal(t, float32(3.0), testIntf[1].(float32), "should be the same")

	var reduce Float32Slice
	reduce = test2.Reduce(func(k int, v float32, agg interface{}) interface{} {
		result := agg.(Float32Slice)
		if v == float32(3.0) {
			result = append(result, v)
			vv := v + 1.0
			result = append(result, vv)
		}
		return result
	}, Float32Slice{}).(Float32Slice)

	assert.Len(t, reduce, 2, "reduce len should be 2")
	assert.Equal(t, float32(3.0), reduce[0], "val should be 3.0")
	assert.Equal(t, float32(4.0), reduce[1], "val should be 4.0")

	// test mapAsync
	// mapAsync might not return the result in the same order
	var ret Float32Slice
	ret = reduce.MapAsync(func(k int, v float32, done chan [2]interface{}) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- [2]interface{}{k, v + 1.0}
	})

	assert.Len(t, ret, 2, "test your be of len 1")
	assert.Equal(t, float32(5.0), ret[1], "Index 1 should be 5.0")

	var retIntf []interface{}
	retIntf = reduce.MapAsyncInterface(func(k int, v float32, done chan [2]interface{}) {
		fmt.Println("async map", k, v)
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- [2]interface{}{k, Float32Slice{v}}
	})

	assert.Len(t, retIntf, len(reduce), "len of retIntf should be same as reduce")
	assert.IsType(t, Float32Slice{}, retIntf[0], "Should of type float32Slice")

	reduceAsync := ret.ReduceAsync(func(k int, v float32, agg *lists.AsyncAggregator) {
		if v == float32(4.0) {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- []float32{2.1}
			return
		}
		agg.Done <- <-agg.Agg
	}).([]float32)

	assert.Len(t, reduceAsync, 1, "should be of len 1")
	assert.Equal(t, float32(2.1), reduceAsync[0], "should be equal to 2.1")

	assert.IsType(t, []float32{}, reduce.Cast(), "cast should give original type")
}
