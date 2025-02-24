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
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/oceanbase/configserver/lib/codec"
	"github.com/oceanbase/configserver/lib/net"
	"github.com/oceanbase/configserver/lib/trace"
)

var invalidActionOnce sync.Once
var invalidActionFunc func(*gin.Context)

func getInvalidActionFunc() func(*gin.Context) {
	invalidActionOnce.Do(func() {
		invalidActionFunc = handlerFunctionWrapper(invalidAction)
	})
	return invalidActionFunc
}

func getServerIdentity() string {
	ip, _ := net.GetLocalIpAddress()
	return ip
}

func handlerFunctionWrapper(f func(context.Context, *gin.Context) *ApiResponse) func(*gin.Context) {
	fn := func(c *gin.Context) {
		tStart := time.Now()
		traceId := trace.RandomTraceId()
		ctxlog := trace.ContextWithTraceId(traceId)
		log.WithContext(ctxlog).Infof("handle request: %s %s", c.Request.Method, c.Request.RequestURI)
		response := f(ctxlog, c)
		cost := time.Now().Sub(tStart).Milliseconds()
		response.TraceId = traceId
		response.Cost = cost
		response.Server = getServerIdentity()
		responseJson, err := codec.MarshalToJsonString(response)
		if err != nil {
			log.WithContext(ctxlog).Errorf("response: %s", "response serialization error")
			c.JSON(http.StatusInternalServerError, NewErrorResponse(errors.Wrap(err, "serialize response")))
		} else {
			log.WithContext(ctxlog).Infof("response: %s", responseJson)
			c.String(response.Code, string(responseJson))
		}
	}
	return fn
}

func invalidAction(ctxlog context.Context, c *gin.Context) *ApiResponse {
	log.WithContext(ctxlog).Error("invalid action")
	return NewIllegalArgumentResponse(errors.New("invalid action"))
}

func getHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		action := c.Query("Action")
		switch action {
		case "ObRootServiceInfo":
			getObRootServiceGetFunc()(c)

		case "GetObProxyConfig":
			getObProxyConfigFunc()(c)

		case "GetObRootServiceInfoUrlTemplate":
			getObProxyConfigWithTemplateFunc()(c)

		case "ObIDCRegionInfo":
			getObIdcRegionInfoFunc()(c)

		default:
			getInvalidActionFunc()(c)
		}
	}
	return gin.HandlerFunc(fn)
}

func postHandler() gin.HandlerFunc {

	fn := func(c *gin.Context) {
		action, _ := c.GetQuery("Action")
		switch action {
		case "ObRootServiceInfo":
			getObRootServicePostFunc()(c)

		case "GetObProxyConfig":
			getObProxyConfigFunc()(c)

		case "GetObRootServiceInfoUrlTemplate":
			getObProxyConfigWithTemplateFunc()(c)
		default:
			getInvalidActionFunc()(c)
		}
	}

	return gin.HandlerFunc(fn)
}

func deleteHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		action, _ := c.GetQuery("Action")
		switch action {
		case "ObRootServiceInfo":
			getObRootServiceDeleteFunc()(c)
		default:
			getInvalidActionFunc()(c)
		}
	}
	return gin.HandlerFunc(fn)
}
