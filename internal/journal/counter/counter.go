package counter

import (
	"expvar"
	"strconv"
	"sync"
)

var (
	InputRequests     *counter
	OutputRequests    *counter
	SuccessRequests   *counter
	FailedRequests    *counter
	ErrorRequests     *counter
	HitCacheRequests  *counter
	MissCacheRequests *counter
)

type counter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *counter) Increase() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *counter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func init() {
	InputRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("InputRequests", InputRequests)

	OutputRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("OutputRequests", OutputRequests)

	SuccessRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("SuccessRequests", SuccessRequests)

	FailedRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("FailedRequests", FailedRequests)

	ErrorRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("ErrorRequests", ErrorRequests)

	HitCacheRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("HitCacheRequests", HitCacheRequests)

	MissCacheRequests = &counter{m: &sync.RWMutex{}}
	expvar.Publish("MissCacheRequests", MissCacheRequests)
}
