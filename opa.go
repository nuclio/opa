/*
Copyright 2025 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package opa provides a Go client library for Open Policy Agent (OPA) with support for HTTP-based policy queries.
//
// The package supports multiple client types:
//   - HTTPClient: Production client for communicating with OPA over HTTP
//   - NopClient: Always returns true, useful for development/testing
//   - MockClient: Test client using testify/mock for unit testing
//
// Example usage:
//
//	config := &opa.Config{
//		ClientKind:           opa.ClientKindHTTP,
//		Address:             "http://localhost:8181",
//		PermissionQueryPath: "/v1/data/authz/allow",
//		RequestTimeout:      10,
//	}
//
//	client := opa.CreateOpaClient(logger, config)
//	allowed, err := client.QueryPermissions("resource1", opa.ActionRead, &opa.PermissionOptions{
//		MemberIds: []string{"user123"},
//	})
package opa

import (
	"context"
)

// Client represents an OPA client that can query permissions.
type Client interface {
	// QueryPermissions queries permission for a single resource.
	QueryPermissions(context.Context, string, Action, *PermissionOptions) (bool, error)

	// QueryPermissionsMultiResources queries permissions for multiple resources at once.
	// Returns a slice of booleans where each index corresponds to the resource at the same index.
	QueryPermissionsMultiResources(context.Context, []string, Action, *PermissionOptions) ([]bool, error)
}
