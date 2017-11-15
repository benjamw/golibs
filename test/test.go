package test

import (
	"context"
	"runtime"
	"strings"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"

	"github.com/benjamw/golibs/db"
)

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
		panic("InitCtx failed to get a context")
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
	db.DeleteMultiK(ctx, keys)
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
