// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Rollbar API Data Format Reference:
//  https://rollbar.com/docs/api/items_post/#data-format

package rollbar_v1

type Payload struct {
	AccessToken string `json:"access_token"`
	Data        *Data  `json:"data"`
}

type Data struct {
	Environment string                 `json:"environment"`
	Body        *Body                  `json:"body"`
	Level       string                 `json:"level,omitempty"`
	Timestamp   int64                  `json:"timestamp,omitempty"`
	CodeVersion string                 `json:"code_version,omitempty"`
	Platform    string                 `json:"platform,omitempty"`
	Language    string                 `json:"language,omitempty"`
	Framework   string                 `json:"framework,omitempty"`
	Context     string                 `json:"context,omitempty"`
	Request     *Request               `json:"request,omitempty"`
	Person      *Person                `json:"person,omitempty"`
	Server      *Server                `json:"server,omitempty"`
	Client      *Client                `json:"client,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
	Fingerprint string                 `json:"fingerprint,omitempty"`
	Title       string                 `json:"title,omitempty"`
	UUID        string                 `json:"uuid,omitempty"`
	Notifier    *Notifier              `json:"notifier,omitempty"`
}

type Body struct {
	Telemetry   *Telemetry    `json:"telemetry,omitempty"`
	Trace       *Trace        `json:"trace,omitempty"`
	TraceChain  []interface{} `json:"trace_chain,omitempty"`
	Message     *Message      `json:"message,omitempty"`
	CrashReport *CrashReport  `json:"crash_report,omitempty"`
}

type Telemetry struct {
	Level       string        `json:"level"`
	Type        string        `json:"type"`
	Source      string        `json:"source"`
	TimestampMs int           `json:"timestamp_ms"`
	Body        TelemetryBody `json:"body"`
}

type TelemetryBody struct {
	EndTimestampMs   int    `json:"end_timestamp_ms"`
	Method           string `json:"method"`
	StartTimestampMs int    `json:"start_timestamp_ms"`
	StatusCode       string `json:"status_code"`
	Subtype          string `json:"subtype"`
	URL              string `json:"url"`
}

type Trace struct {
	Frames    []*Frame   `json:"frames"`
	Exception *Exception `json:"exception"`
}

type Message struct {
	Body string `json:"body"`
}

type CrashReport struct {
	Raw string `json:"raw"`
}

type Exception struct {
	Class       string `json:"class"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

type Frame struct {
	Filename    string   `json:"filename"`
	Lineno      int      `json:"lineno,omitempty"`
	Colno       int      `json:"colno,omitempty"`
	Method      string   `json:"method,omitempty"`
	Code        string   `json:"code,omitempty"`
	ClassName   string   `json:"class_name,omitempty"`
	Context     *Context `json:"context,omitempty"`
	Argspec     []string `json:"argspec,omitempty"`
	Varargspec  string   `json:"varargspec,omitempty"`
	Keywordspec string   `json:"keywordspec,omitempty"`
	Locals      *Locals  `json:"locals,omitempty"`
}

type Context struct {
	Post []interface{} `json:"post"`
	Pre  []string      `json:"pre"`
}

type Locals struct {
	Args    []interface{} `json:"args"`
	Kwargs  *Kwargs       `json:"kwargs"`
	Request string        `json:"request"`
	User    string        `json:"user"`
}

type Kwargs struct {
	Level string `json:"level"`
}

type Request struct {
	GET         []interface{} `json:"GET"`
	POST        []interface{} `json:"POST"`
	Body        string        `json:"body"`
	Headers     *Headers      `json:"headers"`
	Method      string        `json:"method"`
	Params      *Params       `json:"params"`
	QueryString string        `json:"query_string"`
	URL         string        `json:"url"`
	UserIP      string        `json:"user_ip"`
}

type Headers struct {
	Accept  string `json:"Accept"`
	Referer string `json:"Referer"`
}

type Params struct {
	Action     string `json:"action"`
	Controller string `json:"controller"`
}

type Person struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Server struct {
	Branch      string `json:"branch"`
	CodeVersion string `json:"code_version"`
	Host        string `json:"host"`
	Root        string `json:"root"`
	Sha         string `json:"sha"`
}

type Client struct {
	Javascript *Javascript `json:"javascript"`
}

type Javascript struct {
	Browser             string `json:"browser"`
	CodeVersion         string `json:"code_version"`
	GuessUncaughtFrames bool   `json:"guess_uncaught_frames"`
	SourceMapEnabled    bool   `json:"source_map_enabled"`
}

type Notifier struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
