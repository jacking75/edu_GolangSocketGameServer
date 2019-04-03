package gohipernetFake

import (
	"fmt"
	"runtime"

	"github.com/davecgh/go-spew/spew"
)

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		Logger.Error(fmt.Sprintf("%v", x))

		i := 0
		funcName, file, line, ok := runtime.Caller(i)

		for ok {
			msg := fmt.Sprintf("PrintPanicStack. [func]: %s, [file]: %s, [line]: %d\n", runtime.FuncForPC(funcName).Name(), file, line)
			Logger.Error(msg)
			IExportLog("Error", msg)
			//Logger.Error("PrintPanicStack", zap.Int("N", i), zap.String("func",runtime.FuncForPC(funcName).Name()), zap.String("file", file), zap.Int("line", line))
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			msg := fmt.Sprintf("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
			Logger.Error(msg)
			IExportLog("Error", msg)
		}
	}
}
