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

func TestHeavyLiftingIntSlice(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test := []int{}
	for i := 0; i < 200; i++ {
		test = append(test, i)
	}
	ctResult := 0
	result := IntSlice(test).ReduceAsync(
		func(k int, v int, agg *lists.AsyncAggregator) {
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
	resultMapAsync := IntSlice(test).MapAsync(func(k int, v int, done chan [2]int) {
		randNum := random(0, 200)
		time.Sleep(time.Duration(randNum) * time.Millisecond)
		done <- [2]int{k, v}
	}, 100).Cast()

	assert.Len(t, test, len(resultMapAsync), fmt.Sprintf("Result len should be %d", len(test)))
	for i := 0; i < len(test); i++ {
		assert.Equal(t, test[i], resultMapAsync[i], "Values in map async should be same as in original slice")
	}

	filtered := IntSlice(resultMapAsync).Filter(func(k int, v int) bool {
		return k < 100
	})

	assert.Len(t, filtered, 100, "len after filter should be 100")
	assert.True(t, filtered.IsLast(len(filtered)-1), "is last should be true")

	resultIntf := filtered.MapAsyncInterface(func(k int, v int, done chan [2]interface{}) {
		// do some async
		go func() {
			// write response to channel
			// index must be first element
			done <- [2]interface{}{k, v}
		}()
	}, 100)

	assert.Len(t, resultIntf, 100, "len should be 100")
	assert.IsType(t, 0, resultIntf[0], "type of values in resultIntf should be int")

	fmt.Println("Heavy lifting async map success")
}

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
	ret = reduce.MapAsync(func(k int, v int, done chan [2]int) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- [2]int{k, v}
	})

	assert.Len(t, ret, 2, "test your be of len 1")
	assert.Equal(t, 3, ret[1], "index 1 should be '3'")

	var retIntf InterfaceSlice
	retIntf = reduce.MapAsyncInterface(func(k int, v int, done chan [2]interface{}) {
		fmt.Println("async map", k, v)
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		done <- [2]interface{}{k, IntSlice{v}}
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
