/*
*	Caching but for pods and namespaces for the cloud-bunny project
*
*   Ideas : Count high restart count after 24h
*
*
 */

package gargamel

import (
	"fmt"
	"sync"
	"time"

	"github.com/G00dBunny/Gargamel/smurfs"
)

const (
	NoExpiration Expiration = -1
	DefaultExpiration  Expiration = 0 
	// HourExpiration  Expiration = 1<<42 //from doc approx 3600 seconds -> 1h

)

type Expiration time.Duration

type Pods struct {
	pods 		[]*Pod
	expiration *Expiration
}

type Pod struct{
	Name string
}

type Namespace struct{
	Name string
}

// makes the cache invisible but methods usable -> separate external Cache to internal cache :)
type Cache struct {
	*cache
}


type cache struct {
	namespaces 		map[*Namespace]*Pods
	expirationTime 	*Expiration
	lock 			sync.Mutex
}

func New(expiration *Expiration) *Cache{
	return &Cache{
		cache: &cache{
			namespaces: 	make(map[*Namespace]*Pods),
			expirationTime: expiration,
		},
	}
}	


// TODO : add a force parameter for updating? 

func (c *cache) Set (namespace *Namespace, listPods []*Pod) error {
	c.lock.Lock()

	if _,exists := c.namespaces[namespace]; exists {
		c.lock.Unlock()
		return fmt.Errorf("item %s already exists", namespace.Name)
	}

	c.namespaces[namespace] = &Pods{
		pods: listPods,
	}
	c.lock.Unlock()

	return nil
}

func (c * cache) set (namespace *Namespace, pod *Pod, expiration *Expiration) error {
	var e int64

	if expiration == smurfs.MakePointer(DefaultExpiration){
		expiration = c.expirationTime
	}

	if *expiration > 0 {
		e = time.Now().Add(time.Duration(*expiration)).UnixNano()
	}


	podlist := c.namespaces[namespace]
	podlist.pods = append(podlist.pods[:], pod)
	podlist.expiration = smurfs.MakePointer(Expiration(e))


	c.namespaces[namespace] = podlist


	return nil
}

func (c *cache) Add (namespace *Namespace, pod *Pod, expiration *Expiration) error {
	c.lock.Lock()

	if _,exists := c.namespaces[namespace]; !exists {
		c.lock.Unlock()
		return fmt.Errorf("item %s does not exist", namespace.Name)
	}

	c.set(namespace, pod, expiration)

	c.lock.Unlock()

	return nil
}

func (c *cache) Delete (namespace string, pod string){}

//TODO : Add a flushing option

//TODO : Add a expiration smurf that will expire all namespaces

//TODO : Add a replacer smurf that will replace data


func Refresh() *Cache{
	return &Cache{}
}