// correlation package provides functions for propagating correlation IDs with context.
//
// Copyright 2019 KaaIoT Technologies, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package correlation

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

type key int

const (
	correlationIDKey        key    = 1234097
	correlationIDHeaderName string = "Correlation-ID"
)

// ContextWithCorrelationID sets correlation ID to context for further use. When no correlation ID is provided,
// the function generates a new correlation ID.
func ContextWithCorrelationID(parent context.Context, correlationID ...string) context.Context {
	var id string
	if len(correlationID) == 0 || len(correlationID[0]) == 0 {
		id = GenerateID()
	} else {
		id = correlationID[0]
	}

	return context.WithValue(parent, correlationIDKey, id)
}

// ID returns correlation ID from the context, if found. Otherwise the return value is an empty string.
func ID(ctx context.Context) string {
	id, _ := ctx.Value(correlationIDKey).(string)
	return id
}

// WithCorrelationID is an HTTP handler factory function for use in the handlers chain.
// The handler extracts "Correlation-ID" header from HTTP request and saves it in the request context.
// If there is no "Correlation-ID" header found in the request, a new correlation ID is generated.
func WithCorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(correlationIDHeaderName)
		if correlationID == "" {
			correlationID = GenerateID()
		}

		next.ServeHTTP(w, r.WithContext(ContextWithCorrelationID(r.Context(), correlationID)))
	})
}

// SetCorrelationID stored in the context in the HTTP request.
func SetCorrelationID(ctx context.Context, r *http.Request) {
	id := ID(ctx)
	if len(id) == 0 {
		id = GenerateID()
	}

	r.Header.Set(correlationIDHeaderName, id)
}

// GenerateID generates correlation ID
func GenerateID() string {
	return xid.New().String()
}
