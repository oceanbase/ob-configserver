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
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCounter(t *testing.T) {
	count := new(Counter)
	Convey("counter after init", t, func() {
		So(count.sessionCount, ShouldEqual, 0)
	})

	count.incr()
	Convey("counter after incr", t, func() {
		So(count.sessionCount, ShouldEqual, 1)
	})

	count.decr()
	Convey("counter after decr", t, func() {
		So(count.sessionCount, ShouldEqual, 0)
	})
}

func TestHttpServer(t *testing.T) {
	server := &HttpServer{
		Counter: new(Counter),
		Router:  gin.Default(),
		Server: &http.Server{
			Addr: ":0",
		},
	}

	w := httptest.NewRecorder()
	server.UseCounter()
	server.Router.GET("/foo", fooHandler)
	end := make(chan bool, 1)
	handler := func(w http.ResponseWriter, r *http.Request) {
		server.Router.ServeHTTP(w, r)
		time.Sleep(time.Second)
	}
	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	go func() {
		handler(w, req)
		end <- true
	}()

	time.Sleep(10 * time.Millisecond)
	t.Run("handle a 1 second request", func(t *testing.T) {
		Convey("session count should be 1", t, func() {
			So(server.Counter.sessionCount, ShouldEqual, 1)
		})

		err := server.Shutdown(context.Background())
		Convey("server shutdown should fail", t, func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "server shutdown failed")
		})
	})

	<-end
	t.Run("handle request end", func(t *testing.T) {
		Convey("session count should be 0", t, func() {
			So(server.Counter.sessionCount, ShouldEqual, 0)
		})
		err := server.Shutdown(context.Background())
		Convey("server shutdown should success", t, func() {
			So(err, ShouldBeNil)
		})
	})
}

func fooHandler(c *gin.Context) {
	time.Sleep(time.Second)
}
