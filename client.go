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
	debugClient    *httpClient
	infoClient     *httpClient
	errorClient    *httpClient
	warnClient     *httpClient
	criticalClient *httpClient
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
}

var defaultHTTPClient = httpClient{
	client:      http.DefaultClient,
	endpoint:    api.DefaultEndpoint,
	logger:      nilLogger{},
	environment: "development",
	platform:    runtime.GOOS,
}

// New creates a new REST rollbar API client.
//
// The `token` is required, other optional parameters can be passed using the
// various `With...` functions.
func New(token string, options ...Option) Rollbar {
	client := defaultHTTPClient
	client.token = token
	if debug, err := strconv.ParseBool(os.Getenv("ROLLBAR_DEBUG")); err == nil && debug {
		client.debug = debug
	}

	for _, o := range options {
		o(&client)
	}
	if _, ok := client.logger.(nilLogger); client.debug && ok {
		client.logger = traceLogger{os.Stderr}
	}
	if client.serverHost == "" {
		client.serverHost, _ = os.Hostname()
	}

	return &Client{
		debugClient:    &client,
		infoClient:     &client,
		errorClient:    &client,
		warnClient:     &client,
		criticalClient: &client,
	}
}

// payload creates the rollbar payload data.
func (c *httpClient) payload(level Level, err error) *api.Payload {
	title := "<nil>"
	if err != nil {
		title = err.Error()
	}
	stack := CreateStack(3)

	data := &api.Data{
		Environment: c.environment,
		Body:        errorBody(err, stack),
		Level:       string(level),
		Timestamp:   time.Now().Unix(),
		Platform:    c.platform,
		Language:    language,
		Server: &api.Server{
			Host: c.serverHost,
			Root: c.serverRoot,
		},
		Fingerprint: stack.Fingerprint(),
		Title:       title,
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
func (c *httpClient) post(pctx context.Context, payload *api.Payload) (*api.Response, error) {
	if c.token == "" {
		return nil, errors.New("empty token")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode payload")
	}
	if c.debug {
		out, _ := json.MarshalIndent(payload, "", "  ")
		c.logger.Debugf(pctx, string(out)+"\n")
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("failed to create new POST request")
	}

	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	req.Header.Set("Content-Type", "application/json")
	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to POST to rollbar")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received response: %s", resp.Status)
	}

	return c.parseResponse(ctx, resp.Body), nil
}

// parseResponse parses the rollbar API response.
func (c *httpClient) parseResponse(ctx context.Context, rdr io.Reader) *api.Response {
	buf := new(bytes.Buffer)
	io.Copy(buf, rdr)

	c.logger.Debugf(ctx, "-----> %s (response)\n", c.endpoint)
	m := new(api.Response)
	if err := json.Unmarshal(buf.Bytes(), m); err != nil {
		c.logger.Debugf(ctx, "failed to unmarshal payload: %s\n", err)
		c.logger.Debugf(ctx, "%s\n", buf.String())
	} else {
		formatted, _ := json.MarshalIndent(m, "", "  ")
		c.logger.Debugf(ctx, "%s\n", formatted)
	}
	c.logger.Debugf(ctx, "<----- %s (response)\n", c.endpoint)

	return m
}
