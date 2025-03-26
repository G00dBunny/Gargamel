package main

import (
	"fmt"
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
	namespaces 		map[string]Pod
	expirationTime 	time.Duration 
	lock 			sync.Mutex
}

func New(expiration time.Duration) *Cache{
	return &Cache{
		cache: &cache{
			namespaces: 		make(map[string]Pod),
			expirationTime: expiration,
		},
	}
}	


// TODO : add a force parameter for updating? 
func (c *cache) Set (namespace string, pod []string) error {
	c.lock.Lock()

	if _,exists := c.namespaces[namespace]; exists {
		c.lock.Unlock()
		return fmt.Errorf("item %s already exists", namespace)
	}

	c.namespaces[namespace] = Pod{
		pods: pod,
	}
	c.lock.Unlock()

	return nil
}

func (c *cache) Add (namespace string, pod string) {}

func (c *cache) Delete (namespace string, pod string){}

//TODO : Add a flushing option

//TODO : Add a janitor smurf using a tick that will check if items are expired

//TODO : Add a expiration smurf that will expire all namespaces

//TODO : Add a replacer smurf that will replace data


func Refresh() *Cache{
	return &Cache{}
}