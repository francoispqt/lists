package slices

import "github.com/francoispqt/lists"

// IntSlice is a custom type for a slice of int
type IntSlice []int

// Contains method determines whether a slice includes a certain element, returning true or false as appropriate.
func (c IntSlice) Contains(s int) bool {
	for _, v := range c {
		if v == s {
			return true
		}
	}
	return false
}

// ForEach method executes a provided func once for each slice element.
func (c IntSlice) ForEach(cb func(int, int)) {
	for k, v := range c {
		cb(k, v)
	}
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of string (original type).
// For asynchronicity, see MapAsync.
func (c IntSlice) Map(cb func(int, int) int) IntSlice {
	var ret = make([]int, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

// MapInterface method creates a new slice with the results of calling a provided func on every element in the calling array.
// Returns a slice of interfaces.
// For asynchronicity, see MapAsyncInterface.
func (c IntSlice) MapInterface(cb func(int, int) interface{}) InterfaceSlice {
	var ret = make([]interface{}, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

// MapAsync method creates a new slice with the results of calling a provided go routine on every element in the calling array.
// Runs asynchronously and needs gives a chan [2]interface{} to return results.
// To keep initial order, the first elemt of th []interface{} written to the chan must be the key. The second element muse be a string.
// Returns a StringSlice (original type).
// If you want to map to a slice of different type, see MapAsyncInterface.
func (c IntSlice) MapAsync(cb func(int, int, chan [2]int), maxConcurrency ...int) IntSlice {

	var maxConc = lists.DEFAULT_CONC
	if len(maxConcurrency) == 1 {
		maxConc = maxConcurrency[0]
	}

	var mapChan = make(chan [2]int, len(c))
	var doing chan struct{}
	var ret = make([]int, len(c))
	var i = 0
	var received = 0

	if maxConc > 0 {
		doing = make(chan struct{}, maxConc)
	} else {
		doing = make(chan struct{})
	}

	for {

		// if maxConc == 0 go ahead or if maxConc is higher than 0 length of chan doing is lower than maxConc && counter is lower than lenght of slice continue
		// else start reading from the chan to decrease concurrency
		if maxConc == 0 || (len(doing) < maxConc && i < len(c)) {
			v := c[i]
			go cb(i, v, mapChan)
			i++
			if maxConc > 0 {
				doing <- struct{}{}
			}
			if maxConc == 0 && i == len(c) {
				break
			}
		} else {

			// start reading my chan
			intf := <-mapChan
			received++

			if len(intf) > 1 {
				ret[intf[0]] = intf[1]
			} else {
				ret[intf[0]] = 0
			}

			// reading doing to continue the loop
			<-doing

			if received == len(c) {
				return ret
			}
		}
	}

	// is max concurenccy is 0, means no limit
	// so we only start reading the result chan here
	if maxConc == 0 {
		ct := 0
		for intf := range mapChan {
			if len(intf) > 1 {
				ret[intf[0]] = intf[1]
			} else {
				ret[intf[0]] = 0
			}
			ct++
			if ct >= len(c) {
				break
			}
		}
	}

	return ret
}

// MapAsyncInterface method creates a new slice with the results of calling a provided go routine on every element in the calling array.
// Runs asynchronously and needs gives a chan [2]interface{} to return results.
// To keep initial order, the first elemt of th []interface{} written to the chan must be the key. The second element muse be a string.
// Returns InterfaceSlice.
// If you know the result will be of original type, user MapAsync.
func (c IntSlice) MapAsyncInterface(cb func(int, int, chan [2]interface{}), maxConcurrency ...int) InterfaceSlice {
	var maxConc = lists.DEFAULT_CONC
	if len(maxConcurrency) == 1 {
		maxConc = maxConcurrency[0]
	}

	var mapChan = make(chan [2]interface{}, len(c))
	var doing chan struct{}
	var ret = make([]interface{}, len(c))
	var i = 0
	var received = 0

	if maxConc > 0 {
		doing = make(chan struct{}, maxConc)
	} else {
		doing = make(chan struct{})
	}

	for {

		// if maxConc == 0 go ahead or if maxConc is higher than 0 length of chan doing is lower than maxConc && counter is lower than lenght of slice continue
		// else start reading from the chan to decrease concurrency
		if maxConc == 0 || (len(doing) < maxConc && i < len(c)) {
			v := c[i]
			go cb(i, v, mapChan)
			i++
			if maxConc > 0 {
				doing <- struct{}{}
			}
			if maxConc == 0 && i == len(c) {
				break
			}
		} else {

			// start reading my chan
			intf := <-mapChan
			received++

			if len(intf) > 1 {
				ret[intf[0].(int)] = intf[1]
			} else {
				ret[intf[0].(int)] = nil
			}

			// reading doing to continue the loop
			<-doing

			if received == len(c) {
				return ret
			}
		}
	}

	// is max concurenccy is 0, means no limit
	// so we only start reading the result chan here
	if maxConc == 0 {
		ct := 0
		for intf := range mapChan {
			if len(intf) > 1 {
				ret[intf[0].(int)] = intf[1]
			} else {
				ret[intf[0].(int)] = nil
			}
			ct++
			if ct >= len(c) {
				break
			}
		}
	}

	return ret
}

// Reduce method applies a func against an accumulator and each element in the slice (from left to right) to reduce it to a single value of any type.
// If no accumulator is passed as second argument, default accumulator will be nil
// Returns an interface.
// For asynchronicity, see ReduceAsync.
func (c IntSlice) Reduce(cb func(int, int, interface{}) interface{}, defAgg ...interface{}) interface{} {
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
func (c IntSlice) ReduceAsync(cb func(int, int, *lists.AsyncAggregator), defAgg ...interface{}) interface{} {
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
func (c IntSlice) IsLast(i int) bool {
	return i == len(c)-1
}

// Indexes returns a slice of ints with including the indexes of the IntSlice
func (c IntSlice) Indexes() []int {
	var indexes = []int{}
	for i := 0; i < len(c); i++ {
		indexes = append(indexes, i)
	}
	return indexes
}

// Filter method creates a new slice with all elements that pass the test implemented by the provided function.
func (c IntSlice) Filter(cb func(k int, v int) bool) IntSlice {
	var ret = make([]int, 0)
	for k, v := range c {
		if cb(k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Cast explicitly cast the StringSlice to a []int type
func (c IntSlice) Cast() []int {
	var dest []int
	dest = c
	return dest
}
