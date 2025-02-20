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
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	// log "github.com/sirupsen/logrus"

	"github.com/oceanbase/configserver/model"
)

const (
	CONFIG_URL_FORMAT = "%s/services?Action=ObRootServiceInfo&ObCluster=%s"
)

var obProxyConfigOnce sync.Once
var obProxyConfigFunc func(*gin.Context)
var obProxyConfigWithTemplateOnce sync.Once
var obProxyConfigWithTemplateFunc func(*gin.Context)

func getObProxyConfigFunc() func(*gin.Context) {
	obProxyConfigOnce.Do(func() {
		obProxyConfigFunc = handlerFunctionWrapper(getObProxyConfig)
	})
	return obProxyConfigFunc
}

func getObProxyConfigWithTemplateFunc() func(*gin.Context) {
	obProxyConfigWithTemplateOnce.Do(func() {
		obProxyConfigWithTemplateFunc = handlerFunctionWrapper(getObProxyConfigWithTemplate)
	})
	return obProxyConfigWithTemplateFunc
}

func getServiceAddress() string {
	return fmt.Sprintf("http://%s:%d", GetConfigServer().Config.Vip.Address, GetConfigServer().Config.Vip.Port)
}

func isVersionOnly(c *gin.Context) (bool, error) {
	ret := false
	var err error
	versionOnly, ok := c.GetQuery("VersionOnly")
	if ok {
		ret, err = strconv.ParseBool(versionOnly)
	}
	return ret, err
}

func getObProxyConfig(ctxlog context.Context, c *gin.Context) *ApiResponse {
	var response *ApiResponse
	client := GetConfigServer().Client

	versionOnly, err := isVersionOnly(c)
	if err != nil {
		return NewIllegalArgumentResponse(errors.Wrap(err, "invalid parameter, failed to parse versiononly"))
	}

	rootServiceInfoUrlMap := make(map[string]*model.RootServiceInfoUrl)
	clusters, err := client.ObCluster.Query().All(context.Background())
	if err != nil {
		return NewErrorResponse(errors.Wrap(err, "query ob clusters"))
	}

	for _, cluster := range clusters {
		rootServiceInfoUrlMap[cluster.Name] = &model.RootServiceInfoUrl{
			ObCluster: cluster.Name,
			Url:       fmt.Sprintf(CONFIG_URL_FORMAT, getServiceAddress(), cluster.Name),
		}
	}
	rootServiceInfoUrls := make([]*model.RootServiceInfoUrl, 0, len(rootServiceInfoUrlMap))
	for _, info := range rootServiceInfoUrlMap {
		rootServiceInfoUrls = append(rootServiceInfoUrls, info)
	}
	obProxyConfig, err := model.NewObProxyConfig(getServiceAddress(), rootServiceInfoUrls)
	if err != nil {
		response = NewErrorResponse(errors.Wrap(err, "generate obproxy config"))
	} else {
		if versionOnly {
			response = NewSuccessResponse(model.NewObProxyConfigVersionOnly(obProxyConfig.Version))
		} else {
			response = NewSuccessResponse(obProxyConfig)
		}
	}
	return response
}

func getObProxyConfigWithTemplate(ctxlog context.Context, c *gin.Context) *ApiResponse {
	var response *ApiResponse
	client := GetConfigServer().Client

	versionOnly, err := isVersionOnly(c)
	if err != nil {
		return NewIllegalArgumentResponse(errors.Wrap(err, "invalid parameter, failed to parse versiononly"))
	}

	clusterMap := make(map[string]interface{})
	clusters, err := client.ObCluster.Query().All(context.Background())

	if err != nil {
		return NewErrorResponse(errors.Wrap(err, "query ob clusters"))
	}

	for _, cluster := range clusters {
		clusterMap[cluster.Name] = nil
	}
	clusterNames := make([]string, 0, len(clusterMap))
	for clusterName := range clusterMap {
		clusterNames = append(clusterNames, clusterName)
	}

	obProxyConfigWithTemplate, err := model.NewObProxyConfigWithTemplate(getServiceAddress(), clusterNames)

	if err != nil {
		response = NewErrorResponse(errors.Wrap(err, "generate obproxy config with template"))
	} else {
		if versionOnly {
			response = NewSuccessResponse(model.NewObProxyConfigVersionOnly(obProxyConfigWithTemplate.Version))
		} else {
			response = NewSuccessResponse(obProxyConfigWithTemplate)
		}
	}
	return response
}
