package slices

import (
	"fmt"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/stretchr/testify/assert"
)

func TestInterfaceSlice(t *testing.T) {
	var test InterfaceSlice
	t1 := StringSlice{"test"}
	t2 := StringSlice{"test"}
	test = []interface{}{t1, t2}

	indexes := test.Indexes()
	assert.Equal(t, []int{0, 1}, indexes, "indexes should contain the indexes")

	assert.True(t, test.Contains(t1), "test should contain hello")
	assert.False(t, test.Contains(StringSlice{"test3"}), "test should contain hello")

	forEachT := 0
	test.ForEach(func(k int, v interface{}) {
		forEachT++
	})

	assert.Equal(t, 2, forEachT, "foreach should have updated forEachT to be 2")

	test = test.Map(func(k int, v interface{}) interface{} {
		return StringSlice{fmt.Sprintf("test%d", k)}
	})

	assert.Equal(t, StringSlice{"test0"}, test[0], "should be the same")
	assert.Equal(t, StringSlice{"test1"}, test[1], "should be the same")

	var reduce InterfaceSlice
	reduce = test.Reduce(func(k int, v interface{}, agg interface{}) interface{} {
		result := agg.(InterfaceSlice)
		if v.(StringSlice)[0] == "test1" {
			result = append(result, v)
			vv := StringSlice{"test3"}
			result = append(result, vv)
		}
		return result
	}, InterfaceSlice{}).(InterfaceSlice)

	assert.Len(t, reduce, 2, "reduce len should be 2")
	assert.Equal(t, StringSlice{"test1"}, reduce[0], "val should be 3.0")
	assert.Equal(t, StringSlice{"test3"}, reduce[1], "val should be 4.0")

	// test mapAsync
	// mapAsync might not return the result in the same order
	fmt.Println(reduce)
	var ret InterfaceSlice
	ret = reduce.MapAsync(func(k int, v interface{}, done chan []interface{}) {
		if k == 0 {
			time.Sleep(time.Second * 1)
		}
		fmt.Println("writing to done", k, v)
		done <- []interface{}{k, v}
	})

	assert.Len(t, ret, 2, "test your be of len 2")
	assert.Equal(t, StringSlice{"test1"}, ret[0], "index 0 should be 'StringSlice{\"test1\"}', mapping is async but needs to map back to original index")

	reduceAsync := ret.ReduceAsync(func(k int, v interface{}, agg *lists.AsyncAggregator) {
		if v.(StringSlice)[0] == "test3" {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- StringSlice{"test4"}
			return
		}
		agg.Done <- <-agg.Agg
	}).(StringSlice)

	assert.Len(t, reduceAsync, 1, "should be of len 1")
	assert.Equal(t, "test4", reduceAsync[0], "should be equal to 'StringSlice{\"test4\"}'")

	assert.True(t, ret.IsLast(1), "1 should be last index")
	assert.IsType(t, []interface{}{}, reduce.Cast(), "cast should give original type")
}
