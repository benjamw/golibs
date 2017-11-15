package db

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/golibs/random"
)

func TestMain(m *testing.M) {
	InitCtx()
	runVal := m.Run()
	ReleaseCtx()
	os.Exit(runVal)
}

// TestSave should be run first, because Save gets used a lot in subsequent tests
func TestSave(t *testing.T) {
	defer ResetDB()
	ctx := GetCtx()

	// everything that needs to be tested gets run in this function
	p := createFoo(ctx, t)

	if p.String == "" {
		t.Fatal("TestSave returned a foo without a string")
	}

	if p.Int == 0 {
		t.Fatal("TestSave returned a foo without an int")
	}
}

func TestLoad(t *testing.T) {
	defer ResetDB()
	ctx := GetCtx()
	var m Foo

	// test loading a non-existent key
	key := datastore.NewKey(ctx, new(Foo).EntityType(), "", 100, nil)

	found, err := Load(ctx, key, &m)
	if err == nil {
		t.Fatal("Load did not throw an error when the object did not exist")
	}
	if found {
		t.Fatal("Load found an object when the object did not exist")
	}

	// test loading an existing key
	p := createFoo(ctx, t)

	found, err = Load(ctx, p.GetKey(), &m)
	if err != nil {
		t.Fatalf("Load threw an error when the object was supposed to exist. Error: %v", err)
	}
	if !found {
		t.Fatal("Load did not find an object when the object was supposed to exist")
	}
}

func TestLoadS(t *testing.T) {
	defer ResetDB()
	ctx := GetCtx()
	var m Foo

	// test loading a non-existent key
	key := datastore.NewKey(ctx, new(Foo).EntityType(), "", 100, nil)

	found, err := LoadS(ctx, key.Encode(), &m)
	if err == nil {
		t.Fatal("LoadS did not throw an error when the object did not exist")
	}
	if found {
		t.Fatal("LoadS found an object when the object did not exist")
	}

	// test loading an existing key
	p := createFoo(ctx, t)

	found, err = LoadS(ctx, p.GetKey().Encode(), &m)
	if err != nil {
		t.Fatalf("LoadS threw an error when the object was supposed to exist. Error: %v", err)
	}
	if !found {
		t.Fatal("LoadS did not find an object when the object was supposed to exist")
	}
}

func TestLoadInt(t *testing.T) {
	defer ResetDB()
	ctx := GetCtx()
	var m Foo

	// test loading a non-existent key
	found, err := LoadInt(ctx, 100, &m)
	if err == nil {
		t.Fatal("LoadInt did not throw an error when the object did not exist")
	}
	if found {
		t.Fatal("LoadInt found an object when the object did not exist")
	}

	// test loading an existing key
	p := createFoo(ctx, t)

	found, err = LoadInt(ctx, p.GetKey().IntID(), &m)
	if err != nil {
		t.Fatalf("LoadInt threw an error when the object was supposed to exist. Error: %v", err)
	}
	if !found {
		t.Fatal("LoadInt did not find an object when the object was supposed to exist")
	}
}

func TestLoadMulti(t *testing.T) {
	defer ResetDB()
	ctx := GetCtx()
	var m []Model

	// test loading existing keys
	keys := make([]*datastore.Key, 0)
	numFoos := 5
	for i := numFoos; i > 0; i-- {
		p := createFoo(ctx, t)
		keys = append(keys, p.GetKey())
	}

	m = new(Foo).Prepare(numFoos)

	found, err := LoadMulti(ctx, keys, m)
	if err != nil {
		t.Fatalf("LoadMulti threw an error when the objects were supposed to exist. Error: %v", err)
	}
	if found != numFoos {
		t.Fatalf("LoadMulti did not find the correct number of objects. Found: %d; Wanted: %d", found, numFoos)
	}
}

// HELPER FUNCTIONS

func createFullFoo(ctx context.Context, t *testing.T, s string, i int64) Foo {
	file, line, funct := GetCaller()

	thing := Foo{
		String: s,
		Int:    i,
	}
	if err := Save(ctx, &thing); err != nil {
		t.Fatalf("Could not save the test Foo. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	// perform a Get to force the key to be applied so it's available in queries
	var get Foo // for use in the forcing Get, not actual data
	err := datastore.Get(ctx, thing.GetKey(), &get)
	if err != nil {
		t.Fatalf("Could not get the test Foo. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
	}

	return thing
}

func createFoo(ctx context.Context, t *testing.T) Foo {
	s := random.Stringn(10)
	i := random.Intn(100, 1000)

	return createFullFoo(ctx, t, s, i)
}

// HELPER STRUCTS
// Because model can't be imported as it imports db

type base struct {
	key   *datastore.Key `datastore:"-"`
	isNew bool           `datastore:"-"`
}

func (b *base) EntityType() string {
	return "BASE"
}

func (b *base) GetKey() *datastore.Key {
	return b.key
}

func (b *base) SetKey(key *datastore.Key) error {
	b.key = key
	return nil
}

func (b *base) IsNew() bool {
	if b.GetKey() == nil {
		return true
	}
	return b.isNew
}

func (b *base) SetIsNew(isNew bool) {
	b.isNew = isNew
}

func (b *base) PreSave(ctx context.Context) error {
	return nil
}

func (b *base) PostSave(ctx context.Context) error {
	return nil
}

func (b *base) PostLoad(ctx context.Context) error {
	fullKey := b.GetKey()
	if fullKey == nil {
		err := &MissingKeyError{}
		return err
	}

	return nil
}

func (b *base) Transform(ctx context.Context, pl datastore.PropertyList) error {
	return nil
}

func (b *base) PreDelete(ctx context.Context) error {
	return nil
}

// Prepare gets a properly sized []Model ready for use in LoadMutliX
func (b *base) Prepare(num int) []Model {
	return nil
}

type Foo struct {
	base
	String string `datastore:",noindex"`
	Int    int64  `datastore:",noindex"`
}

func (m *Foo) EntityType() string {
	return "Foo"
}

func (m *Foo) PreSave(c context.Context) error {
	if m.GetKey() == nil {
		m.SetIsNew(true)
		m.SetKey(datastore.NewIncompleteKey(c, m.EntityType(), nil))
	}
	return nil
}

func (m *Foo) Prepare(n int) []Model {
	f := make([]*Foo, n)
	r := make([]Model, n)
	for k := range f {
		v := new(Foo)
		r[k] = Model(v)
	}

	return r
}

// HELPER TEST FUNCS
// Because test can't be imported as it imports db

var (
	ctx      context.Context
	doneFunc func()
)

// InitCtx initializes the testing context
func InitCtx() {
	c, done, err := aetest.NewContext()
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic("I got in one little fight and my momma got scared")
	}

	ctx = c
	doneFunc = done
}

// GetCtx returns the testing context
func GetCtx() context.Context {
	return ctx
}

// ReleaseCtx processes the doneFunc that clears and releases the testing context
func ReleaseCtx() {
	doneFunc()
}

// ResetDB clears ALL elements from the datastore
func ResetDB() {
	keys, _ := datastore.NewQuery("").KeysOnly().GetAll(ctx, nil)
	DeleteMultiK(ctx, keys)
}

// GetCaller looks up the stack until it finds a func with a name that starts with "Test..."
// and returns the containing file (and two dirs up), the line number, and the name of that func
//
// It can be used in ancillary functions that throw t.Error and t.Fatal to lead one to the location
// of the actual error causing function call instead of the literal location of the t.Error or t.Fatal call
// as follows (in test file, in ancillary func):
//	file, line, funct := test.GetCaller()
// 	t.Fatalf("An error occurred. Func: %s; File: %s; Line: %d; Error: %v", funct, file, line, err)
func GetCaller() (file string, line int, name string) {
	pc := make([]uintptr, 10) // 10 should be enough...
	runtime.Callers(2, pc)    // skip 2 gets us to the caller of this function (GetCaller)

	var fileFunc string
	var lineFunc int
	var nameFunc string
	for _, v := range pc {
		fun := runtime.FuncForPC(v - 1)
		if fun == nil {
			file = "ERROR FINDING FILE"
			line = 0
			name = "???"

			return
		}

		fileFunc, lineFunc = fun.FileLine(v - 1) // -1 because v is for the line below the actual function call
		nameFunc = fun.Name()

		// get actual file name and 2 dirs up "golibs/password/password_test.go" from
		// something like "/home/stuff/golibs/password/password_test.go"
		fileSplit := strings.Split(fileFunc, "/")
		fileSplit = fileSplit[len(fileSplit)-3:]
		fileFunc = strings.Join(fileSplit, "/")

		// get actual function name "TestTokenChangePassword" from
		// something like "github.com/benjamw/golibs/password.TestEncode"
		nameSplit := strings.Split(nameFunc, "/")
		nameFunc = nameSplit[len(nameSplit)-1]
		nameSplit = strings.Split(nameFunc, ".")
		nameFunc = nameSplit[len(nameSplit)-1]

		// stop here if a proper "Test..." function is found
		if "Test" == nameFunc[:4] {
			break
		}
	}

	// if a "Test..." function wasn't found 10 steps above this one,
	// this returns the last function found in the stack
	file = fileFunc
	line = lineFunc
	name = nameFunc

	return
}
