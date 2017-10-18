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

type Stack []*api.Frame

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
	return name[end+1 : len(name)]
}

func funcName(pc uintptr) string {
	return funcNameFromFunc(runtime.FuncForPC(pc))
}
