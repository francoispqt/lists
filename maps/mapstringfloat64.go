package maps

import (
	"fmt"

	"github.com/francoispqt/lists"
)

type MapStringFloat64 map[string]float64

// Contains method determines whether a slice includes a certain element, returning true or false as appropriate.
func (c MapStringFloat64) Contains(s float64) bool {
	for _, v := range c {
		if v == s {
			return true
		}
	}
	return false
}

// ForEach method executes a provided func once for each slice element.
func (c MapStringFloat64) ForEach(cb func(string, float64)) {
	for k, v := range c {
		cb(k, v)
	}
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of string (original type).
// For asynchronicity, see MapAsync.
func (c MapStringFloat64) Map(cb func(string, float64) float64) MapStringFloat64 {
	var ret = make(map[string]float64, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of interfaces.
// For asynchronicity, see MapAsyncInterface.
func (c MapStringFloat64) MapInterface(cb func(string, float64) interface{}) MapStringInterface {
	var ret = make(map[string]interface{}, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

// MapAsync method creates a new slice with the results of calling a provided go routine on every element in the calling array.
// Runs asynchronously and needs gives a chan []interface{} to return results.
// To keep initial order, the first elemt of th []interface{} written to the chan must be the key. The second element muse be a string.
// Returns a StringSlice (original type).
// If you want to map to a slice of different type, see MapAsyncInterface.
func (c MapStringFloat64) MapAsync(cb func(string, float64, chan [2]interface{}), maxConcurrency ...int) MapStringFloat64 {
	mapChan := make(chan [2]interface{}, len(c))
	for k, v := range c {
		go cb(k, v, mapChan)
	}
	var ret = make(map[string]float64, len(c))
	ct := 0
	for intf := range mapChan {
		fmt.Println(intf)
		if len(intf) > 1 {
			ret[intf[0].(string)] = intf[1].(float64)
		} else {
			ret[intf[0].(string)] = 0
		}
		ct++
		if ct == len(c) {
			close(mapChan)
		}
	}
	return ret
}

// MapAsyncInterface method creates a new slice with the results of calling a provided go routine on every element in the calling array.
// Runs asynchronously and needs gives a chan []interface{} to return results.
// To keep initial order, the first elemt of th []interface{} written to the chan must be the key. The second element muse be a string.
// Returns InterfaceSlice.
// If you know the result will be of original type, user MapAsync.
// @Todo implement max concurrency (in case a lot of requests for example)
func (c MapStringFloat64) MapAsyncInterface(cb func(string, float64, chan [2]interface{})) MapStringInterface {
	mapChan := make(chan [2]interface{}, len(c))
	for k, v := range c {
		go cb(k, v, mapChan)
	}
	var ret = map[string]interface{}{}
	for intf := range mapChan {
		if len(intf) == 2 {
			ret[intf[0].(string)] = intf[1]
			if len(ret) == len(c) {
				close(mapChan)
			}
		}
	}
	return ret
}

// Reduce method applies a func against an accumulator and each element in the slice (from left to right) to reduce it to a single value of any type.
// If no accumulator is passed as second argument, default accumulator will be nil
// Returns an interface.
// For asynchronicity, see ReduceAsync.
func (c MapStringFloat64) Reduce(cb func(string, float64, interface{}) interface{}, defAgg ...interface{}) interface{} {
	var agg interface{}
	if len(defAgg) == 0 {
		agg = nil
	} else {
		agg = defAgg[0]
	}
	for k, v := range c {
		agg = cb(k, v, agg)
	}
	return agg
}

// Reduce method applies a go routinge against an accumulator and each element in the slice (from left to right) to reduce it to a single value of any type.
// Returns an interface.
// For synchronicity, see Reduce.
func (c MapStringFloat64) ReduceAsync(cb func(string, float64, *lists.AsyncAggregator), defAgg ...interface{}) interface{} {
	agg := &lists.AsyncAggregator{
		Done: make(chan interface{}, len(c)),
		Agg:  make(chan interface{}, len(c)),
	}
	if len(defAgg) == 0 {
		agg.Agg <- nil
	} else {
		agg.Agg <- defAgg[0]
	}
	for k, v := range c {
		go cb(k, v, agg)
		agg.Agg <- <-agg.Done
	}
	return <-agg.Agg
}

// IsLast checks if the index passed is the last of the slice
func (c MapStringFloat64) IsLast(k string) bool {
	cL := len(c)
	ct := 0
	for kk, _ := range c {
		ct++
		if ct == cL && k == kk {
			return true
		}
	}
	return false
}

// Indexes returns a slice of ints with including the indexes of the StringSlice
func (c MapStringFloat64) Indexes() []string {
	var indexes = []string{}
	for k, _ := range c {
		indexes = append(indexes, k)
	}
	return indexes
}

// Filter method creates a new slice with all elements that pass the test implemented by the provided function.
func (c MapStringFloat64) Filter(cb func(k string, v float64) bool) MapStringFloat64 {
	var ret = make(map[string]float64, 0)
	for k, v := range c {
		if cb(k, v) {
			ret[k] = v
		}
	}
	return ret
}

// Cast explicitly cast the StringSlice to a map[string]string type
func (c MapStringFloat64) Cast() map[string]float64 {
	var dest map[string]float64
	dest = c
	return dest
}