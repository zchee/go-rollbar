// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"net/http"

	api "github.com/zchee/go-rollbar/api/v1"
	"golang.org/x/net/context"
)

// Call respesents a fluent style call options.
type Call interface {
	Request(*http.Request) Call
	Person(string, string, string) Call
	Custom(map[string]interface{}) Call
	UUID(string) Call
	Do(context.Context) (*api.Response, error)
}

type callOption struct {
	err    error
	req    *http.Request
	person *api.Person
	custom map[string]interface{}
	id     string
}

func joinPayload(payload *api.Payload, opt callOption) {
	if opt.req != nil {
		payload.Data.Request = errorRequest(opt.req)
	}
	if opt.person != nil {
		payload.Data.Person = opt.person
	}
	if opt.custom != nil {
		payload.Data.Custom = opt.custom
	}
	if opt.id != "" {
		payload.Data.UUID = opt.id
	}
}

// DebugCall represents a calls the debug level stack trace.
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

// Request is a data about the request this event occurred in.
func (c *DebugCall) Request(req *http.Request) Call {
	c.req = req
	return c
}

// Person is the user affected by this event. Will be indexed by ID, username, and email.
// People are stored in Rollbar keyed by ID. If you send a multiple different usernames/emails for the
// same ID, the last received values will overwrite earlier ones.
func (c *DebugCall) Person(id, username, email string) Call {
	if id == "" { // id is required
		return c
	}

	c.person = &api.Person{
		ID:       id,
		Username: username,
		Email:    email,
	}
	return c
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

// Do executes the call to access rollbar endpoint.
func (c *DebugCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(DebugLevel, c.err)
	joinPayload(payload, c.callOption)
	req, err := c.client.newRequest(payload)
	if err != nil {
		return nil, err
	}
	var m api.Response
	err = c.client.Do(ctx, req, &m)
	return &m, err
}

// InfoCall represents a calls the info level stack trace.
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

// Request is a data about the request this event occurred in.
func (c *InfoCall) Request(req *http.Request) Call {
	c.req = req
	return c
}

// Person is the user affected by this event. Will be indexed by ID, username, and email.
// People are stored in Rollbar keyed by ID. If you send a multiple different usernames/emails for the
// same ID, the last received values will overwrite earlier ones.
func (c *InfoCall) Person(id, username, email string) Call {
	if id == "" { // id is required
		return c
	}

	c.person = &api.Person{
		ID:       id,
		Username: username,
		Email:    email,
	}
	return c
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

// Do executes the call to access rollbar endpoint.
func (c *InfoCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(InfoLevel, c.err)
	joinPayload(payload, c.callOption)
	req, err := c.client.newRequest(payload)
	if err != nil {
		return nil, err
	}
	var m api.Response
	err = c.client.Do(ctx, req, &m)
	return &m, err
}

// ErrorCall represents a calls the error level stack trace.
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

// Request is a data about the request this event occurred in.
func (c *ErrorCall) Request(req *http.Request) Call {
	c.req = req
	return c
}

// Person is the user affected by this event. Will be indexed by ID, username, and email.
// People are stored in Rollbar keyed by ID. If you send a multiple different usernames/emails for the
// same ID, the last received values will overwrite earlier ones.
func (c *ErrorCall) Person(id, username, email string) Call {
	if id == "" { // id is required
		return c
	}

	c.person = &api.Person{
		ID:       id,
		Username: username,
		Email:    email,
	}
	return c
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

// Do executes the call to access rollbar endpoint.
func (c *ErrorCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(ErrorLevel, c.err)
	joinPayload(payload, c.callOption)
	req, err := c.client.newRequest(payload)
	if err != nil {
		return nil, err
	}
	var m api.Response
	err = c.client.Do(ctx, req, &m)
	return &m, err
}

// WarnCall represents a calls the warning level stack trace.
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

// Request is a data about the request this event occurred in.
func (c *WarnCall) Request(req *http.Request) Call {
	c.req = req
	return c
}

// Person is the user affected by this event. Will be indexed by ID, username, and email.
// People are stored in Rollbar keyed by ID. If you send a multiple different usernames/emails for the
// same ID, the last received values will overwrite earlier ones.
func (c *WarnCall) Person(id, username, email string) Call {
	if id == "" { // id is required
		return c
	}

	c.person = &api.Person{
		ID:       id,
		Username: username,
		Email:    email,
	}
	return c
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

// Do executes the call to access rollbar endpoint.
func (c *WarnCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(WarnLevel, c.err)
	joinPayload(payload, c.callOption)
	req, err := c.client.newRequest(payload)
	if err != nil {
		return nil, err
	}
	var m api.Response
	err = c.client.Do(ctx, req, &m)
	return &m, err
}

// CriticalCall represents a calls the critical level stack trace.
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

// Request is a data about the request this event occurred in.
func (c *CriticalCall) Request(req *http.Request) Call {
	c.req = req
	return c
}

// Person is the user affected by this event. Will be indexed by ID, username, and email.
// People are stored in Rollbar keyed by ID. If you send a multiple different usernames/emails for the
// same ID, the last received values will overwrite earlier ones.
func (c *CriticalCall) Person(id, username, email string) Call {
	if id == "" { // id is required
		return c
	}

	c.person = &api.Person{
		ID:       id,
		Username: username,
		Email:    email,
	}
	return c
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

// Do executes the call to access rollbar endpoint.
func (c *CriticalCall) Do(ctx context.Context) (*api.Response, error) {
	payload := c.client.payload(CriticalLevel, c.err)
	joinPayload(payload, c.callOption)
	req, err := c.client.newRequest(payload)
	if err != nil {
		return nil, err
	}
	var m api.Response
	err = c.client.Do(ctx, req, &m)
	return &m, err
}
