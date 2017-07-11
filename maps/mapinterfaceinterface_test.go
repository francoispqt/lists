package maps

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/francoispqt/lists"
	"github.com/francoispqt/lists/slices"
	"github.com/stretchr/testify/assert"
)

func makeMapInterfaceInterface() MapInterfaceInterface {
	myMap := make(map[interface{}]interface{}, 500)
	for i := 0; i <= 499; i++ {
		iAk := strconv.Itoa(i + 1)
		iAv := strconv.Itoa(i)
		myMap[iAk] = iAv
	}
	return myMap
}

func TestHeavyLiftingMapInterfaceInterface(t *testing.T) {

	myMap := makeMapInterfaceInterface()
	// max concurrency is set to 20
	// test is relying on external api, we don't need to stress it too much
	result := myMap.MapAsync(func(k, v interface{}, done chan [2]interface{}) {
		// do some async
		go func() {
			// write response to channel
			// index must be first element
			done <- [2]interface{}{k, v}
		}()
	}, 100)
	assert.Len(t, result, 500, "len should be 500")
	for k, v := range result {
		assert.True(t, (k != "" && v != ""), "None of the walue should be zero val")
	}

	filtered := result.Filter(func(k interface{}, v interface{}) bool {
		kInt, err := strconv.Atoi(k.(string))
		if err != nil {
			panic(err)
		}
		return kInt <= 100
	})

	assert.Len(t, filtered, 100, "len after filter should be 100")

	ctForEach := 0
	filtered.ForEach(func(k interface{}, v interface{}) {
		ctForEach++
	})

	assert.Equal(t, ctForEach, 100, "forEach counter should be 100")

	fmt.Println("Done testing heavy lifting")
}

func TestMapInterfaceInterface(t *testing.T) {
	var test MapInterfaceInterface
	test = map[interface{}]interface{}{
		"hello": "world",
		"foo":   "bar",
	}

	var indexes slices.InterfaceSlice
	indexes = test.Indexes()
	assert.True(t, indexes.Contains("hello"), "indexes should be equal")
	assert.True(t, indexes.Contains("foo"), "indexes should be equal")

	assert.True(t, test.Contains("world"), "should contain world")
	assert.False(t, test.Contains("coffee"), "should contain world")

	// test contains with slice
	var testSlice MapInterfaceInterface
	testSlice = map[interface{}]interface{}{
		"hello": []string{"world"},
		"foo":   []string{"bar"},
		"bar":   [2]string{"bar", "hello"},
		"hey":   TesStruct{Foo: "bar"},
	}
	assert.True(t, testSlice.Contains([]string{"world"}), "should contain []string{\"world\"}")
	assert.True(t, testSlice.Contains([2]string{"bar", "hello"}), "should contain [2]string{\"bar\",\"hello\"}")
	assert.True(t, testSlice.Contains(TesStruct{Foo: "bar"}), "should contain TesStruct{Foo: \"bar\"}")
	assert.False(t, testSlice.Contains([]string{"hello"}), "should not contain []string{\"hello\"}")

	var test2 MapInterfaceInterface
	test2 = test.Map(func(k interface{}, v interface{}) interface{} {
		vStr := v.(string)
		vStr += " world"
		return vStr
	})

	assert.Equal(t, "world world", test2["hello"], "should be the same")
	assert.Equal(t, "bar world", test2["foo"], "should be the same")

	var reduce map[string]string
	reduce = test2.Reduce(
		func(k interface{}, v interface{}, agg interface{}) interface{} {
			result := agg.(map[string]string)
			if k.(string) == "hello" {
				result[k.(string)] = v.(string)
				vv := v.(string) + " !"
				result[k.(string)+"world"] = vv
			}
			return result
		},
		map[string]string{},
	).(map[string]string)

	assert.Equal(t, "world world", reduce["hello"], "should be the same")
	assert.Equal(t, "world world !", reduce["helloworld"], "should be the same")

	mapAsync := test2.MapAsync(func(k interface{}, v interface{}, done chan [2]interface{}) {
		if k == "hello" {
			time.Sleep(time.Second * 1)
			done <- [2]interface{}{k, "foobar"}
		} else {
			done <- [2]interface{}{k, "hello world"}
		}
	})

	assert.Len(t, mapAsync, 2, "mapAsync should be of len 1")
	assert.Equal(t, mapAsync["hello"], "foobar", "should be the same")

	fmt.Println(test2)
	redAsyncIntf := test2.ReduceAsync(func(k interface{}, v interface{}, agg *lists.AsyncAggregator) {
		if strings.Contains(v.(string), "world world") {
			time.Sleep(time.Second * 1)
			<-agg.Agg
			agg.Done <- MapInterfaceInterface{"foo": "bar"}
			return
		}
		agg.Done <- <-agg.Agg
	}).(MapInterfaceInterface)

	fmt.Println(redAsyncIntf, "redAsyncIntf")

	assert.IsType(t, MapInterfaceInterface{}, redAsyncIntf, "should be the same")
	assert.IsType(t, map[interface{}]interface{}{"": ""}, test2.Cast(), "should be the same")
}
