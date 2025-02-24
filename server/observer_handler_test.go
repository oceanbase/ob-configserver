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
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/oceanbase/configserver/config"
	"github.com/oceanbase/configserver/ent"
)

const testRootServiceJson = "{\"Type\":\"PRIMARY\",\"ObClusterId\":1,\"ObRegionId\":1,\"ObCluster\":\"c1\",\"ObRegion\":\"c1\",\"ReadonlyRsList\":[],\"RsList\":[{\"address\":\"1.1.1.1:2882\",\"role\":\"LEADER\",\"sql_port\":2881}],\"timestamp\":1649435362283000}"

func TestGetRootServiceInfoParamOldVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObRegion=c1&ObRegionId=1", nil)

	obRootServiceInfoParam, err := getCommonParam(c)
	require.Equal(t, "c1", obRootServiceInfoParam.ObCluster)
	require.Equal(t, int64(1), obRootServiceInfoParam.ObClusterId)
	require.Equal(t, 0, obRootServiceInfoParam.Version)
	require.True(t, err == nil)

}

func TestGetRootServiceInfoParamVersion2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&ObClusterId=1&version=2", nil)

	obRootServiceInfoParam, err := getCommonParam(c)
	require.Equal(t, "c1", obRootServiceInfoParam.ObCluster)
	require.Equal(t, int64(1), obRootServiceInfoParam.ObClusterId)
	require.Equal(t, 2, obRootServiceInfoParam.Version)
	require.True(t, err == nil)

}

func TestGetObRootServiceInfo(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	client.ObCluster.
		Create().
		SetName("c1").
		SetObClusterID(1).
		SetType("PRIMARY").
		SetRootserviceJSON(testRootServiceJson).
		OnConflict().
		SetRootserviceJSON(testRootServiceJson).
		Exec(context.Background())

	response := getObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestGetObRootServiceInfoV2(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&version=2", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}
	client.ObCluster.
		Create().
		SetName("c1").
		SetObClusterID(1).
		SetType("PRIMARY").
		SetRootserviceJSON(testRootServiceJson).
		OnConflict().
		SetRootserviceJSON(testRootServiceJson).
		Exec(context.Background())

	response := getObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestGetObRootServiceInfoV2WithObClusterId(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&ObClusterId=1&version=2", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	client.ObCluster.
		Create().
		SetName("c1").
		SetObClusterID(1).
		SetType("PRIMARY").
		SetRootserviceJSON(testRootServiceJson).
		OnConflict().
		SetRootserviceJSON(testRootServiceJson).
		Exec(context.Background())

	response := getObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestGetObRootServiceInfoNoResult(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c2&ObClusterId=2&version=2", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	response := getObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusNotFound, response.Code)
}

func TestCreateOrUpdateObRootServiceInfo(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&ObClusterId=1&version=2", bytes.NewBuffer([]byte(testRootServiceJson)))

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	response := createOrUpdateObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteObRootServiceInfo(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&ObClusterId=1&version=2", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	client.ObCluster.
		Create().
		SetName("c1").
		SetObClusterID(1).
		SetType("PRIMARY").
		SetRootserviceJSON(testRootServiceJson).
		OnConflict().
		SetRootserviceJSON(testRootServiceJson).
		Exec(context.Background())

	response := deleteObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteObRootServiceInfoVersion1(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&ObClusterId=1", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	client.ObCluster.
		Create().
		SetName("c1").
		SetObClusterID(1).
		SetType("PRIMARY").
		SetRootserviceJSON(testRootServiceJson).
		OnConflict().
		SetRootserviceJSON(testRootServiceJson).
		Exec(context.Background())

	response := deleteObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusBadRequest, response.Code)
}

func TestDeleteObRootServiceInfoWithoutClusterId(t *testing.T) {
	// test gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "http://1.1.1.1:8080/services?Action=ObRootServiceInfo&ObCluster=c1&version=2", nil)

	// mock db client
	client, _ := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client.Schema.Create(context.Background())

	configServerConfig, _ := config.ParseConfigServerConfig("../etc/config.yaml")
	configServer = &ConfigServer{
		Config: configServerConfig,
		Client: client,
	}

	client.ObCluster.
		Create().
		SetName("c1").
		SetObClusterID(1).
		SetType("PRIMARY").
		SetRootserviceJSON(testRootServiceJson).
		OnConflict().
		SetRootserviceJSON(testRootServiceJson).
		Exec(context.Background())

	response := deleteObRootServiceInfo(context.Background(), c)
	require.Equal(t, http.StatusBadRequest, response.Code)
}
