[![Build Status](https://travis-ci.org/francoispqt/lists.svg?branch=master)](https://travis-ci.org/francoispqt/lists)
[![Go Report Card](https://goreportcard.com/badge/github.com/francoispqt/lists)](https://goreportcard.com/report/github.com/francoispqt/lists)
[![codecov](https://codecov.io/gh/francoispqt/lists/branch/master/graph/badge.svg)](https://codecov.io/gh/francoispqt/lists)

# Lists
Utilities for slices, maps, arrays, strings, chans...

Lists is a go package with multiple utilities for all sorts of lists.
Contributions are welcome !

## GoDoc
https://godoc.org/github.com/francoispqt/lists

**Maps**
https://godoc.org/github.com/francoispqt/lists/maps

**Slices**
https://godoc.org/github.com/francoispqt/lists/slices

## Examples
### Maps
```go
import "github.com/francoispqt/lists/maps"
```
Examples are shown with MapStringString which is a map[string]string

**Mapping**
```go
func main() {
    someMap := map[string]string{
        "hello": "world",
        "foo": "bar",
    }

    result := maps.MapStringString(someMap).Map(func(k int, v string) string {
		v += " !"
		return v
	})

    fmt.Println(result) // map[hello:world ! foo:bar !]
}
```

**Reduction**
```go
func reduction(k string, v string, agg interface{}) interface{} {
	result := agg.(maps.MapStringString)
	if v == "world" {
		result[k] = v + " !"
	}
	return result
}

func main() {
	someMap := map[string]string{
		"hello": "world",
		"foo":   "bar",
	}

	result := maps.MapStringString(someMap).Reduce(reduction, make(maps.MapStringString, 1))

	fmt.Println(result) // map[hello:world !]
}
```

**Async Mapping**
In the async version of the Map method, the func passed is called as a go routine (go func).
To map correctly the values to their initial index in the slice, you must write to the channel a slice of interfaces with the first element in the slice being the index.
```go
func main() {
	someMap := map[string]string{
		"hello": "world",
		"foo":   "bar",
	}

	result := maps.MapStringString(someMap).MapAsync(func(k string, v string, done chan [2]string) {
		if k == "hello" {
			time.Sleep(time.Second * 1)
		}
		done <- [2]string{k, v + " !"}
	})

	fmt.Println(result) // map[foo:bar ! hello:world !]
}
```

**Async Reduction**
```go
func asyncReduction(k string, v string, agg *lists.AsyncAggregator) {
	resultIntf := <-agg.Agg
	if v == "world" {
		result := resultIntf.(map[string]string)
		result[v] = "hello !"
		agg.Done <- result
		return
	}
	agg.Done <- resultIntf
}

func main() {
	someMap := map[string]string{
		"hello": "world",
		"foo":   "bar",
	}

	result := maps.MapStringString(someMap).ReduceAsync(asyncReduction, make(map[string]string, 0))

	fmt.Println(result) // map[world:hello !]
}
```

### Slices
```go
import "github.com/francoispqt/lists/slices"
```
Examples are shown with StringSlice which is a []string

**Mapping**
```go
func main() {
    someSlice := []string{"hello", "foo"}

    result := slices.StringSlice(someSlice).Map(func(k int, v string) string {
		v += " world !"
		return v
	})

    fmt.Println(result) // [hello world ! foo world !]
}
```

**Reduction**
```go
func reduction(k int, v string, agg interface{}) interface{} {
    result := agg.(slices.StringSlice)
    if v == "hello" {
        result[0] = v + " world !"
    }
    return result
}

func main() {
    someSlice := []string{"hello", "foo"}

    result := slices.StringSlice(someSlice).Reduce(reduction, make(slices.StringSlice, 1))

    fmt.Println(result) // [hello world !]
}
```

**Async Mapping**
In the async version of the Map method, the func passed is called as a go routine (go func).
To map correctly the values to their initial index in the slice, you must write to the channel a slice of interfaces with the first element in the slice being the index.
```go
func main() {
    someSlice := []string{"hello", "foo"}

    result := slices.StringSlice(someSlice).MapAsync(func(k int, v string, done chan [2]interface{}) {
        if k == 0 {
            time.Sleep(time.Second * 1)
        }
        done <- []interface{}{k, v + " !"}
    })

    fmt.Println(result) // [hello ! foo !]
}
```

**Async Reduction**
```go
func asyncReduction(k int, v string, agg *lists.AsyncAggregator) {
	resultIntf := <-agg.Agg
	if v == "hello" {
		result := resultIntf.(map[string]string)
		result[v] = "world !"
		agg.Done <- result
		return
	}
	agg.Done <- resultIntf
}

func main() {
    someSlice := []string{"hello", "foo"}

    result := slices.StringSlice(someSlice).ReduceAsync(asyncReduction, make(map[string]string, 0))

    fmt.Println(result) // map[hello:world !]
}
```
## Usage and Documentation
1. **[Maps](#maps-1)**
	1. Contains
	2. ForEach
	3. Map
	4. MapInterface
	5. MapAsync
	6. MapAsyncInterface
	7. Reduce
	8. ReduceAsync
	9. Indexes
	10. Filter
	11. Cast

2. **[Slices](#slices-1)**
	1. Contains
	2. ForEach
	3. Map
	4. MapInterface
	5. MapAsync
	6. MapAsyncInterface
	7. Reduce
	8. ReduceAsync
	9. Indexes
	10. Filter
	11. Cast

## Maps
GoDoc: https://godoc.org/github.com/francoispqt/lists/maps

Maps have 6 different types:

* maps.MapInterfaceInterface
* maps.MapStringString
* maps.MapStringInterface
* maps.MapStringInt
* maps.MapStringFloat64
* maps.MapStringFloat32

Functions explained below will use MapStringString, specificities for certain types will be mentioned.

### Contains
Contains method determines whether a slice includes a certain element, returning true or false as appropriate.

```go
var someMap maps.MapStringString
someMap = map[string]string{
	"hello": "world",
	"foo":   "bar",
}

test.Contains("world") // true
test.Contains("coffee") // false
```

### ForEach
ForEach method executes a provided func once for each slice element.
```go
var someMap MapStringString
someMap = map[string]string{
	"hello": "world",
	"foo":   "bar",
}
someSlice.ForEach(func(k, v string) {
	fmt.Println(k, v)
})
```
### Map
Map method creates a new map with the results of calling a provided func on every element in the calling map.
Returns a map of original type.
For asynchronicity, see MapAsync.
For returning a map with interfaces elements, see MapInterface.

```go
var someMap MapStringString
someMap = map[string]string{
	"hello": "world",
	"foo":   "bar",
}

result = someSlice.Map(func(k, v string) string {
	v += " !"
	return v
})

fmt.Println(result) // map[hello: world !, foo: bar !]
```
### MapInterface
MapInterface method creates a new map with the results of calling a provided func on every element in the calling array.
Returns a map of interfaces indexed by strings for all MapString.
This method does not exist for MapStringInterface and MapInterfaceInterface.
For asynchronicity, see MapAsyncInterface.

```go
var someMap MapStringString
someMap = map[string]string{
	"hello": "world",
	"foo":   "bar",
}

result := someSlice.MapInterface(func(k, v string) interface{} {
	return 1 // we return a different type just to show usage
})

fmt.Println(result) // map[hello: 1, foo: 1]
```

### MapAsync
MapAsync method creates a new map with the results of calling a provided go routine on every element in the calling array.
Runs asynchronously and gives a chan [2]string to return results for all MapStringString and [2]interface{} for all other types.
To keep initial order, the first elemt of the [2]interface{} (or [2]string) written to the chan must be the key. The second element muse be a the destionation type.
Returns the original type (example: MapStringString.Map() returns MapStringString).
If you want to map to a map of different types, see MapAsyncInterface.

**Concurrency**
A second optional argument can be passed to MapAsync, it must be an int. This argument sets a max concurrency to the Mapping, meaning when the number of go routine call will equal to the max concurrency, it will wait for values to be consumed from the "done" chan before calling the next go routines.

```go
var someMap maps.MapStringString
someMap = maps.MapStringString{
	"someURI": "http://someuri.com"
	...
}
// let say some slice has 500 elements
// we want to get the result of all requests
// to do it as quick as possible, let's do it in parallel
// but we must make sure the number of file handlers open are not too high or it will panic
// so we use maxConcurrency and we set it to 100

result := someSlice.MapAsync(func(k string, v string, done chan [2]string) {
		// make get request
		rs, err := http.Get(v)
		log.Printf("calling :", v)

		if err != nil {
			panic(err)
		}
		defer rs.Body.Close()

		bodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			panic(err)
		}

		bodyString := string(bodyBytes)
		log.Printf("got response :", bodyString)

		// write response to channel
		// index must be first element
		done <- [2]string{k, bodyString}
}, 100)
```

### MapAsyncInterface
Map async interface is the same as the MapAsync, except that you can map any value to the original index.
Meaning a map[string]string will return MapStringInterface which is map[string]interface{}.

```go
// here for example we return a map string containing the bytes of the response body or an error
result := someMap.MapAsync(func(k string, v string, done chan [2]interface{}) {
		// make get request
		rs, err := http.Get(v)
		log.Printf("calling :", v)

		if err != nil {
			panic(err)
		}
		defer rs.Body.Close()

		bodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			done <- [2]interface{}{k, err}
			return
		}

		done <- [2]interface{}{k, bodyBytes}
}, 100)
```

### Reduce
Reduce method applies a func against an accumulator and each element in the map to reduce it to a single value of any type.
If no accumulator is passed as second argument, default accumulator will be nil
Returns an interface.
For asynchronicity, see ReduceAsync.

```go
var someMap maps.MapStringString
someMap = maps.MapStringString{
	"1": "2",
	"2": "3",
}

result := someMap.Reduce(
	func(k string, v string, agg interface{}) interface{} {
		result := agg.(int)
		intK, _ := strconv.Atoi(k) // we should check the error :)
		intV, _ := strconv.Atoi(v)
		result += intK + intV
		return result
	},
	0,
).(int) // casting the result directly (as we know what's in there)

fmt.Println(result) // 8
```
### ReduceAsync
ReduceAsync method applies a go routine against an accumulator and each element in the map to reduce it to a single value of any type.
The accumulator is a *lists.AsyncAggregator, it is a pointer to a stuct containing two chan :
```go
type AsyncAggregator struct {
	Agg  chan interface{}
	Done chan interface{}
}
```
It returns an interface.
For synchronicity, see Reduce.

You must get the current state of the aggregator by reading from AsyncAggregator.Agg chan
```go
result := <- AsyncAggregator.Agg
```
Then write the next state of the aggregator to the Done channel
```go
AsyncAggregator.Done <- result
```
If you don't want to change the sate of the aggregator between iterations, just write to the Done, the reading from the Agg
```go
AsyncAggregator.Done <- <-AsyncAggregator.Agg
```

```go
var someMap maps.MapStringString
someMap = map[string]string{
	"hello": "world",
	"foo":   "bar",
}

result := someMap.ReduceAsync(func(k, v string, agg *lists.AsyncAggregator) {
	result := <-agg.Agg
	if k == "foo" {
		time.Sleep(time.Second * 1)
		agg.Done <- "bar"
		return
	}
	agg.Done <- result
}).(string)

fmt.Println(result) // bar
```



### Indexes
The Indexes method returns a slice of a given map's indexes (keys).

```go
var someMap MapStringString
someMap = map[string]string{
	"hello": "world",
	"foo":   "bar",
}

var indexes slices.StringSlice
indexes = someMap.Indexes()

fmt.Println(indexes) // [hello, foo]
```

### Filter
Filter method creates a map with all elements that pass the test implemented by the provided function.

```go
someMap := maps.MapStringString{
	"hello": "world",
	"foo": "bar",
}

result := someMap.Filter(func(k, v string) bool {
	return v != "world"
})

fmt.Println(result) // map[foo: bar]
```

### Cast
Cast method explicitly casts the map to its original type.
For example MapStringString.Cast() returns map[string]string.

```go
someMap := maps.MapStringString{
	"hello": "world",
	"foo": "bar",
}

someMap.Cast()
```

## Slices
GoDoc: https://godoc.org/github.com/francoispqt/lists/slices

Slices have 5 different types:

* slices.InterfaceSlice
* slices.StringSlice
* slices.Intslice
* slices.Float32Slice
* slices.Float34Slice

Functions explained below will use StringSlice, specificities for certain types will be mentioned.

### Contains
Contains method determines whether a slice includes a certain element, returning true or false as appropriate.

```go
var someSlice slices.StringSlice
someSlice = []string{"world","bar"}

test.Contains("world") // true
test.Contains("coffee") // false
```

### ForEach
ForEach method executes a provided func once for each slice element.

```go
var someSlice slices.StringSlice
someSlice = []string{"world","bar"}

someSlice.ForEach(func(k, v string) {
	fmt.Println(k, v)
})
```

### Map
Map method creates a new slice with the results of calling a provided func on every element in the calling map.
Returns a map of original type.
For asynchronicity, see MapAsync.
For returning a map with interfaces elements, see MapInterface.

```go
var someSlice slices.StringSlice
someSlice = []string{"world","bar"}

result = someSlice.Map(func(k, v string) string {
	v += " !"
	return v
})

fmt.Println(result) // [world !, bar !]
```

### MapInterface
MapInterface method creates a new map with the results of calling a provided func on every element in the calling array.
Returns a map of interfaces indexed by strings for all MapString.
This method does not exist for MapStringInterface and MapInterfaceInterface.
For asynchronicity, see MapAsyncInterface.

```go
var someSlice slices.StringSlice
someSlice = []string{"world","bar"}

result := someSlice.MapInterface(func(k, v string) interface{} {
	return 1 // we return a different type just to show usage
})

fmt.Println(result) // [1, 1]
```

### MapAsync
MapAsync method creates a new map with the results of calling a provided go routine on every element in the calling array.
Runs asynchronously and gives a chan [2]string to return results for all MapStringString and [2]interface{} for all other types.
To keep initial order, the first elemt of the [2]interface{} (or [2]string) written to the chan must be the key. The second element muse be a the destionation type.
Returns the original type (example: MapStringString.Map() returns MapStringString).
If you want to map to a map of different types, see MapAsyncInterface.

**Concurrency**
A second optional argument can be passed to MapAsync, it must be an int. This argument sets a max concurrency to the Mapping, meaning when the number of go routine call will equal to the max concurrency, it will wait for values to be consumed from the "done" chan before calling the next go routines.

```go
var someSlice slices.StringSlice
someSlice = []string{"http://someuri.com","http://someotheruri.com"}

// let say some slice has 500 elements
// we want to get the result of all requests
// to do it as quick as possible, let's do it in parallel
// but we must make sure the number of file handlers open are not too high or it will panic
// so we use maxConcurrency and we set it to 100

result := someSlice.MapAsync(func(k string, v string, done chan [2]interface{}) {
		// make get request
		rs, err := http.Get(v)
		log.Printf("calling :", v)

		if err != nil {
			panic(err)
		}
		defer rs.Body.Close()

		bodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			panic(err)
		}

		bodyString := string(bodyBytes)
		log.Printf("got response :", bodyString)

		// write response to channel
		// index must be first element
		done <- [2]interface{k, bodyString}
}, 100)
```


### MapAsyncInterface
Map async interface is the same as the MapAsync, except that you can map any value to the original index.
Meaning a map[string]string will return InterfaceSlice which is a[]interface{}.

```go
// here for example we return a map string containing the bytes of the response body or an error
result := someSlice.MapAsync(func(k string, v string, done chan [2]interface{}) {
		// make get request
		rs, err := http.Get(v)
		log.Printf("calling :", v)

		if err != nil {
			panic(err)
		}
		defer rs.Body.Close()

		bodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			done <- [2]interface{}{k, err}
			return
		}

		done <- [2]interface{}{k, bodyBytes}
}, 100)
```

### Reduce
Reduce method applies a func against an accumulator and each element in the map to reduce it to a single value of any type.
If no accumulator is passed as second argument, default accumulator will be nil
Returns an interface.
For asynchronicity, see ReduceAsync.

```go
var someSlice slices.StringSlice
someSlice = []string{"1","2"}

result := someSlice.Reduce(
	func(k string, v string, agg interface{}) interface{} {
		result := agg.(int)
		intV, _ := strconv.Atoi(v) // we should check the error :)
		result += intV
		return result
	},
	0,
).(int) // casting the result directly (as we know what's in there)

fmt.Println(result) // 3
```

### ReduceAsync
ReduceAsync method applies a go routine against an accumulator and each element in the map to reduce it to a single value of any type.
The accumulator is a *lists.AsyncAggregator, it is a pointer to a stuct containing two chan :
```go
type AsyncAggregator struct {
	Agg  chan interface{}
	Done chan interface{}
}
```
It returns an interface.
For synchronicity, see Reduce.

You must get the current state of the aggregator by reading from AsyncAggregator.Agg chan
```go
result := <- AsyncAggregator.Agg
```
Then write the next state of the aggregator to the Done channel
```go
AsyncAggregator.Done <- result
```
If you don't want to change the sate of the aggregator between iterations, just write to the Done, the reading from the Agg
```go
AsyncAggregator.Done <- <-AsyncAggregator.Agg
```

```go
var someSlice slices.StringSlice
someSlice = []string{"foo","bar"}

result := someSlice.ReduceAsync(func(k int, v string, agg *lists.AsyncAggregator) {
	result := <-agg.Agg
	if v == "foo" {
		result = map[string]string{}
		result[v] = "bar"
		time.Sleep(time.Second * 1)
		agg.Done <- result
		return
	}
	agg.Done <- result
}).(map[string]string)

fmt.Println(result) // map[foo: bar]
```

### Indexes
The Indexes method returns a slice of a given map's indexes (keys).

```go
var someSlice slices.StringSlice
someSlice = []string{"foo","bar"}

var indexes slices.IntSlice
indexes = someSlice.Indexes()

fmt.Println(indexes) // [0, 1]
```

### Filter
Filter method creates a slice with all elements that pass the test implemented by the provided function.

```go
var someSlice slices.StringSlice
someSlice = []string{"foo","bar"}

result := someSlice.Filter(func(k, v string) bool {
	return v != "foo"
})

fmt.Println(result) // [bar]
```

### Cast
Cast method explicitly casts the slice to its original type.
For example StringSlice.Cast() returns []string.

```go
var someSlice slices.StringSlice
someSlice = []string{"foo","bar"}

someSlice.Cast()
```

## Tests

The package is thoroughly tested, although it could take a little cleaning and commenting.
Coverage is at 94% for slices and 97% for maps.
Running test takes around a 50 seconds, depending on the computer because it does a lot of async and some sleep.

You can run test for the whole package by running at the root
```bash
go test ./... -cover
```
Or individually in each package
```bash
cd maps
go test -cover
```
```bash
cd slices
go test -cover
```

## Todo

Add Sort method on slices.
Add Find method on slices and maps.
