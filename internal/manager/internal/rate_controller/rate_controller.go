package rate_controller

// RateController will control the rate of messages that can go through
// this package is required for limiting the number of messages passed to socket clients
// the difference feature this package provides over sync.WaitGroup is the ability to use keys
// so that calling ReleaseOne with same key multiple times won't have any effect
type RateController interface {
	WaitOne(key string)
	ReleaseOne(key string)
}

func New(parallelism int) RateController {
	c := &waitGroupRateController{
		parallelism: parallelism,
		kmap:        make(map[string]bool),
	}

	c.l2.Lock()

	return c
}
