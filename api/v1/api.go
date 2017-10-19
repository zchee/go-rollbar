// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar_v1

const (
	// DefaultEndpoint default of Rollbar v1 API endpoint.
	DefaultEndpoint = "https://api.rollbar.com/api/1/item/"
)

// Payload represents a Rollbar REST API payload.
type Payload struct {
	// AccessToken an access token with scope "post_server_item" or "post_client_item".
	//
	// A post_client_item token must be used if the "platform" is "browser", "android", "ios", "flash", or "client"
	// A post_server_item token should be used for other platforms.
	AccessToken string `json:"access_token"`
	// Data is a payload main data.
	Data *Data `json:"data"`
}

// Data is a payload main data.
type Data struct {
	// Environment is the name of the environment in which this occurrence was seen.
	// A string up to 255 characters. For best results, use "production" or "prod" for your
	// production environment.
	// You don't need to configure anything in the Rollbar UI for new environment names;
	// we'll detect them automatically.
	Environment string `json:"environment"`
	Body        *Body  `json:"body"`
	// Level is the severity level. One of: "critical", "error", "warning", "info", "debug"
	// Defaults to "error" for exceptions and "info" for messages.
	// The level of the *first* occurrence of an item is used as the item's level.
	Level string `json:"level,omitempty"`
	// Timestamp is a when this occurred, as a unix timestamp.
	Timestamp int64 `json:"timestamp,omitempty"`
	// CodeVersion is a string, up to 40 characters, describing the version of the application code
	// Rollbar understands these formats:
	//  - semantic version (i.e. "2.1.12")
	//  - integer (i.e. "45")
	//  - git SHA (i.e. "3da541559918a808c2402bba5012f6c60b27661c")
	// If you have multiple code versions that are relevant, those can be sent inside "client" and "server"
	// (see those sections below)
	// For most cases, just send it here.
	CodeVersion string `json:"code_version,omitempty"`
	// Platform is the platform on which this occurred. Meaningful platform names:
	// "browser", "android", "ios", "flash", "client", "heroku", "google-app-engine"
	// If this is a client-side event, be sure to specify the platform and use a post_client_item access token.
	Platform string `json:"platform,omitempty"`
	// Language is the name of the language your code is written in.
	// This can affect the order of the frames in the stack trace. The following languages set the most
	// recent call first - 'ruby', 'javascript', 'php', 'java', 'objective-c', 'lua'
	// It will also change the way the individual frames are displayed, with what is most consistent with
	// users of the language.
	Language string `json:"language,omitempty"`
	// Framework is the name of the framework your code uses.
	Framework string `json:"framework,omitempty"`
	// Context is an identifier for which part of your application this event came from.
	// Items can be searched by context (prefix search)
	// For example, in a Rails app, this could be `controller#action`.
	// In a single-page javascript app, it could be the name of the current screen or route.
	Context string   `json:"context,omitempty"`
	Request *Request `json:"request,omitempty"`
	Person  *Person  `json:"person,omitempty"`
	Server  *Server  `json:"server,omitempty"`
	Client  *Client  `json:"client,omitempty"`
	// Custom is the any arbitrary metadata you want to send. "custom" itself should be an object.
	Custom map[string]interface{} `json:"custom,omitempty"`
	// Fingerprint is a string controlling how this occurrence should be grouped. Occurrences with the same
	// fingerprint are grouped together. See the "Grouping" guide for more information.
	// Should be a string up to 40 characters long; if longer than 40 characters, we'll use its SHA1 hash.
	// If omitted, we'll determine this on the backend.
	Fingerprint string `json:"fingerprint,omitempty"`
	// Title is a string that will be used as the title of the Item occurrences will be grouped into.
	// Max length 255 characters.
	// If omitted, we'll determine this on the backend.
	Title string `json:"title,omitempty"`
	// UUID is a string, up to 36 characters, that uniquely identifies this occurrence.
	// While it can now be any latin1 string, this may change to be a 16 byte field in the future.
	// We recommend using a UUID4 (16 random bytes).
	// The UUID space is unique to each project, and can be used to look up an occurrence later.
	// It is also used to detect duplicate requests. If you send the same UUID in two payloads, the second
	// one will be discarded.
	// While optional, it is recommended that all clients generate and provide this field.
	UUID     string    `json:"uuid,omitempty"`
	Notifier *Notifier `json:"notifier,omitempty"`
}

// Body is the main data being sent. It can either be a message, an exception, or a crash report.
type Body struct {
	Telemetry *Telemetry `json:"telemetry,omitempty"`
	Trace     *Trace     `json:"trace,omitempty"`
	// TraceChain is the used for exceptions with inner exceptions or causes.
	TraceChain  []*Trace     `json:"trace_chain,omitempty"`
	Message     *Message     `json:"message,omitempty"`
	CrashReport *CrashReport `json:"crash_report,omitempty"`
}

// Telemetry only applicable if you are sending telemetry data.
type Telemetry struct {
	// Level is the severity level of the telemetry data. One of: "critical", "error", "warning", "info", "debug".
	Level string `json:"level"`
	// Type is the type of telemetry data. One of: "log", "network", "dom", "navigation", "error", "manual".
	Type string `json:"type"`
	// Source is the source of the telemetry data. Usually "client" or "server".
	Source string `json:"source"`
	// TimestampMs is the when this occurred, as a unix timestamp in milliseconds.
	TimestampMs int           `json:"timestamp_ms"`
	Body        TelemetryBody `json:"body"`
}

// TelemetryBody is the key-value pairs for the telemetry data point. See "body" key below.
//
// If type above is "log", body should contain "message" key.
//
// If type above is "network", body should contain "method", "url", and "status_code" keys.
//
// If type above is "dom", body should contain "element" key.
//
// If type above is "navigation", body should contain "from" and "to" keys.
//
// If type above is "error", body should contain "message" key.
type TelemetryBody struct {
	EndTimestampMs   int    `json:"end_timestamp_ms"`
	Method           string `json:"method"`
	StartTimestampMs int    `json:"start_timestamp_ms"`
	StatusCode       string `json:"status_code"`
	Subtype          string `json:"subtype"`
	URL              string `json:"url"`
}

// Trace is the stack trace data.
type Trace struct {
	Frames    []*Frame   `json:"frames"`
	Exception *Exception `json:"exception"`
}

// Message only one of "trace", "trace_chain", "message", or "crash_report" should be present.
// Presence of a "message" key means that this payload is a log message.
type Message struct {
	Body string `json:"body"`
}

// CrashReport only one of "trace", "trace_chain", "message", or "crash_report" should be present.
type CrashReport struct {
	Raw string `json:"raw"`
}

// Exception is an object describing the exception instance.
type Exception struct {
	// Class is the exception class name.
	Class string `json:"class"`
	// Description is the exception message, as a string.
	Description string `json:"description"`
	// Message is an alternate human-readable string describing the exception.
	// Usually the original exception message will have been machine-generated;
	// you can use this to send something custom.
	Message string `json:"message"`
}

// Frame is the stack frames.
type Frame struct {
	// Filename is the filename including its full path.
	Filename string `json:"filename"`
	// Lineno is the line number as an integer.
	Lineno int `json:"lineno,omitempty"`
	// Colno is the column number as an integer.
	Colno int `json:"colno,omitempty"`
	// Method is the method or function name.
	Method string `json:"method,omitempty"`
	// Code is the line of code.
	Code string `json:"code,omitempty"`
	// ClassName is a string containing the class name.
	// Used in the UI when the payload's top-level "language" key has the value "java".
	ClassName string   `json:"class_name,omitempty"`
	Context   *Context `json:"context,omitempty"`
	// Argspec is the list of the name of the arguments to the method/function call.
	Argspec []string `json:"argspec,omitempty"`
	// Varargspec is the if the function call takes an arbitrary number of unnamed positional arguments,
	// the name of the argument that is the list containing those arguments.
	// For example, in Python, this would typically be "args" when "*args" is used.
	// The actual list will be found in locals.
	Varargspec string `json:"varargspec,omitempty"`
	// Keywordspec if the function call takes an arbitrary number of keyword arguments, the name
	// of the argument that is the object containing those arguments.
	// For example, in Python, this would typically be "kwargs" when "**kwargs" is used.
	// The actual object will be found in locals.
	Keywordspec string  `json:"keywordspec,omitempty"`
	Locals      *Locals `json:"locals,omitempty"`
}

// Context Additional code before and after the "code" line.
type Context struct {
	// Pre is the list of lines of code before the "code" line
	Pre []string `json:"pre"`
	// Post is the list of line of code after the "code" line.
	Post []string `json:"post"`
}

// Locals is the object of local variables for the method/function call.
// The values of variables from argspec, vararspec and keywordspec
// can be found in locals.
type Locals struct {
	Request string            `json:"request"`
	User    string            `json:"user"`
	Args    []interface{}     `json:"args"`
	Kwargs  map[string]string `json:"kwargs"`
}

// Request is the data about the request this event occurred in.
type Request struct {
	// URL is full URL where this event occurred.
	URL string `json:"url"`
	// Method is the request method.
	Method string `json:"method"`
	// Headers object containing the request headers.
	// Header names should be formatted like they are in HTTP.
	Headers map[string][]string `json:"headers"`
	Params  *Params             `json:"params"`
	// GET query string params.
	GET map[string][]string `json:"GET"`
	// QueryString is the raw query string.
	QueryString string `json:"query_string"`
	// POST POST params.
	POST map[string][]string `json:"POST"`
	// Body is the raw POST body.
	Body string `json:"body"`
	// UserIP is the user's IP address as a string.
	// Can also be the special value "$remote_ip", which will be replaced with the source IP of the API request.
	// Will be indexed, as long as it is a valid IPv4 address.
	UserIP string `json:"user_ip"`
}

// Params any routing parameters. (i.e. for use with Rails Routes)
type Params struct {
	Controller string `json:"controller"`
	Action     string `json:"action"`
}

// Person is the user affected by this event. Will be indexed by ID, username, and email.
// People are stored in Rollbar keyed by ID. If you send a multiple different usernames/emails for the
// same ID, the last received values will overwrite earlier ones.
type Person struct {
	// ID is a string up to 40 characters identifying this user in your system.
	ID string `json:"id"`
	// Username a string up to 255 characters.
	Username string `json:"username,omitempty"`
	// Email a string up to 255 characters.
	Email string `json:"email,omitempty"`
}

// Server is a data about the server related to this event.
type Server struct {
	// Host is the server hostname. Will be indexed.
	Host string `json:"host,omitempty"`
	// Root is the path to the application code root, not including the final slash.
	// Used to collapse non-project code when displaying tracebacks.
	Root string `json:"root,omitempty"`
	// Branch name of the checked-out source control branch. Defaults to "master"
	Branch string `json:"branch,omitempty"`
	// CodeVersion string describing the running code version on the server.
	// See note about "code_version" above.
	CodeVersion string `json:"code_version,omitempty"`
	// Deprecated: Sha is Git SHA of the running code revision. Use the full sha.
	Sha string `json:"sha,omitempty"`
}

// Client is the data about the client device this event occurred on.
// As there can be multiple client environments for a given event (i.e. Flash running inside
// an HTML page), data should be namespaced by platform.
type Client struct {
	Javascript *Javascript `json:"javascript"`
}

// Javascript is the Rollbar understands the following field.
type Javascript struct {
	// Browser is the user agent string.
	Browser string `json:"browser"`
	// CodeVersion is the string describing the running code version in javascript
	// See note about "code_version" above.
	CodeVersion string `json:"code_version"`
	// SourceMapEnabled is the set to true to enable source map deobfuscation.
	// See the "Source Maps" guide for more details.
	SourceMapEnabled bool `json:"source_map_enabled"`
	// GuessUncaughtFrames is the set to true to enable frame guessing.
	// See the "Source Maps" guide for more details.
	GuessUncaughtFrames bool `json:"guess_uncaught_frames"`
}

// Notifier is the describes the library used to send this event.
type Notifier struct {
	// Name name of the library.
	Name string `json:"name"`
	// Version library version string.
	Version string `json:"version"`
}
