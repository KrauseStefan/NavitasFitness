package TestHelper

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

func GetContext() context.Context {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		panic(err)
	}
	defer done()
	return ctx
}

type AssertObj struct {
	t        *testing.T
	leftSide interface{}
}

func logError(t *testing.T, str string) {
	t.Error(str)
}

func printLineAndFunction(t *testing.T) {
	stackLines := bytes.Split(debug.Stack(), []byte("\n"))

	lines := stackLines[7:9]

	logError(t, string(bytes.Join(lines, []byte("\n"))))
}

func (a *AssertObj) Equals(rightSide interface{}) {
	leftSide := a.leftSide

	// nil does not equal nil (nil == nil) => false
	if !reflect.DeepEqual(leftSide, rightSide) && (leftSide != nil && rightSide != nil) {
		logError(a.t, fmt.Sprintln("Fail:", leftSide, "did not equal", rightSide))

		if leftSide != nil && rightSide != nil {
			typeLeft := reflect.TypeOf(leftSide).String()
			typeRight := reflect.TypeOf(rightSide).String()

			if typeLeft != typeRight {
				logError(a.t, "Type of leftside: "+typeLeft)
				logError(a.t, "Type of rightside: "+typeRight)
			}
			//} else {
			//	a.t.Log("Type of leftside: " + reflect.TypeOf(leftSide).String())
			//	a.t.Log("Type of rightside: " + reflect.TypeOf(rightSide).String())
		}
	}

	if a.t.Failed() {
		printLineAndFunction(a.t)
		a.t.FailNow()
	}
}

func Assert(t *testing.T, leftSide interface{}) *AssertObj {
	valueLeft := reflect.ValueOf(leftSide)
	if leftSide != nil && valueLeft.Kind() == reflect.Bool {
		if !valueLeft.Bool() {
			printLineAndFunction(t)
			logError(t, "Assert value was not true")
			t.FailNow()
		}
		return nil
	}

	assertObj := new(AssertObj)
	assertObj.t = t
	assertObj.leftSide = leftSide
	return assertObj
}
