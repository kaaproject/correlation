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
	"reflect"
	"testing"
)

func TestContextWithCorrelationID(t *testing.T) {
	tests := []struct {
		name          string
		parent        context.Context
		correlationID []string
		want          context.Context
		generated     bool
	}{
		{
			name:          "context with normal correlation ID",
			parent:        context.Background(),
			correlationID: []string{"corrID-123"},
			want:          context.WithValue(context.Background(), correlationIDKey, "corrID-123"),
		},
		{
			name:          "context with empty correlation ID",
			parent:        context.Background(),
			correlationID: []string{""},
			want:          context.WithValue(context.Background(), correlationIDKey, ""),
			generated:     true,
		},
		{
			name:          "context with no correlation ID",
			parent:        context.Background(),
			correlationID: nil,
			want:          context.WithValue(context.Background(), correlationIDKey, ""),
			generated:     true,
		},
	}
	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			got := ContextWithCorrelationID(tc.parent, tc.correlationID...)
			if tc.generated {
				if got.Value(correlationIDKey) == "" {
					t.Error("ContextWithCorrelationID() must be not empty string")

				}
			} else {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("ContextWithCorrelationID() = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestID(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "retrieve normal correlation ID from the context",
			ctx:  context.WithValue(context.Background(), correlationIDKey, "corrID-123"),
			want: "corrID-123",
		},
		{
			name: "no correlation ID in the context",
			ctx:  context.Background(),
			want: "",
		},
	}
	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			if got := ID(tc.ctx); got != tc.want {
				t.Errorf("ID() = %v, want %v", got, tc.want)
			}
		})
	}
}
