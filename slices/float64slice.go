package slices

import "github.com/francoispqt/lists"

// Float64Slice is a custom type for a slice of float64
type Float64Slice []float64

// Contains method determines whether a slice includes a certain element, returning true or false as appropriate.
func (c Float64Slice) Contains(s float64) bool {
	for _, v := range c {
		if v == s {
			return true
		}
	}
	return false
}

// ForEach method executes a provided func once for each slice element.
func (c Float64Slice) ForEach(cb func(int, float64)) {
	for k, v := range c {
		cb(k, v)
	}
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of string (original type).
// For asynchronicity, see MapAsync.
func (c Float64Slice) Map(cb func(int, float64) float64) Float64Slice {
	var ret = make([]float64, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of interfaces.
// For asynchronicity, see MapAsyncInterface.
func (c Float64Slice) MapInterface(cb func(int, float64) interface{}) InterfaceSlice {
	var ret = make([]interface{}, len(c))
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
func (c Float64Slice) MapAsync(cb func(int, float64, chan []interface{})) Float64Slice {
	mapChan := make(chan []interface{}, len(c))
	for k, v := range c {
		go cb(k, v, mapChan)
	}
	var ret = make([]float64, len(c))
	ct := 0
	for intf := range mapChan {
		if len(intf) > 1 {
			ret[intf[0].(int)] = intf[1].(float64)
		} else {
			ret[intf[0].(int)] = 0
		}
		ct++
		if ct == len(ret) {
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
func (c Float64Slice) MapAsyncInterface(cb func(int, float64, chan []interface{})) InterfaceSlice {
	mapChan := make(chan []interface{}, len(c))
	for k, v := range c {
		go cb(k, v, mapChan)
	}
	var ret = make([]interface{}, len(c))
	ct := 0
	for intf := range mapChan {
		if len(intf) > 1 {
			ret[intf[0].(int)] = intf[1]
		} else {
			ret[intf[0].(int)] = nil
		}
		ct++
		if ct == len(ret) {
			close(mapChan)
		}
	}

	return ret
}

// Reduce method applies a func against an accumulator and each element in the slice (from left to right) to reduce it to a single value of any type.
// If no accumulator is passed as second argument, default accumulator will be nil
// Returns an interface.
// For asynchronicity, see ReduceAsync.
func (c Float64Slice) Reduce(cb func(int, float64, interface{}) interface{}, defAgg ...interface{}) interface{} {
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
func (c Float64Slice) ReduceAsync(cb func(int, float64, *lists.AsyncAggregator), defAgg ...interface{}) interface{} {
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
func (c Float64Slice) IsLast(i int) bool {
	return i == len(c)-1
}

// Indexes returns a slice of ints with including the indexes of the Float64Slice
func (c Float64Slice) Indexes() []int {
	var indexes = []int{}
	for i := 0; i < len(c); i++ {
		indexes = append(indexes, i)
	}
	return indexes
}

// Filter method creates a new slice with all elements that pass the test implemented by the provided function.
func (c Float64Slice) Filter(cb func(k int, v float64) bool) Float64Slice {
	var ret = make([]float64, 0)
	for k, v := range c {
		if cb(k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Cast explicitly cast the StringSlice to a []float64 type
func (c Float64Slice) Cast() []float64 {
	var dest []float64
	dest = c
	return dest
}
