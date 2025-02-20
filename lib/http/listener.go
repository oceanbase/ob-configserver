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

package http

import (
	"context"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Listener struct {
	tcpListener  *net.TCPListener
	unixListener *net.UnixListener
	mux          *http.ServeMux
	srv          *http.Server
}

func NewListener() *Listener {
	mux := http.NewServeMux()
	return &Listener{
		mux: mux,
		srv: &http.Server{Handler: mux},
	}
}

func (l *Listener) StartTCP(addr string) error {
	tcpListener, err := NewTcpListener(addr)
	if err != nil {
		return err
	}
	go func() {
		_ = l.srv.Serve(tcpListener)
		log.Info("http tcp server exited")
	}()
	l.tcpListener = tcpListener
	return nil
}

func NewTcpListener(addr string) (*net.TCPListener, error) {
	cfg := net.ListenConfig{}
	listener, err := cfg.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return nil, err
	}
	return listener.(*net.TCPListener), nil
}

func (l *Listener) StartSocket(path string) error {
	listener, err := NewSocketListener(path)
	if err != nil {
		return err
	}

	go func() {
		_ = l.srv.Serve(listener)
		log.Info("http socket server exited")
	}()
	l.unixListener = listener
	return nil
}

func NewSocketListener(path string) (*net.UnixListener, error) {
	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return nil, err
	}
	return net.ListenUnix("unix", addr)
}

func (l *Listener) AddHandler(path string, h http.Handler) {
	l.mux.Handle(path, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Connection", "close")
		h.ServeHTTP(writer, request)
	}))
}

func (l *Listener) Close() {
	//var err error
	_ = l.srv.Close()
	if l.tcpListener != nil {
		_ = l.tcpListener.Close()
		//if err != nil {
		//	log.WithError(err).Warn("close tcpListener got error")
		//}
	}
	if l.unixListener != nil {
		_ = l.unixListener.Close()
		//if err != nil {
		//	log.WithError(err).Warn("close unixListener got error")
		//}
	}
}
