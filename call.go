// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	api "github.com/zchee/go-rollbar/api/v1"
	"golang.org/x/net/context"
)

type Call interface {
	Custom(map[string]interface{}) Call
	UUID(string) Call
	Do(context.Context) (*api.Response, error)
}

type callOption struct {
	err    error
	custom map[string]interface{}
	id     string
}

func joinPayload(payload *api.Payload, opt callOption) {
	if opt.custom != nil {
		payload.Data.Custom = opt.custom
	}
	if opt.id != "" {
		payload.Data.UUID = opt.id
	}
}

type DebugCall struct {
	client *httpClient
	callOption
}

// Debug sends the error to rollbar with debug level.
func (c *client) Debug(err error) Call {
	var call DebugCall
	call.client = c.debugClient
	call.err = err
	return &call
}

// Custom is any arbitrary metadata you want to send. "custom" itself should be an object.
func (c *DebugCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

// UUID a string, up to 36 characters, that uniquely identifies this occurrence.
// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
// We recommend using a UUID4 (16 random bytes).
// The UUID space is unique to each project, and can be used to look up an occurrence later.
// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
// one will be discarded.
// While optional, it is recommended that all clients generate and provide this field.
func (c *DebugCall) UUID(id string) Call {
	c.id = id
	return c
}

func (c *DebugCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(DebugLevel, c.err)
	joinPayload(payload, c.callOption)
	return c.client.post(ctx, payload)
}

type InfoCall struct {
	client *httpClient
	callOption
}

// Info sends the error to rollbar with info level.
func (c *client) Info(err error) Call {
	var call InfoCall
	call.client = c.infoClient
	call.err = err
	return &call
}

// Custom is any arbitrary metadata you want to send. "custom" itself should be an object.
func (c *InfoCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

// UUID a string, up to 36 characters, that uniquely identifies this occurrence.
// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
// We recommend using a UUID4 (16 random bytes).
// The UUID space is unique to each project, and can be used to look up an occurrence later.
// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
// one will be discarded.
// While optional, it is recommended that all clients generate and provide this field.
func (c *InfoCall) UUID(id string) Call {
	c.id = id
	return c
}

func (c *InfoCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(InfoLevel, c.err)
	joinPayload(payload, c.callOption)
	return c.client.post(ctx, payload)
}

type ErrorCall struct {
	client *httpClient
	callOption
}

// Error sends the error to rollbar with error level.
func (c *client) Error(err error) Call {
	var call ErrorCall
	call.client = c.errorClient
	call.err = err
	return &call
}

// Custom is any arbitrary metadata you want to send. "custom" itself should be an object.
func (c *ErrorCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

// UUID a string, up to 36 characters, that uniquely identifies this occurrence.
// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
// We recommend using a UUID4 (16 random bytes).
// The UUID space is unique to each project, and can be used to look up an occurrence later.
// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
// one will be discarded.
// While optional, it is recommended that all clients generate and provide this field.
func (c *ErrorCall) UUID(id string) Call {
	c.id = id
	return c
}

func (c *ErrorCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(ErrorLevel, c.err)
	joinPayload(payload, c.callOption)
	return c.client.post(ctx, payload)
}

type WarnCall struct {
	client *httpClient
	callOption
}

// Warn sends the error to rollbar with warning level.
func (c *client) Warn(err error) Call {
	var call WarnCall
	call.client = c.warnClient
	call.err = err
	return &call
}

// Custom is any arbitrary metadata you want to send. "custom" itself should be an object.
func (c *WarnCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

// UUID a string, up to 36 characters, that uniquely identifies this occurrence.
// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
// We recommend using a UUID4 (16 random bytes).
// The UUID space is unique to each project, and can be used to look up an occurrence later.
// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
// one will be discarded.
// While optional, it is recommended that all clients generate and provide this field.
func (c *WarnCall) UUID(id string) Call {
	c.id = id
	return c
}

func (c *WarnCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(WarnLevel, c.err)
	joinPayload(payload, c.callOption)
	return c.client.post(ctx, payload)
}

type CriticalCall struct {
	client *httpClient
	callOption
}

// Critical sends the error to rollbar with critical level.
func (c *client) Critical(err error) Call {
	var call CriticalCall
	call.client = c.criticalClient
	call.err = err
	return &call
}

// Custom is any arbitrary metadata you want to send. "custom" itself should be an object.
func (c *CriticalCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

// UUID a string, up to 36 characters, that uniquely identifies this occurrence.
// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
// We recommend using a UUID4 (16 random bytes).
// The UUID space is unique to each project, and can be used to look up an occurrence later.
// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
// one will be discarded.
// While optional, it is recommended that all clients generate and provide this field.
func (c *CriticalCall) UUID(id string) Call {
	c.id = id
	return c
}

func (c *CriticalCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(CriticalLevel, c.err)
	joinPayload(payload, c.callOption)
	return c.client.post(ctx, payload)
}
