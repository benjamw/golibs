package hooks

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/benjamw/golibs/test"
)

type TestListener struct {
	H func(context.Context, string) (bool, error)
}

func (h *TestListener) Do(ctx context.Context, p ...interface{}) (bool, error) {
	var ok bool

	var s string
	if s, ok = p[0].(string); !ok {
		panic("second parameter of test doer is of invalid type")
	}

	return h.H(ctx, s)
}

func TestMain(m *testing.M) {
	test.InitCtx()
	runVal := m.Run()
	test.ReleaseCtx()
	os.Exit(runVal)
}

func TestRegister(t *testing.T) {
	reset()
	defer reset()

	l := "Register"

	Register(l, &TestListener{})

	if _, ok := registry[l]; !ok {
		t.Fatal("TestRegister: registry is empty")
	}
}

func TestReset(t *testing.T) {
	reset()
	defer reset()

	l := "Reset"

	if _, ok := priorities[l]; len(priorities) > 0 || ok {
		t.Fatal("TestReset: priorities is not empty")
	}

	if _, ok := container[l]; len(container) > 0 || ok {
		t.Fatal("TestReset: container is not empty")
	}

	if _, ok := registry[l]; ok {
		t.Fatal("TestReset: registry is not empty")
	}

	Register(l, &TestListener{})
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, 10)

	if _, ok := priorities[l]; !ok || len(priorities[l]) == 0 {
		t.Fatal("TestReset: priorities is empty")
	}

	if _, ok := container[l]; !ok || len(container[l]) == 0 {
		t.Fatal("TestReset: container is empty")
	}

	if _, ok := registry[l]; !ok {
		t.Fatal("TestReset: registry is empty")
	}

	reset()

	if _, ok := priorities[l]; len(priorities) > 0 || ok {
		t.Fatal("TestReset: refreshed priorities is not empty")
	}

	if _, ok := container[l]; len(container) > 0 || ok {
		t.Fatal("TestReset: refreshed container is not empty")
	}

	if _, ok := registry[l]; len(registry) > 0 || ok {
		t.Fatal("TestReset: refreshed registry is not empty")
	}
}

func TestListen(t *testing.T) {
	reset()
	defer reset()

	l := "Listen"

	Register(l, &TestListener{})

	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, -10)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, 100)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, 1)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, -100)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, 10)

	if len(priorities[l]) != 5 {
		t.Fatalf("TestListen: incorrect number of priorities. Wanted: 5; Got: %d", len(priorities[l]))
	}

	if len(container[l]) != 5 {
		t.Fatalf("TestListen: incorrect number of container. Wanted: 5; Got: %d", len(container[l]))
	}

	var i = -9999
	for _, v := range priorities[l] {
		if v <= i {
			t.Fatalf("TestListen: priorities are not sorted. %d is not less than %d", v, i)
		}

		i = v
	}
}

func TestDo(t *testing.T) {
	ctx := test.GetCtx()
	reset()
	defer reset()

	l := "Do"

	Happened := false

	Register(l, &TestListener{})

	// test doing with no listeners (nothing should happen)
	Do(l, ctx, "test")

	// Note priority order here...
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return false, nil }}, 4)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, nil }}, 1)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) {
		Happened = true
		return true, nil
	}}, 2)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) { return true, fmt.Errorf("error") }}, 3)
	Listen(l, &TestListener{func(ctx context.Context, s string) (bool, error) {
		t.Fatal("TestDo: ran a listener that should have been skipped")
		return true, nil
	}}, 5)

	cont := Do(l, ctx, "test")
	if cont {
		t.Fatal("TestDo: returned a true continue flag when it should have halted")
	}

	if !Happened {
		t.Fatal("TestDo: did not run a listener that should have been run")
	}
}

func TestDoPanic(t *testing.T) {
	ctx := test.GetCtx()
	reset()
	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("TestDoPanic did not panic")
		} else {
			if _, ok := r.(HookError); !ok {
				t.Fatalf("TestDoPanic panicked with the wrong type: %T: %v", r, r)
			}
		}
	}()

	l := "DoPanic"

	// running Do with no registrars should panic
	Do(l, ctx, "test")
}
