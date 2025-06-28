package xruntime

import (
	"fmt"
	"runtime"
	"strings"
)

// StackTrace returns a slice of strings, each of which is a text representation of a frame in the stack trace. 'skip'
// is the number of stack frames to skip before processing, with 0 identifying the caller of StackTrace.
func StackTrace(skip int) []string {
	var pcs [128]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	return PCsToStackTrace(pcs[:n])
}

// PCsToStackTrace converts a slice of program counters (PCs) into a slice of strings representing the stack trace.
func PCsToStackTrace(pcs []uintptr) []string {
	stack := make([]string, 0, len(pcs))
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		if !more {
			break
		}
		stack = append(stack, fmt.Sprintf("[%s] %s:%d", f.Function, StackTracePath(f.File), f.Line))
	}
	return stack
}

// StackTracePath returns the last directory and the file name component of the path. Generally, this allows the stack
// trace entries to contain enough of a path to be useful and to allow the IDE to find and open the file, but doesn't
// include long full paths.
func StackTracePath(path string) string {
	if i := strings.LastIndexByte(path, '/'); i != -1 {
		if i = strings.LastIndexByte(path[:i], '/'); i != -1 {
			return path[i+1:]
		}
	}
	return path
}
