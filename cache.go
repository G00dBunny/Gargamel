package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	NoExpiration Expiration = -1
	DefaultExpiration  Expiration = 1<<42 //from doc approx 3600 seconds -> 1h

)

type Expiration time.Duration

type Pods struct {
	pods 		[]Pod
}

type Pod struct{
	Name string
	Expiration 	int64
}

type Namespace struct{
	Name string
}

// makes the cache invisible but methods usable -> separate external Cache to internal cache :)
type Cache struct {
	*cache
}

type cache struct {
	namespaces 		map[Namespace]Pods
	expirationTime 	time.Duration 
	lock 			sync.Mutex
}

func New(expiration time.Duration) *Cache{
	return &Cache{
		cache: &cache{
			namespaces: 		make(map[Namespace]Pods),
			expirationTime: expiration,
		},
	}
}	


// TODO : add a force parameter for updating? 
func (c *cache) Set (namespace Namespace, listPods []Pod) error {
	c.lock.Lock()

	if _,exists := c.namespaces[namespace]; exists {
		c.lock.Unlock()
		return fmt.Errorf("item %s already exists", namespace)
	}

	c.namespaces[namespace] = Pods{
		pods: listPods,
	}
	c.lock.Unlock()

	return nil
}

func (c * cache) set (namespace Namespace, pod string) error {
	if _,exists :=  c.namespaces[namespace] ; exists {

	}

	return nil
}

func (c *cache) Add (namespace Namespace, pod string) error {
	c.lock.Lock()

	if _,exists := c.namespaces[namespace]; exists {
		c.lock.Unlock()
		return fmt.Errorf("item %s already exists", namespace)
	}

	return nil
}

func (c *cache) Delete (namespace string, pod string){}

//TODO : Add a flushing option

//TODO : Add a janitor smurf using a tick that will check if items are expired

//TODO : Add a expiration smurf that will expire all namespaces

//TODO : Add a replacer smurf that will replace data


func Refresh() *Cache{
	return &Cache{}
}