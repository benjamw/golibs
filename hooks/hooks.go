package hooks

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"google.golang.org/appengine/log"
)

var (
	// registry holds the list of registered hooks and listener type for that hook
	registry map[string]string

	// container holds the list of registered listeners and their priorities for each hook
	container map[string]map[int][]Doer

	// priorities holds a cached list of sorted priorities for each hook
	priorities map[string][]int
)

type Doer interface {
	// Do should check the parameters and make sure they comply with
	// the signature of the actual Doer function.
	//
	// The return values are a continue flag and any error.
	//
	// If continue is false, all subsequent Doers will not be called.
	// If continue is true, but the error is not nil, the error will
	// be logged and processing will continue.
	Do(context.Context, ...interface{}) (bool, error)
}

func init() {
	reset()
}

// reset everything
func reset() {
	registry = make(map[string]string, 0)
	container = make(map[string]map[int][]Doer, 0)
	priorities = make(map[string][]int, 0)
}

// Register a new hook and Doer
func Register(hook string, h Doer) {
	if _, ok := registry[hook]; ok {
		panic(HookError{fmt.Sprintf("duplicate hook name (%s) registered", hook)})
	}

	registry[hook] = reflect.TypeOf(h).String()
}

// Listen for a given hook with the given Doer on the given priority
func Listen(hook string, h Doer, priority int) {
	r := reflect.TypeOf(h).String()
	if registry[hook] != r {
		panic(HookError{fmt.Sprintf("%s listener is listening with wrong doer", hook)})
	}

	if container[hook] == nil {
		container[hook] = make(map[int][]Doer, 0)
	}
	if container[hook][priority] == nil {
		container[hook][priority] = make([]Doer, 0)
	}
	container[hook][priority] = append(container[hook][priority], h)

	// sort and cache the priorities
	n := 0
	p := make([]int, len(container[hook]))
	for k := range container[hook] {
		p[n] = int(k)
		n++
	}

	sort.Ints(p)

	priorities[hook] = p
}

// Do the given hook with the given parameters
// return the last continue flag
func Do(hook string, ctx context.Context, p ...interface{}) bool {
	if _, ok := registry[hook]; !ok {
		panic(HookError{fmt.Sprintf("%s hook not found in registry", hook)})
	}

	for _, k := range priorities[hook] {
		for _, v := range container[hook][k] {
			cont, err := v.Do(ctx, p...)
			if err != nil {
				log.Infof(ctx, fmt.Sprintf("a %s hook threw an error: %v", hook, err))
			}

			if !cont {
				return false
			}
		}
	}

	return true
}
