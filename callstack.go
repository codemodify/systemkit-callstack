package callstack

import (
	"encoding/json"
	"runtime"
	"strings"
)

type Frame struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"function"`
}

func (thisRef Frame) String() string {
	data, _ := json.Marshal(thisRef)
	return string(data)
}

// Get - return current call frame
func Get() Frame {
	callStack := GetFramesWithSkip(4)

	if len(callStack) == 0 {
		return Frame{}
	}

	return callStack[0]
}

// GetFrames - return call stack, no skips
func GetFrames() []Frame {
	return GetFramesWithSkip(4)
}

// GetFramesWithSkip - return call stack, skip N frames
func GetFramesWithSkip(skip int) []Frame {
	return NativeFramesToFrames(GetNativeFrames(skip))
}

// GetNativeFrames - return Go native call stack
func GetNativeFrames(skip int) []uintptr {
	const callStackDepth = 50 // most relevant context seem to appear near the top of the stack
	var callStackBuffer = make([]uintptr, callStackDepth)
	callStackSize := runtime.Callers(skip, callStackBuffer)
	return callStackBuffer[:callStackSize]
}

// NativeFramesToFrames - converts native call stack to []Frame
func NativeFramesToFrames(callStack []uintptr) []Frame {
	frames := []Frame{}

	callStackFrames := runtime.CallersFrames(callStack)
	for {
		frame, ok := callStackFrames.Next()
		if !ok {
			break
		}

		pkg, fn := splitPackageFuncName(frame.Function)
		if frameFilter(pkg, fn, frame.File, frame.Line) {
			frames = frames[:0]
			continue
		}

		frames = append(frames, Frame{
			File: frame.File,
			Line: frame.Line,
			Func: fn,
		})
	}

	return frames
}

func splitPackageFuncName(funcName string) (string, string) {
	var packageName string
	if ind := strings.LastIndex(funcName, "/"); ind > 0 {
		packageName += funcName[:ind+1]
		funcName = funcName[ind+1:]
	}
	if ind := strings.Index(funcName, "."); ind > 0 {
		packageName += funcName[:ind]
		funcName = funcName[ind+1:]
	}
	return packageName, funcName
}

func frameFilter(packageName, funcName string, file string, line int) bool {
	return packageName == "runtime" && funcName == "panic"
}
