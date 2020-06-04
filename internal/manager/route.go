package manager

import (
	"strings"
)

// rules for pattern matching
// * can be replaced with anything event empty string
// example:
// bar/* can be matched against bar/foo or simply bar
// trailing /* wouldn't have any effect
// example:
// bar/*/* is the same as simply bar

type Route struct {
	Path       string
	segments   []string
	cache      map[string]bool
	subscriber *Subscriber
}

func ProcessPath(path string) []string {

	// trim extra suffix
	for {
		if path[len(path)-1] != '/' && path[len(path)-1] != '*' {
			break
		}
		path = strings.TrimSuffix(path, "/")
		path = strings.TrimSuffix(path, "*")
	}

	// split by '/'
	segments := strings.Split(path, "/")

	return segments
}

func (r *Route) DoesMatch(p string) bool {

	// check if path has been processed before
	_, ok := r.cache[p]

	// if cache exists
	if ok {
		return true
	}

	// message segments
	msgSegments := ProcessPath(p)

	// subscriber segment segments
	subSegments := r.segments

	// if message segment is smaller than subscriber segment
	if len(msgSegments) < len(subSegments) {
		return false
	}

	// if subscriber is empty
	if len(subSegments) == 0 {
		return true
	}

	for i := 0; i < len(subSegments); i++ {
		subSegment := subSegments[i]
		msgSegment := msgSegments[i]

		if subSegment == "*" {
			continue
		}

		if subSegment != msgSegment {
			return false
		}

	}

	return true
}
