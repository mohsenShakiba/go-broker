package rate_controller

import (
	"fmt"
	"sync"
)

type waitGroupRateController struct {
	parallelism int
	l1          sync.Mutex
	l2          sync.Mutex
	l3          sync.Mutex
	isBlocked   bool
	kmap        map[string]bool
}

func (c *waitGroupRateController) WaitOne(key string) {
	c.l1.Lock()

	// if not available
	if c.parallelism <= 0 {
		c.l3.Lock()
		c.isBlocked = true
		c.l1.Unlock()
		fmt.Println("now blocked")
		c.l2.Lock()
		c.l3.Unlock()
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

		c.l3.Lock()
		if c.isBlocked {
			fmt.Println("now unblocked")
			c.isBlocked = false
			c.l2.Unlock()
		}
		c.l3.Unlock()

		c.parallelism += 1
		delete(c.kmap, key)
	}

}
