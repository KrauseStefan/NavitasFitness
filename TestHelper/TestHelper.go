package TestHelper

import (
	"appengine_internal"

	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"
)

type Spy struct {
	argument_1 []interface{}
	argument_2 []interface{}
	argument_3 []interface{}
	callCount  uint
}

func (s *Spy) RegisterCall() {
	s.callCount += 1
}

func (s *Spy) CallCount() int {
	return int(s.callCount)
}

func prependToArgumentSlice(argSlice []interface{}, item interface{}) []interface{} {
	if len(argSlice) == 0 {
		argSlice = make([]interface{}, 0, 10)
	}

	return append([]interface{}{item}, argSlice...)
}

func (s *Spy) RegisterArg1(arg1 interface{}) {
	s.argument_1 = prependToArgumentSlice(s.argument_1, arg1)
}

func (s *Spy) RegisterArg2(arg1 interface{}, arg2 interface{}) {
	s.RegisterArg1(arg1)
	s.argument_2 = prependToArgumentSlice(s.argument_2, arg2)
}

func (s *Spy) RegisterArg3(arg1 interface{}, arg2 interface{}, arg3 interface{}) {
	s.RegisterArg2(arg1, arg2)
	s.argument_3 = prependToArgumentSlice(s.argument_3, arg3)
}

func (s *Spy) GetLatestArg1() interface{} {
	if len(s.argument_1) > 0 {
		return s.argument_1[0]
	}
	return nil
}

type ContextMock struct {
	OptionalId int
	//req  *http.Request
	//done chan struct{} // Closed when the context has expired.
}

func (c *ContextMock) Debugf(format string, args ...interface{})    {}
func (c *ContextMock) Infof(format string, args ...interface{})     {}
func (c *ContextMock) Warningf(format string, args ...interface{})  {}
func (c *ContextMock) Errorf(format string, args ...interface{})    {}
func (c *ContextMock) Criticalf(format string, args ...interface{}) {}
func (c *ContextMock) Call(service, method string, in, out appengine_internal.ProtoMessage, opts *appengine_internal.CallOptions) error {
	return nil
}
func (c *ContextMock) FullyQualifiedAppID() string {
	return ""
}
func (c *ContextMock) Request() interface{} {
	return nil
}

//func assert(t *testing.T, result bool) {
//	if !result {
//		t.Fail()
//	}
//}

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
	if !reflect.DeepEqual(leftSide, rightSide ) && (leftSide != nil && rightSide != nil) {
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

func LogObject(obj interface{}, message string) {
	json, _ := json.Marshal(obj)
	fmt.Println(message, string(json))
}
