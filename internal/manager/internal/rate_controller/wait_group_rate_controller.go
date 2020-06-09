package rate_controller

import (
	"sync"
)

type waitGroupRateController struct {
	parallelism int
	l1          sync.Mutex
	ch          chan bool
	kmap        map[string]bool
}

func (c *waitGroupRateController) WaitOne(key string) {
	c.l1.Lock()
	c.kmap[key] = true
	c.l1.Unlock()
	c.ch <- true
}

func (c *waitGroupRateController) ReleaseOne(key string) {
	c.l1.Lock()
	_, ok := c.kmap[key]
	c.l1.Unlock()

	if ok {
		<-c.ch
	}
}
