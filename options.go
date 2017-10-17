// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"net/http"
)

// Option defines an interface of optional parameters to the
// `rollbar.New` constructor.
type Option interface {
	Key() key
	Value() interface{}
}

type option struct {
	key   key
	value interface{}
}

func (o *option) Key() key {
	return o.key
}
func (o *option) Value() interface{} {
	return o.value
}

type key int

const (
	keyHTTPClient key = 1 << iota
	keyEndpoint
	keyLogger
	keyDebug
	keyEnvironment
	keyPlatform
	keyCodeVersion
	keyServerHost
	keyServerRoot
)

// WithClient allows you to specify an net/http.Client object to
// use to communicate with the Rollbar endpoints.
//
// For example, if you need to use this in Google App Engine, you can pass it the
// result of `urlfetch.Client`.
func WithClient(cl *http.Client) Option {
	return &option{
		key:   keyHTTPClient,
		value: cl,
	}
}

// WithEndpoint allows you to specify an alternate API endpoint.
// The default is DefaultEndpoint.
func WithEndpoint(s string) Option {
	return &option{
		key:   keyEndpoint,
		value: s,
	}
}

// WithLogger specifies the logger object to be used.
// If not specified and `WithDebug` is enabled, then a default
// logger which writes to os.Stderr.
func WithLogger(l Logger) Option {
	return &option{
		key:   keyLogger,
		value: l,
	}
}

// WithDebug specifies that we want to run in debugging mode.
// You can set this value manually to override any existing global
// defaults.
//
// If one is not specified, the default value is false, or the
// value specified in ROLLBAR_DEBUG environment variable.
func WithDebug(b bool) Option {
	return &option{
		key:   keyDebug,
		value: b,
	}
}

// WithEnvironment name of the environment in which this occurrence was seen.
//
// A string up to 255 characters. For best results, use "production" or "prod" for your
// production environment.
// You don't need to configure anything in the Rollbar UI for new environment names;
// we'll detect them automatically.
func WithEnvironment(env string) Option {
	return &option{
		key:   keyEnvironment,
		value: env,
	}
}

// WithPlatform name of platform on which this occurred.
//
// Meaningful platform names:
//  "browser", "android", "ios", "flash", "client", "heroku", "google-app-engine"
// If this is a client-side event, be sure to specify the platform and use a post_client_item access token.
func WithPlatform(value string) Option {
	return &option{
		key:   keyPlatform,
		value: value,
	}
}

// WithCodeVersion is a string, up to 40 characters, describing the version of the application code
//
// Rollbar understands these formats:
//  - semantic version (i.e. "2.1.12")
//  - integer (i.e. "45")
//  - git SHA (i.e. "3da541559918a808c2402bba5012f6c60b27661c")
func WithCodeVersion(version string) Option {
	return &option{
		key:   keyCodeVersion,
		value: version,
	}
}

// WithServerHost is the server hostname. Will be indexed.
func WithServerHost(hostname string) Option {
	return &option{
		key:   keyServerHost,
		value: hostname,
	}
}

// WithServerRoot is the path to the application code root. Not including the final slash.
// Used to collapse non-project code when displaying tracebacks.
func WithServerRoot(root string) Option {
	return &option{
		key:   keyServerRoot,
		value: root,
	}
}

type ErrorOption interface {
	Key() errorKey
	Value() interface{}
}

type errorOption struct {
	key   errorKey
	value interface{}
}

func (o *errorOption) Key() errorKey {
	return o.key
}
func (o *errorOption) Value() interface{} {
	return o.value
}

type errorKey int

const (
	keyCustom errorKey = 1 << iota
	keyUUID
)

// WithCustom is any arbitrary metadata you want to send. "custom" itself should be an object.
func WithCustom(version string) ErrorOption {
	return &errorOption{
		key:   keyCustom,
		value: version,
	}
}

// WithUUID a string, up to 36 characters, that uniquely identifies this occurrence.
// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
// We recommend using a UUID4 (16 random bytes).
// The UUID space is unique to each project, and can be used to look up an occurrence later.
// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
// one will be discarded.
// While optional, it is recommended that all clients generate and provide this field.
func WithUUID(uuid string) ErrorOption {
	return &errorOption{
		key:   keyUUID,
		value: uuid,
	}
}
