package lists

type AsyncAggregator struct {
	Agg  chan interface{}
	Done chan interface{}
}
