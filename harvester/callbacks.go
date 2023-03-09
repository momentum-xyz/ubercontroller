package harvester

import (
	"sync"
)

type Callback *func(p any)

type Callbacks struct {
	Mu   sync.RWMutex
	Data map[BCType]map[HarvesterEvent]map[Callback]bool
}

func NewCallbacks() *Callbacks {
	return &Callbacks{
		Mu:   sync.RWMutex{},
		Data: make(map[BCType]map[HarvesterEvent]map[Callback]bool),
	}
}

func (c *Callbacks) Add(bcType BCType, event HarvesterEvent, f Callback) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.Data[bcType] == nil {
		c.Data[bcType] = map[HarvesterEvent]map[Callback]bool{}
	}
	if c.Data[bcType][event] == nil {
		c.Data[bcType][event] = map[Callback]bool{}
	}

	c.Data[bcType][event][f] = true
}

func (c *Callbacks) Remove(bcType BCType, event HarvesterEvent, f Callback) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.Data[bcType][event][f]; ok == true {
		delete(c.Data[bcType][event], f)
	}
}

func (c *Callbacks) Trigger(bcType BCType, event HarvesterEvent, p any) {
	if _, ok := c.Data[bcType][event]; ok == false {
		return
	}

	for pf := range c.Data[bcType][event] {
		f := *pf
		go f(p)
	}
}