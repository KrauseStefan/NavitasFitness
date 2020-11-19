package TestHelper

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"

	"golang.org/x/net/context"
)

func GetContext() context.Context {
	ctx := context.Background()
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

	valueLeft := reflect.ValueOf(leftSide)
	if valueLeft.Kind() == reflect.Bool {
		printLineAndFunction(a.t)
		logError(a.t, "Do not use .Equals with boolean expressions")
		a.t.FailNow()
	}

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
		}
	}

	if a.t.Failed() {
		printLineAndFunction(a.t)
		a.t.FailNow()
	}
}

func (a *AssertObj) Contains(rightSide interface{}) {
	leftSide := a.leftSide

	valueLeft := reflect.ValueOf(leftSide)
	if valueLeft.Kind() != reflect.Slice && valueLeft.Kind() != reflect.Array {
		printLineAndFunction(a.t)
		logError(a.t, ".Contains can only be used with slices")
		a.t.FailNow()
		return
	}

	for i := 0; i < valueLeft.Len(); i++ {
		elem := valueLeft.Index(i).Interface()

		// nil does not equal nil (nil == nil) => false
		if reflect.DeepEqual(elem, rightSide) || (elem == nil && rightSide == nil) {
			return
		}
	}

	logError(a.t, fmt.Sprintln("Fail:", leftSide, "did not contain", rightSide))

	if leftSide != nil && rightSide != nil {
		typeLeft := reflect.TypeOf(leftSide).String()
		typeRight := reflect.TypeOf(rightSide).String()

		if typeLeft != typeRight {
			logError(a.t, "Type of array/slice: "+typeLeft)
			logError(a.t, "Type of rightside: "+typeRight)
		}
	}

	printLineAndFunction(a.t)
	a.t.FailNow()
}

func Assert(t *testing.T, leftSide interface{}) *AssertObj {
	valueLeft := reflect.ValueOf(leftSide)
	if leftSide != nil && valueLeft.Kind() == reflect.Bool {
		if !valueLeft.Bool() {
			printLineAndFunction(t)
			logError(t, "Assert value was not true")
			t.FailNow()
		}
	}

	assertObj := new(AssertObj)
	assertObj.t = t
	assertObj.leftSide = leftSide
	return assertObj
}
