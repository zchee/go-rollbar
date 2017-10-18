// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"golang.org/x/net/context"
)

type Call interface {
	Custom(map[string]interface{}) Call
	UUID(string) Call
	Do(context.Context) error
}

type callOption struct {
	err    error
	custom map[string]interface{}
	id     string
}

type DebugService struct {
	client *httpClient
}

// Debug sends the error to rollbar with debug level.
func (c *Client) Debug(err error) Call {
	var call DebugCall
	call.service = c.debug
	call.err = err
	return &call
}

type InfoService struct {
	client *httpClient
}

// Info sends the error to rollbar with info level.
func (c *Client) Info(err error) Call {
	var call InfoCall
	call.service = c.info
	call.err = err
	return &call
}

type ErrorService struct {
	client *httpClient
}

// Error sends the error to rollbar with error level.
func (c *Client) Error(err error) Call {
	var call ErrorCall
	call.service = c.errorService
	call.err = err
	return &call
}

type WarnService struct {
	client *httpClient
}

// Warn sends the error to rollbar with warning level.
func (c *Client) Warn(err error) Call {
	var call WarnCall
	call.service = c.warn
	call.err = err
	return &call
}

type CriticalService struct {
	client *httpClient
}

// Critical sends the error to rollbar with critical level.
func (c *Client) Critical(err error) Call {
	var call CriticalCall
	call.service = c.critical
	call.err = err
	return &call
}

type DebugCall struct {
	service *DebugService
	callOption
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

func (c *DebugCall) Do(ctx context.Context) error {
	payload := c.service.client.payload(DebugLevel, c.err)
	payload.Data.Custom = c.custom
	payload.Data.UUID = c.id
	return c.service.client.post(ctx, payload)
}

type InfoCall struct {
	service *InfoService
	callOption
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

func (c *InfoCall) Do(ctx context.Context) error {
	payload := c.service.client.payload(InfoLevel, c.err)
	payload.Data.Custom = c.custom
	payload.Data.UUID = c.id
	return c.service.client.post(ctx, payload)
}

type ErrorCall struct {
	service *ErrorService
	callOption
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

func (c *ErrorCall) Do(ctx context.Context) error {
	payload := c.service.client.payload(ErrorLevel, c.err)
	payload.Data.Custom = c.custom
	payload.Data.UUID = c.id
	return c.service.client.post(ctx, payload)
}

type WarnCall struct {
	service *WarnService
	callOption
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

func (c *WarnCall) Do(ctx context.Context) error {
	payload := c.service.client.payload(WarnLevel, c.err)
	payload.Data.Custom = c.custom
	payload.Data.UUID = c.id
	return c.service.client.post(ctx, payload)
}

type CriticalCall struct {
	service *CriticalService
	callOption
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

func (c *CriticalCall) Do(ctx context.Context) error {
	payload := c.service.client.payload(CriticalLevel, c.err)
	payload.Data.Custom = c.custom
	payload.Data.UUID = c.id
	return c.service.client.post(ctx, payload)
}
