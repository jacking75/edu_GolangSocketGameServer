package gohipernetFake

import (
	"fmt"
	"runtime"

	"github.com/davecgh/go-spew/spew"
)

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		logError("", 0, fmt.Sprintf("%v", x))

		i := 0
		funcName, file, line, ok := runtime.Caller(i)

		for ok {
			msg := fmt.Sprintf("PrintPanicStack. [func]: %s, [file]: %s, [line]: %d\n", runtime.FuncForPC(funcName).Name(), file, line)
			logError("", 0, msg)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			msg := fmt.Sprintf("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
			logError("", 0, msg)
		}
	}
}
