// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/pkg/errors"
	api "github.com/zchee/go-rollbar/api/v1"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type Client struct {
	debug        *DebugService
	info         *InfoService
	errorService *ErrorService
	warn         *WarnService
	critical     *CriticalService
}

type httpClient struct {
	token       string
	client      *http.Client
	endpoint    string
	debug       bool
	logger      Logger
	environment string
	platform    string
	codeVersion string
	serverHost  string
	serverRoot  string
	id          string
	custom      map[string]interface{}
}

// New creates a new REST rollbar API client.
//
// The `token` is required, other optional parameters can be passed using the
// various `With...` functions.
func New(token string, options ...Option) Rollbar {
	cl := http.DefaultClient
	endpoint := api.DefaultEndpoint
	debug, _ := strconv.ParseBool(os.Getenv("ROLLBAR_DEBUG"))
	environment := "development"
	var (
		logger      Logger = nilLogger{}
		platform    string
		codeVersion string
		serverHost  string
		serverRoot  string
	)

	for _, o := range options {
		switch o.Key() {
		case keyHTTPClient:
			cl = o.Value().(*http.Client)
		case keyEndpoint:
			endpoint = o.Value().(string)
		case keyDebug:
			debug = o.Value().(bool)
		case keyLogger:
			logger = o.Value().(Logger)
		case keyEnvironment:
			environment = o.Value().(string)
		case keyPlatform:
			platform = o.Value().(string)
		case keyCodeVersion:
			codeVersion = o.Value().(string)
		case keyServerHost:
			serverHost = o.Value().(string)
		case keyServerRoot:
			serverRoot = o.Value().(string)
		}
	}

	if _, ok := logger.(nilLogger); debug && ok {
		logger = traceLogger{os.Stderr}
	}
	if platform == "" {
		platform = runtime.GOOS
	}
	if serverHost == "" {
		serverHost, _ = os.Hostname()
	}

	client := &httpClient{
		token:       token,
		client:      cl,
		endpoint:    endpoint,
		debug:       debug,
		logger:      logger,
		environment: environment,
		platform:    platform,
		codeVersion: codeVersion,
		serverHost:  serverHost,
		serverRoot:  serverRoot,
	}

	return &Client{
		debug:        &DebugService{client: client},
		info:         &InfoService{client: client},
		errorService: &ErrorService{client: client},
		warn:         &WarnService{client: client},
		critical:     &CriticalService{client: client},
	}
}

// payload creates the rollbar payload data.
func (c *httpClient) payload(level Level, err error) *api.Payload {
	title := "<nil>"
	if err != nil {
		title = err.Error()
	}

	data := &api.Data{
		Title:       title,
		Body:        errorBody(err, CreateStack(4)),
		Environment: c.environment,
		Level:       string(level),
		Timestamp:   time.Now().Unix(),
		Platform:    c.platform,
		Language:    language,
		Server: &api.Server{
			Host: c.serverHost,
			Root: c.serverRoot,
		},
		Notifier: &api.Notifier{
			Name:    Name,
			Version: Version,
		},
	}

	return &api.Payload{
		AccessToken: c.token,
		Data:        data,
	}
}

// post posts payload to rollbar.
func (c *httpClient) post(pctx context.Context, payload *api.Payload) error {
	if c.token == "" {
		return errors.New("empty token")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode payload")
	}
	if c.debug {
		out, _ := json.MarshalIndent(payload, "", "  ")
		c.logger.Debug(pctx, string(out))
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(data))
	if err != nil {
		return errors.New("failed to create new POST request")
	}

	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	req.Header.Set("Content-Type", "application/json")
	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return errors.Wrap(err, "failed to POST to rollbar")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("received response: %s", resp.Status)
	}

	return c.parseResponse(pctx, resp.Body, data)
}

// parseResponse parses the rollbar API response.
func (c *httpClient) parseResponse(ctx context.Context, rdr io.Reader, data interface{}) error {
	if c.debug {
		var buf bytes.Buffer
		io.Copy(&buf, rdr)

		c.logger.Debug(ctx, "-----> %s (response)", c.endpoint)
		var m api.Response
		if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
			c.logger.Debug(ctx, "failed to unmarshal payload: %s", err)
			c.logger.Debug(ctx, "%s", buf.String())
		} else {
			formatted, _ := json.MarshalIndent(m, "", "  ")
			c.logger.Debug(ctx, "%s", formatted)
		}
		c.logger.Debug(ctx, "<----- %s (response)", c.endpoint)
		rdr = &buf
	}
	return json.NewDecoder(rdr).Decode(&data)
}
