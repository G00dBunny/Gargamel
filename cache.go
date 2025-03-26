package main

import (
	"sync"
	"time"
)


type Pod struct{
	pods 		[]string
	Expiration 	int64
}

// it makes the cache invisible but we can still access his methods
type Cache struct {
	*cache
}

type cache struct {
	namespace 		map[string]Pod
	expirationTime 	time.Duration 
	lock 			sync.Mutex
}

func New(expiration time.Duration) *Cache{
	return &Cache{
		cache: &cache{
			namespace: 		make(map[string]Pod),
			expirationTime: expiration,
		},
	}
}	

func (c *cache) Set (namespace string, pod string) {}

func (c *cache) Add (namespace string, pod string) {}

func (c *cache) Delete (namespace string, pod string){}

func Refresh() *Cache{
	return &Cache{}
}