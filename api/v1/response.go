// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v1

// Response represents a rollbar REST API response.
type Response struct {
	Err     int    `json:"err"`
	Result  Result `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
}

// Result result of response data.
type Result struct {
	UUID string `json:"uuid"`
}
