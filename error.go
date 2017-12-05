// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	api "github.com/zchee/go-rollbar/api/v1"
)

// errorBody creates a rollbar error body with a given stack trace.
func errorBody(err error, stack Stack) *api.Body {
	message := "<nil>"
	if err != nil {
		message = err.Error()
	}

	return &api.Body{
		Trace: &api.Trace{
			Frames: stack,
			Exception: &api.Exception{
				Class:   errorClass(err),
				Message: message,
			},
		},
	}
}

// errorClass expands the function(class) name from err.
func errorClass(err error) string {
	if err == nil {
		return "<nil>"
	}

	fn := reflect.TypeOf(err).String()
	switch fn {
	case "":
		return "panic"
	case "*errors.errorString":
		checksum := adler32.Checksum([]byte(err.Error()))
		return fmt.Sprintf("{%x}", checksum)
	default:
		return strings.TrimPrefix(fn, "*")
	}
}

var (
	// TODO(zchee): remove it
	reFilterHeaders = regexp.MustCompile("Authorization")
	reFilterFields  = regexp.MustCompile("password|secret|token")
)

func errorRequest(req *http.Request) *api.Request {
	const remoteIP = "$remote_ip"

	query := filterParams(reFilterFields, req.URL.Query())
	body, _ := ioutil.ReadAll(req.Body)

	return &api.Request{
		URL:         req.URL.String(),
		Method:      req.Method,
		Headers:     filterParams(reFilterHeaders, req.Header),
		GET:         query,
		QueryString: url.Values(query).Encode(),
		POST:        req.Form,
		Body:        string(body),
		UserIP:      remoteIP,
	}
}

func filterParams(pat *regexp.Regexp, values map[string][]string) map[string][]string {
	const redacted = "xxxxxxxxxxxx (redacted)"
	for key := range values {
		if pat.MatchString(key) {
			values[key] = []string{redacted}
		}
	}

	return values
}
