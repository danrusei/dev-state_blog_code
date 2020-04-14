package site

import "log"

// ContentCache create a key, page relation
type ContentCache map[string]interface{}

//Cache holds the site
type Cache struct {
	cache ContentCache
}

//NewCache initiate a new cache
func NewCache() *Cache {
	return &Cache{
		cache: ContentCache{},
	}
}

//GetKeys retrive keys
func (c *Cache) GetKeys() []string {
	var keys []string
	for key := range c.cache {
		keys = append(keys, key)
	}

	return keys
}

//GetValues provide the values
func (c *Cache) GetValues() []interface{} {
	var values []interface{}
	for _, value := range c.cache {
		values = append(values, value)
	}

	return values
}

//Get for a key return the item(page/post)
func (c *Cache) Get(key string) interface{} {
	item, exists := c.cache[key]
	if exists { // Found an item in the cache
		log.Println("cache hit")
		return item
	}

	log.Println("cache miss/stale")

	return item
}

//Set the key for an item
func (c *Cache) Set(key string, item interface{}) {
	c.cache[key] = item
}
