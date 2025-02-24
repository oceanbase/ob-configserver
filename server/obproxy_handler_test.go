/**
 * Copyright 2025 OceanBase
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/oceanbase/configserver/config"
	"github.com/oceanbase/configserver/ent"
)

func TestParseVersionOnlyNormal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=GetObproxyConfig&VersionOnly=true", nil)
	versionOnly, err := isVersionOnly(c)
	require.True(t, versionOnly)
	require.True(t, err == nil)
}

func TestParseVersionOnlyError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=GetObproxyConfig&VersionOnly=abc", nil)
	_, err := isVersionOnly(c)
	require.True(t, err != nil)
}

func TestParseVersionOnlyNotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=GetObproxyConfig", nil)
	c.Params = []gin.Param{{Key: "Version", Value: "abc"}}
	versionOnly, err := isVersionOnly(c)
	require.False(t, versionOnly)
	require.True(t, err == nil)
}

func TestGetObProxyConfig(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=GetObproxyConfig", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	response := getObProxyConfig(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestGetObproxyConfigWithTemplate(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=GetObproxyConfig", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	response := getObProxyConfigWithTemplate(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}
