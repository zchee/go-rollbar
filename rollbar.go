// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

const (
	// Name name of client package.
	Name = "go-rollbar"
	// Version version of client package.
	Version = "0.0.0"

	language = "go"
)

// Client represents a first Client methods.
type Client interface {
	Debug(error) Call
	Info(error) Call
	Error(error) Call
	Warn(error) Call
	Critical(error) Call
}
