package maps

import (
	"reflect"

	"github.com/francoispqt/lists"
)

// InterfaceSlice is a custom type for a slice of interface{}
type MapStringInterface map[string]interface{}

// Contains method determines whether a slice includes a certain element, returning true or false as appropriate.
func (c MapStringInterface) Contains(s interface{}) bool {
	contains := false
	for _, v := range c {
		val := reflect.ValueOf(v)
		sVal := reflect.ValueOf(s)
		switch val.Kind() {
		case reflect.Slice:
			if sVal.Kind() == reflect.Slice {
				contains = compare(v, s)
				if contains {
					return true
				}
			}
			continue
		case reflect.Array:
			if sVal.Kind() == reflect.Array {
				contains = compare(v, s)
				if contains {
					return true
				}
			}
			continue
		default:
			if v == s {
				return true
			}
		}
	}
	return contains
}

// ForEach method executes a provided func once for each slice element.
func (c MapStringInterface) ForEach(cb func(string, interface{})) {
	for k, v := range c {
		cb(k, v)
	}
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of string (original type).
// For asynchronicity, see MapAsync.
func (c MapStringInterface) Map(cb func(string, interface{}) interface{}) MapStringInterface {
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
func (c MapStringInterface) MapAsync(cb func(string, interface{}, chan [2]interface{}), maxConcurrency ...int) MapStringInterface {
	var maxConc = lists.DEFAULT_CONC
	if len(maxConcurrency) == 1 {
		maxConc = maxConcurrency[0]
	}

	mapChan := make(chan [2]interface{}, len(c))
	ret := make(map[string]interface{}, len(c))

	// if maxConc is higher than 0 length of chan doing is lower than maxConc && counter is lower than lenght of slice continue
	// else start reading from the chan to decrease concurrency
	if maxConc > 0 {

		sent := 0
		received := 0
		doing := make(chan struct{}, maxConc)
		indexes := c.Indexes()

		for {

			var k string
			var v interface{}
			if len(indexes) > sent {
				k = indexes[sent]
				v = c[k]
			}

			if len(doing) < maxConc && sent < len(c) {
				go cb(k, v, mapChan)
				sent++
				if maxConc > 0 {
					doing <- struct{}{}
				}
			} else {

				// start reading my chan
				intf := <-mapChan
				received++

				ret[intf[0].(string)] = intf[1]

				// reading doing to continue the loop
				<-doing

				if received == sent {
					close(mapChan)
					return ret
				}
			}
		}
	} else {
		// max concurenccy is 0, means no limit
		// so we only start reading the result chan here
		for k, v := range c {
			go cb(k, v, mapChan)
		}

		ct := 0
		for intf := range mapChan {

			ret[intf[0].(string)] = intf[1]

			ct++
			if ct == len(c) {
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
func (c MapStringInterface) Reduce(cb func(string, interface{}, interface{}) interface{}, defAgg ...interface{}) interface{} {
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
func (c MapStringInterface) ReduceAsync(cb func(string, interface{}, *lists.AsyncAggregator), defAgg ...interface{}) interface{} {
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

// Indexes returns a slice of ints with including the indexes of the StringSlice
func (c MapStringInterface) Indexes() []string {
	var indexes = []string{}
	for k, _ := range c {
		indexes = append(indexes, k)
	}
	return indexes
}

// Filter method creates a new slice with all elements that pass the test implemented by the provided function.
func (c MapStringInterface) Filter(cb func(k string, v interface{}) bool) MapStringInterface {
	var ret = make(map[string]interface{}, 0)
	for k, v := range c {
		if cb(k, v) {
			ret[k] = v
		}
	}
	return ret
}

// Cast explicitly cast the StringSlice to a map[string]string type
func (c MapStringInterface) Cast() map[string]interface{} {
	var dest map[string]interface{}
	dest = c
	return dest
}
