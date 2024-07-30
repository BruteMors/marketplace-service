package tracing

type MapCarrier map[string][]byte

func (c MapCarrier) Get(key string) string {
	if val, ok := c[key]; ok {
		return string(val)
	}
	return ""
}

func (c MapCarrier) Set(key string, value string) {
	c[key] = []byte(value)
}

func (c MapCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
