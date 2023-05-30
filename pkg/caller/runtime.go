package caller

import "runtime"

func GetRunTimeCaller(skip int) string {
	var callerName string

	if pointerCaller, _, _, success := runtime.Caller(skip); success {
		functionCaller := runtime.FuncForPC(pointerCaller)
		if functionCaller != nil {
			callerName = functionCaller.Name()
		}
	}

	return callerName
}
