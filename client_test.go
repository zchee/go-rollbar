// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	api "github.com/zchee/go-rollbar/api/v1"
	"golang.org/x/net/context"
)

func newClient(cl *httpClient) *client {
	return &client{
		debugClient:    cl,
		infoClient:     cl,
		errorClient:    cl,
		warnClient:     cl,
		criticalClient: cl,
	}
}

type dummyLogger struct {
	logger *log.Logger
}

func (l dummyLogger) Debugf(_ context.Context, format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l dummyLogger) Infof(_ context.Context, format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func TestNew(t *testing.T) {
	const testToken = "xxxxxxxxxxxxxxxx"
	httpDummyClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	hostName, err := os.Hostname()
	if err != nil {
		t.Errorf("failed get the OS hostname: %v", err)
	}
	dl := dummyLogger{log.New(ioutil.Discard, "", 0)}

	type args struct {
		token   string
		options []Option
	}
	tests := []struct {
		name       string
		args       args
		wantClient httpClient
	}{
		{
			name: "default",
			args: args{
				token: testToken,
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
		},
		{
			name: "with Client",
			args: args{
				token:   testToken,
				options: []Option{WithClient(httpDummyClient)},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      httpDummyClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
		},
		{
			name: "with Endpoint",
			args: args{
				token:   testToken,
				options: []Option{WithEndpoint("https://endpoint.example.com")},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    "https://endpoint.example.com",
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
		},
		{
			name: "with Logger",
			args: args{
				token:   testToken,
				options: []Option{WithLogger(dl)},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      dl,
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
		},
		{
			name: "with Debug",
			args: args{
				token:   testToken,
				options: []Option{WithDebug(true)},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      traceLogger{os.Stderr}, // if debug is true, also logger is traceLogger
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				debug:       true,
				stackskip:   3,
			},
		},
		{
			name: "with Environment",
			args: args{
				token:   testToken,
				options: []Option{WithEnvironment("production")},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "production",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
		},
		{
			name: "with Platform",
			args: args{
				token:   testToken,
				options: []Option{WithPlatform("google-app-engine")},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    "google-app-engine",
				serverHost:  hostName,
				stackskip:   3,
			},
		},
		{
			name: "with Codeversion",
			args: args{
				token:   testToken,
				options: []Option{WithCodeVersion("2.1.12")},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				codeVersion: "2.1.12",
				stackskip:   3,
			},
		},
		{
			name: "with ServerHost",
			args: args{
				token:   testToken,
				options: []Option{WithServerHost("localhost")},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  "localhost",
				stackskip:   3,
			},
		},
		{
			name: "with ServerRoot",
			args: args{
				token:   testToken,
				options: []Option{WithServerRoot("/app/src")},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				serverRoot:  "/app/src",
				stackskip:   3,
			},
		},
		{
			name: "with ServerBranch",
			args: args{
				token:   testToken,
				options: []Option{WithServerBranch("test-branch")},
			},
			wantClient: httpClient{
				token:        testToken,
				client:       http.DefaultClient,
				endpoint:     api.DefaultEndpoint,
				logger:       nilLogger{},
				environment:  "development",
				platform:     runtime.GOOS,
				serverHost:   hostName,
				serverBranch: "test-branch",
				stackskip:    3,
			},
		},
		{
			name: "with WithStackSkip",
			args: args{
				token:   testToken,
				options: []Option{WithStackSkip(3)},
			},
			wantClient: httpClient{
				token:       testToken,
				client:      http.DefaultClient,
				endpoint:    api.DefaultEndpoint,
				logger:      nilLogger{},
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := newClient(&tt.wantClient)
			got := New(tt.args.token, tt.args.options...)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%q. New(%v, %v) = %v, want %v", tt.name, tt.args.token, tt.args.options, got, want)
			}
		})
	}
}

func Test_httpClient_payload(t *testing.T) {
	if strings.HasPrefix(runtime.Version(), "devel") {
		t.Skip("go version is devel")
	}
	const testToken = "xxxxxxxxxxxxxxxx"
	hostName, err := os.Hostname()
	if err != nil {
		t.Errorf("failed get the OS hostname: %v", err)
	}
	now := time.Now().Unix()

	type fields struct {
		token        string
		environment  string
		platform     string
		codeVersion  string
		serverHost   string
		serverRoot   string
		serverBranch string
		stackskip    int
	}
	type args struct {
		level Level
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *api.Payload
	}{
		{
			name: "default",
			fields: fields{
				token:       testToken,
				environment: "development",
				platform:    runtime.GOOS,
				serverHost:  hostName,
				stackskip:   3,
			},
			args: args{
				level: DebugLevel,
				err:   errors.New("default error"),
			},
			want: &api.Payload{
				AccessToken: testToken,
				Data: &api.Data{
					Environment: "development",
					Body: &api.Body{
						Trace: testTrace,
					},
					Level:     "debug",
					Timestamp: now,
					Platform:  runtime.GOOS,
					Language:  language,
					Server: &api.Server{
						Host: hostName,
					},
					Fingerprint: testFingerprint,
					Title:       "default error",
					Notifier: &api.Notifier{
						Name:    "go-rollbar",
						Version: "0.0.0",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpClient{
				token:        tt.fields.token,
				environment:  tt.fields.environment,
				platform:     tt.fields.platform,
				codeVersion:  tt.fields.codeVersion,
				serverHost:   tt.fields.serverHost,
				serverRoot:   tt.fields.serverRoot,
				serverBranch: tt.fields.serverBranch,
				stackskip:    tt.fields.stackskip,
			}
			if got := c.payload(tt.args.level, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpClient.payload(%v, %v) = %v, want %v", tt.args.level, tt.args.err, got, tt.want)
			}
		})
	}
}
