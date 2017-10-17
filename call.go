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

func (c *DebugCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

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

func (c *InfoCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

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

func (c *ErrorCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

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

func (c *WarnCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

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

func (c *CriticalCall) Custom(custom map[string]interface{}) Call {
	c.custom = custom
	return c
}

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
