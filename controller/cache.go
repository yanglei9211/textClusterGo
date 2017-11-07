// cache 1.0 only cache data

package controller

type Getter func(keys Keys) (Result, error)

type Result map[uint64]interface{}
type Keys []uint64

type entry struct {
	key   uint64
	value interface{}
}

type Cache struct {
	getter Getter
	data   map[uint64]interface{}
}

func NewCache() *Cache {
	c := Cache{}
	c.data = make(map[uint64]interface{})
	return &c
}

func (c *Cache) Get(keys Keys, getter Getter) (Result, error) {
	result := make(Result)
	missedKeys := make(Keys, 0)
	for _, key := range keys {
		if e, ok := c.data[key]; ok {
			result[key] = e
		} else {
			missedKeys = append(missedKeys, key)
		}
	}
	if len(missedKeys) == 0 {
		return result, nil
	}

	missedResult, err := getter(missedKeys)
	if err != nil {
		return result, err
	}
	c.Set(missedResult)
	for k, d := range missedResult {
		result[k] = d
	}
	return result, nil
}

func (c *Cache) Set(data Result) {
	for k, v := range data {
		if _, ok := c.data[k]; !ok {
			c.data[k] = v
		}
	}
}

func (c *Cache) Remove(keys Keys) {
	for _, k := range keys {
		if _, ok := c.data[k]; ok {
			delete(c.data, k)
		}
	}
}
