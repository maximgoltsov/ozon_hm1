package counter

import (
	"encoding/json"
	"expvar"
	"sync"
)

var c *counter

type Request struct {
	Success int
	Failed  int
}

type CounterRequest struct {
	Incoming Request
	Outgoing Request
}

type CounterType struct {
	Requests CounterRequest
	Errors   int
}

type counter struct {
	cnt CounterType
	m   *sync.RWMutex
}

func (c *counter) IncRequestIncomingSuccess() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt.Requests.Incoming.Success++
}

func IncRequestIncomingSuccess() {
	c.IncRequestIncomingSuccess()
}

func (c *counter) IncRequestIncomingFail() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt.Requests.Incoming.Failed++
}

func IncRequestIncomingFail() {
	c.IncRequestIncomingFail()
}

func (c *counter) IncRequestOutgoingSuccess() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt.Requests.Outgoing.Success++
}

func IncRequestOutgoingSuccess() {
	c.IncRequestOutgoingSuccess()
}

func (c *counter) IncRequestOutgoingFail() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt.Requests.Outgoing.Failed++
}

func IncRequestOutgoingFail() {
	c.IncRequestOutgoingFail()
}

func (c *counter) IncErrors() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt.Errors++
}

func IncErrors() {
	c.IncErrors()
}

func (c *counter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	j, err := json.Marshal(c.cnt)
	if err != nil {
		return ""
	}
	return string(j)
}

// func Inc() {
// 	c.Inc()
// }

func init() {
	c = &counter{m: &sync.RWMutex{}}
	expvar.Publish("Counters", c)
}
