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

package opa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nuclio/logger"
	nucliozap "github.com/nuclio/zap"
	"github.com/stretchr/testify/suite"
)

type HTTPClientTestSuite struct {
	suite.Suite
	logger         logger.Logger
	testHTTPServer *httptest.Server
	httpClient     *HTTPClient
}

func (suite *HTTPClientTestSuite) SetupTest() {
	var err error
	suite.logger, err = nucliozap.NewNuclioZapTest("opa-test")
	suite.Require().NoError(err)

	allowPath := "/v1/data/authz/allow"
	filterPath := "/v1/data/authz/filter_allowed"

	// Create test HTTP server
	suite.testHTTPServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case allowPath:
			var permissionRequest PermissionQueryRequest
			err := json.NewDecoder(r.Body).Decode(&permissionRequest)
			suite.Require().NoError(err)

			// For testing, allow if resource starts with "allow"
			allowed := len(permissionRequest.Input.Resource) > 0 && permissionRequest.Input.Resource[0:5] == "allow"

			permissionResponse := PermissionQueryResponse{
				Result: allowed,
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(permissionResponse)
			suite.Require().NoError(err)

		case filterPath:
			var permissionRequest PermissionFilterRequest
			err := json.NewDecoder(r.Body).Decode(&permissionRequest)
			suite.Require().NoError(err)

			// For testing, allow resources that start with "allow"
			var allowedResources []string
			for _, resource := range permissionRequest.Input.Resources {
				if len(resource) > 0 && resource[0:5] == "allow" {
					allowedResources = append(allowedResources, resource)
				}
			}

			permissionResponse := PermissionFilterResponse{
				Result: allowedResources,
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(permissionResponse)
			suite.Require().NoError(err)
		}
	}))

	// Create HTTP client
	suite.httpClient = NewHTTPClient(
		suite.logger,
		suite.testHTTPServer.URL,
		allowPath,
		filterPath,
		5*time.Second,
		true, // Enable verbose logging for tests
		"test-override-value",
	)
}

func (suite *HTTPClientTestSuite) TearDownTest() {
	suite.testHTTPServer.Close()
}

func (suite *HTTPClientTestSuite) TestQueryPermissions_Allow() {
	// Test resource that should be allowed
	allowed, err := suite.httpClient.QueryPermissions(
		"allow-resource",
		ActionRead,
		&PermissionOptions{
			MemberIds: []string{"user1"},
		},
	)

	suite.Require().NoError(err)
	suite.Require().True(allowed)
}

func (suite *HTTPClientTestSuite) TestQueryPermissions_Deny() {
	// Test resource that should be denied
	allowed, err := suite.httpClient.QueryPermissions(
		"deny-resource",
		ActionRead,
		&PermissionOptions{
			MemberIds: []string{"user1"},
		},
	)

	suite.Require().NoError(err)
	suite.Require().False(allowed)
}

func (suite *HTTPClientTestSuite) TestQueryPermissions_WithOverride() {
	// Test with override header value
	allowed, err := suite.httpClient.QueryPermissions(
		"deny-resource", // Would normally be denied
		ActionRead,
		&PermissionOptions{
			MemberIds:           []string{"user1"},
			OverrideHeaderValue: "test-override-value", // This should bypass the check
		},
	)

	suite.Require().NoError(err)
	suite.Require().True(allowed)
}

func (suite *HTTPClientTestSuite) TestQueryPermissionsMultiResources() {
	resources := []string{
		"allow-resource-1",
		"deny-resource-1",
		"allow-resource-2",
		"deny-resource-2",
	}

	permissions, err := suite.httpClient.QueryPermissionsMultiResources(
		context.Background(),
		resources,
		ActionRead,
		&PermissionOptions{
			MemberIds: []string{"user1"},
		},
	)

	suite.Require().NoError(err)
	suite.Require().Equal(4, len(permissions))
	suite.Require().True(permissions[0])  // allow-resource-1
	suite.Require().False(permissions[1]) // deny-resource-1
	suite.Require().True(permissions[2])  // allow-resource-2
	suite.Require().False(permissions[3]) // deny-resource-2
}

func (suite *HTTPClientTestSuite) TestQueryPermissionsMultiResources_WithOverride() {
	resources := []string{
		"allow-resource-1",
		"deny-resource-1",
		"allow-resource-2",
		"deny-resource-2",
	}

	permissions, err := suite.httpClient.QueryPermissionsMultiResources(
		context.Background(),
		resources,
		ActionRead,
		&PermissionOptions{
			MemberIds:           []string{"user1"},
			OverrideHeaderValue: "test-override-value", // This should bypass the check
		},
	)

	suite.Require().NoError(err)
	suite.Require().Equal(4, len(permissions))
	suite.Require().True(permissions[0]) // All should be allowed with override
	suite.Require().True(permissions[1])
	suite.Require().True(permissions[2])
	suite.Require().True(permissions[3])
}

func TestHTTPClientTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientTestSuite))
}
