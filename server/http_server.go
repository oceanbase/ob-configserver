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
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	libhttp "github.com/oceanbase/configserver/lib/http"
)

type HttpServer struct {
	// server will be stopped, new request will be rejected
	Stopping int32
	// current session count, concurrent safely
	Counter *Counter
	// http routers
	Router *gin.Engine
	// address
	Address string
	// http server, call its Run, Shutdown methods
	Server *http.Server
	// stop the http.Server by calling cancel method
	Cancel context.CancelFunc
}

// UseCounter use counter middleware
func (server *HttpServer) UseCounter() {
	server.Router.Use(
		server.counterPreHandlerFunc,
		server.counterPostHandlerFunc,
	)
}

// Run start a httpServer
// when ctx is cancelled, call shutdown to stop the httpServer
func (server *HttpServer) Run(ctx context.Context) {

	server.Server.Handler = server.Router
	if server.Address != "" {
		log.WithContext(ctx).Infof("listen on address: %s", server.Address)
		tcpListener, err := libhttp.NewTcpListener(server.Address)
		if err != nil {
			log.WithError(err).
				Errorf("create tcp listener on address '%s' failed %v", server.Address, err)
			return
		}
		go func() {
			if err := server.Server.Serve(tcpListener); err != nil {
				log.WithError(err).
					Info("tcp server exited")
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			if err := server.Shutdown(ctx); err != nil {
				log.WithContext(ctx).
					WithError(err).
					Error("server shutdown failed!")
				// in a for loop, sleep 100ms
				time.Sleep(time.Millisecond * 100)
			} else {
				log.WithContext(ctx).Info("server shutdown successfully.")
				return
			}
		}
	}
}

// shutdown httpServer can shutdown if sessionCount is 0,
// otherwise, return an error
func (server *HttpServer) Shutdown(ctx context.Context) error {
	atomic.StoreInt32(&(server.Stopping), 1)
	sessionCount := atomic.LoadInt32(&server.Counter.sessionCount)
	if sessionCount > 0 {
		return errors.Errorf("server shutdown failed, cur-session count:%d, shutdown will be success when wait session-count is 0.", sessionCount)
	}
	return server.Server.Close()
}

// counterPreHandlerFunc middleware for httpServer session count, before process a request
func (server *HttpServer) counterPreHandlerFunc(c *gin.Context) {
	if atomic.LoadInt32(&(server.Stopping)) == 1 {
		c.Abort()
		c.JSON(http.StatusServiceUnavailable, "server is shutdowning now.")
		return
	}

	server.Counter.incr()

	c.Next()
}

// counterPostHandlerFunc middleware for httpServer session count, after process a request
func (server *HttpServer) counterPostHandlerFunc(c *gin.Context) {
	c.Next()
	server.Counter.decr()
}

// counter session counter
// when server receive a request, sessionCount +1,
// when the request returns a response, sessionCount -1.
type Counter struct {
	sessionCount int32
	sync.Mutex
}

// incr sessionCount +1 concurrent safely
func (c *Counter) incr() {
	c.Lock()
	c.sessionCount++
	defer c.Unlock()
}

// decr sessionCount -1 concurrent safely
func (c *Counter) decr() {
	c.Lock()
	c.sessionCount--
	defer c.Unlock()
}
