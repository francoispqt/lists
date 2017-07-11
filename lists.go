package lists

const DEFAULT_CONC = 0

type AsyncAggregator struct {
	Agg  chan interface{}
	Done chan interface{}
}
