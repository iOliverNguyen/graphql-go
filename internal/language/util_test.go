package language

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func deepEqual(T *testing.T, A, B interface{}) {
	if !reflect.DeepEqual(A, B) {
		T.Errorf("Expect deep equal `%#v` `%#v`", A, B)
	}
}

func expect(T *testing.T, condition bool, msg string, args ...interface{}) {
	if !condition {
		T.Errorf(msg, args...)
	}
}

func expectPanic(T *testing.T, fn func(), msg string) {
	defer func() {
		err := recover()
		if err == nil {
			T.Errorf("Expect panic with message:\n%v\n---")
			return
		}
		errMsg := fmt.Sprint(err)
		if !strings.Contains(errMsg, msg) {
			T.Errorf("Expect panic with message:\n%v---\nbut got:\n%v---", msg, errMsg)
		}
	}()
	fn()
}
