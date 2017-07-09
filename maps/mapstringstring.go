package maps

type MapStringString map[string]string

func (c MapStringString) Contains(s string) bool {
	for _, v := range c {
		if v == s {
			return true
		}
	}
	return false
}

func (c MapStringString) ForEach(cb func(string, string)) {
	for k, v := range c {
		cb(k, v)
	}
}

func (c MapStringString) Map(cb func(string, string) string) MapStringString {
	var ret = make(map[string]string, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

func (c MapStringString) MapInterface(cb func(string, string) interface{}) MapStringInterface {
	var ret = make(map[string]interface{}, len(c))
	for k, v := range c {
		ret[k] = cb(k, v)
	}
	return ret
}

func (c MapStringString) MapAsync(cb func(string, string, chan [2]string)) MapStringString {
	mapChan := make(chan [2]string, len(c))
	for k, v := range c {
		go cb(k, v, mapChan)
	}
	var ret = map[string]string{}
	for str := range mapChan {
		if len(str) == 2 {
			ret[str[0]] = str[1]
		}
		if len(ret) == len(c) {
			close(mapChan)
		}
	}
	return ret
}

func (c MapStringString) MapAsyncInterface(cb func(string, string, chan [2]interface{})) MapStringInterface {
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

func (c MapStringString) Reduce(cb func(map[string]string, string, string) map[string]string, aggSlice ...map[string]string) MapStringString {
	var agg map[string]string
	if len(aggSlice) > 0 {
		agg = aggSlice[0]
	} else {
		agg = map[string]string{}
	}
	for k, v := range c {
		agg = cb(agg, k, v)
	}
	return agg
}

func (c MapStringString) Reduce(cb func(interface{}, string, string) interface{}, aggSlice ...interface{}) interface{} {
	var agg interface{}
	if len(aggSlice) > 0 {
		agg = aggSlice[0]
	}
	for k, v := range c {
		agg = cb(agg, k, v)
	}
	return agg
}

func (c MapStringString) ReduceAsync(cb func(chan map[string]string, string, string, chan map[string]string)) MapStringString {
	var done = make(chan map[string]string, len(c))
	var agg = make(chan map[string]string, len(c))
	agg <- make(map[string]string)
	for k, v := range c {
		go cb(agg, k, v, done)
		agg <- <-done
	}
	return <-agg
}

func (c MapStringString) ReduceAsync(cb func(chan interface{}, string, string, chan interface{})) interface{} {
	done := make(chan interface{}, len(c))
	agg := make(chan interface{}, len(c))
	agg <- nil
	for k, v := range c {
		go cb(agg, k, v, done)
		agg <- <-done
	}
	return <-agg
}

func (c MapStringString) IsLast(k string) bool {
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

func (c MapStringString) Indexes() []string {
	var indexes = []string{}
	for k, _ := range c {
		indexes = append(indexes, k)
	}
	return indexes
}

func (c MapStringString) Cast() map[string]string {
	var dest map[string]string
	dest = c
	return dest
}
