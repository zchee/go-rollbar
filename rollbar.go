// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"golang.org/x/net/context"
)

const (
	Name    = "go-rollbar"
	Version = "0.0.0"

	language = "go"
)

// Debug sends the error to rollbar with debug level.
func (c *Client) Debug(ctx context.Context, err error) {
	if err := c.post(ctx, DebugLevel, err); err != nil {
		c.logger.Fatal(ctx, "rollbar: Debug: %+v", err)
	}
}

// Info sends the error to rollbar with info level.
func (c *Client) Info(ctx context.Context, err error) {
	if err := c.post(ctx, InfoLevel, err); err != nil {
		c.logger.Fatal(ctx, "rollbar: Info: %+v", err)
	}
}

// Error sends the error to rollbar with error level.
func (c *Client) Error(ctx context.Context, err error) {
	if err := c.post(ctx, ErrorLevel, err); err != nil {
		c.logger.Fatal(ctx, "rollbar: Error: %+v", err)
	}
}

// Warn sends the error to rollbar with warn level.
func (c *Client) Warn(ctx context.Context, err error) {
	if err := c.post(ctx, WarnLevel, err); err != nil {
		c.logger.Fatal(ctx, "rollbar: Warn: %+v", err)
	}
}

// Critical sends the error to rollbar with critical level.
func (c *Client) Critical(ctx context.Context, err error) {
	if err := c.post(ctx, CriticalLevel, err); err != nil {
		c.logger.Fatal(ctx, "rollbar: Critical: %+v", err)
	}
}
