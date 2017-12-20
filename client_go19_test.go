// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

// This stack trace data is Go 1.9.1.

package rollbar

import (
	api "github.com/zchee/go-rollbar/api/v1"
)

var (
	testTrace = &api.Trace{
		Frames: []*api.Frame{
			{
				Filename: "/usr/local/go/src/testing/testing.go",
				Lineno:   746,
				Method:   "testing.tRunner",
			},
			{
				Filename: "/usr/local/go/src/runtime/asm_amd64.s",
				Lineno:   2337,
				Method:   "runtime.goexit",
			},
		},
		Exception: &api.Exception{
			Class:   "{23d90530}",
			Message: "default error",
		},
	}
	testFingerprint = "485dc810"
)
