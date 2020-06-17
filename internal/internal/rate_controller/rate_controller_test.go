package rate_controller

import (
	"testing"
	"time"
)

func TestRateController(t *testing.T) {

	c := New(1)

	c.WaitOne("t1")
	c.ReleaseOne("t1")

	isT2Released := false

	c.WaitOne("t2")
	time.AfterFunc(time.Second, func() {
		isT2Released = true
		c.ReleaseOne("t2")
		t.Logf("released t2")
	})

	go func() {
		c.WaitOne("t3")
		if !isT2Released {
			t.Fatalf("the wait one must not return")
		}
		t.Logf("success release")
	}()

	time.Sleep(time.Second * 2)

}
