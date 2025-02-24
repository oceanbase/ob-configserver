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
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/oceanbase/configserver/config"
	"github.com/oceanbase/configserver/ent"
	"github.com/oceanbase/configserver/lib/trace"
	"github.com/oceanbase/configserver/logger"
)

var configServer *ConfigServer

func GetConfigServer() *ConfigServer {
	return configServer
}

type ConfigServer struct {
	Config *config.ConfigServerConfig
	Server *HttpServer
	Client *ent.Client
}

func NewConfigServer(conf *config.ConfigServerConfig) *ConfigServer {
	server := &ConfigServer{
		Config: conf,
		Server: &HttpServer{
			Counter: new(Counter),
			Router:  gin.Default(),
			Server:  &http.Server{},
			Address: conf.Server.Address,
		},
		Client: nil,
	}
	configServer = server
	return configServer
}

func (server *ConfigServer) Run() error {
	client, err := ent.Open(server.Config.Storage.DatabaseType, server.Config.Storage.ConnectionUrl)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("initialize storage client with config %v", server.Config.Storage))
	}

	server.Client = client

	defer server.Client.Close()

	if err := server.Client.Schema.Create(context.Background()); err != nil {
		return errors.Wrap(err, "create configserver schema")
	}

	// start http server
	ctx, cancel := context.WithCancel(trace.ContextWithTraceId(logger.INIT_TRACEID))
	server.Server.Cancel = cancel

	// register route
	InitConfigServerRoutes(server.Server.Router)

	// run http server
	server.Server.Run(ctx)
	return nil
}
