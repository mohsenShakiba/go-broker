package rate_controller

import "sync"

type waitGroupRateController struct {
	parallelism int
	l1          sync.Mutex
	l2          sync.Mutex
	isBlocked   bool
	kmap        map[string]bool
}

func (c *waitGroupRateController) WaitOne(key string) {
	c.l1.Lock()

	// if not available
	if c.parallelism <= 0 {
		c.isBlocked = true
		c.l1.Unlock()
		c.l2.Lock()
		c.l1.Lock()
		c.kmap[key] = true
		c.parallelism -= 1
		c.l1.Unlock()
		c.l2.Unlock()
	} else {
		c.parallelism -= 1
		c.kmap[key] = true
		c.l1.Unlock()
	}
}

func (c *waitGroupRateController) ReleaseOne(key string) {
	c.l1.Lock()
	defer c.l1.Unlock()

	_, ok := c.kmap[key]

	if ok {

		if c.isBlocked {
			c.l2.Unlock()
			c.isBlocked = false
		}

		c.parallelism += 1
		delete(c.kmap, key)
	}

}
