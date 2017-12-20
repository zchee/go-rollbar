// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"fmt"
	"hash/crc32"
	"path/filepath"
	"runtime"
	"strings"

	api "github.com/zchee/go-rollbar/api/v1"
)

// Stack represents a api.Frame slice.
type Stack []*api.Frame

// CreateStack creates the Stack data except before skip callers.
func CreateStack(skip int) Stack {
	var stack Stack

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = append(stack, &api.Frame{
			Filename: file,
			Method:   funcName(pc),
			Lineno:   line,
		})
	}

	return stack
}

// CreateStackFromCaller creates the Stack data from callers.
func CreateStackFromCaller(callers []uintptr) Stack {
	stack := make(Stack, 0, len(callers))

	for _, caller := range callers {
		if fn := runtime.FuncForPC(caller); fn != nil {
			file, line := fn.FileLine(caller)
			stack = append(stack, &api.Frame{
				Filename: file,
				Method:   funcNameFromFunc(fn),
				Lineno:   line,
			})
		}
	}

	return stack
}

// Fingerprint create a fingerprint that uniqely identify a given message.
// We use the full callstack, including file names. That ensure that there are no false duplicates
// but also means that after changing the code (adding/removing lines), the fingerprints will change.
// It's a trade-off.
func (s Stack) Fingerprint() string {
	h := crc32.NewIEEE()
	for _, frame := range s {
		fmt.Fprintf(h, "%s%s%d", frame.Filename, frame.Method, frame.Lineno)
	}
	return fmt.Sprintf("%x", h.Sum32())
}

func funcNameFromFunc(fn *runtime.Func) string {
	if fn == nil {
		return "???"
	}
	name := fn.Name()
	end := strings.LastIndex(name, string(filepath.Separator))
	return name[end+1:]
}

func funcName(pc uintptr) string {
	return funcNameFromFunc(runtime.FuncForPC(pc))
}
