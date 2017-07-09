# Lists
Utilities for slices, maps, arrays, strings, chans...

Lists is a go package with multiple utilities for all sorts of lists.
Package is still under development, contributions are welcome !

## Examples
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

**Slice reduction**
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

    result := slices.StringSlice(someSlice).MapAsync(func(k int, v string, done chan []interface{}) {
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
### Maps
