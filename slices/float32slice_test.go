package slices

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/stretchr/testify/assert"
)

func TestHeavyLiftingFloat32Slice(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test := []float32{}
	for i := 0; i < 200; i++ {
		test = append(test, float32(i))
	}
	ctResult := 0
	result := Float32Slice(test).ReduceAsync(
		func(k int, v float32, agg *lists.AsyncAggregator) {
			randNum := random(0, 200)
			if randNum < 100 {
				ctResult++
				res := <-agg.Agg
				time.Sleep(time.Duration(randNum) * time.Millisecond)
				resultMap := res.(map[string]string)
				resultMap[strconv.Itoa(k)] = strconv.Itoa(randNum)
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
	resultMapAsync := Float32Slice(test).MapAsync(func(k int, v float32, done chan [2]interface{}) {
		randNum := random(0, 200)
		time.Sleep(time.Duration(randNum) * time.Millisecond)
		done <- [2]interface{}{k, v}
	}, 100).Cast()

	assert.Len(t, test, len(resultMapAsync), fmt.Sprintf("Result len should be %d", len(test)))
	for i := 0; i < len(test); i++ {
		assert.Equal(t, test[i], resultMapAsync[i], "Values in map async should be same as in original slice")
	}

	filtered := Float32Slice(resultMapAsync).Filter(func(k int, v float32) bool {
		return k < 100
	})

	assert.Len(t, filtered, 100, "len after filter should be 100")
	assert.True(t, filtered.IsLast(len(filtered)-1), "is last should be true")

	resultIntf := filtered.MapAsyncInterface(func(k int, v float32, done chan [2]interface{}) {
		// do some async
		go func() {
			// write response to channel
			// index must be first element
			done <- [2]interface{}{k, v}
		}()
	}, 100)

	assert.Len(t, resultIntf, 100, "len should be 100")
	assert.IsType(t, float32(0), resultIntf[0], "type of values in resultIntf should be int")

	fmt.Println("Heavy lifting async map success")
}

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
