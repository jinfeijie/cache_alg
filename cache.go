package cache

type Interface interface {
	Set(key string, Value interface{}) uint
	Get(Key string) interface{}
	Del(key string) uint
	GetAll() interface{}
}

type Cache struct {
	cacheAlg Interface
}

func NewCache(i Interface) *Cache {
	return &Cache{
		cacheAlg: i,
	}
}

func (c *Cache) Set(key string, Value interface{}) uint {
	return c.cacheAlg.Set(key, Value)
}

func (c *Cache) Get(key string) interface{} {
	return c.cacheAlg.Get(key)
}

func (c *Cache) Del(key string) uint {
	return c.cacheAlg.Del(key)
}

func (c *Cache) GetAll() interface{} {
	return c.cacheAlg.GetAll()
}
