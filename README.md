# Lists
Utilities for slices, maps, arrays, strings, chans...

Lists is a go package with multiple utilities for all sorts of lists.
Package is still under development, contributions are welcome !

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
  4. MapAsyncInterface
  5. Reduce
  6. ReduceAsync
  7. Indexes
  8. Filter
  9. Cast

2. **[Slices](#slices-1)**
  1. Contains
  2. ForEach
  3. Map
  4. MapAsyncInterface
  5. Reduce
  6. ReduceAsync
  7. Indexes
  8. Filter
  9. Cast

## Maps
### Contains
### ForEach
### Map
### MapAsyncInterface
### Reduce
### ReduceAsync
### Indexes
### Filter
### Cast

## Slices
### Contains
### ForEach
### Map
### MapAsyncInterface
### Reduce
### ReduceAsync
### Indexes
### Filter
### Cast

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
