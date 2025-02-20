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
	"fmt"
	"testing"
	"time"
)

func TestNewTcpListener(t *testing.T) {
	listener1, err := NewTcpListener("127.0.0.1:9998")
	if err != nil {
		t.Error("run tcp listener failed", err)
	}
	defer listener1.Close()
	listener2, err := NewTcpListener("127.0.0.1:9998")
	if err == nil {
		t.Error("run tcp listener on same port")
	}
	defer listener2.Close()
}

func TestStartTcp(t *testing.T) {
	defer time.Sleep(time.Millisecond)

	listener1 := NewListener()
	err := listener1.StartTCP("127.0.0.1:9998")
	if err != nil {
		t.Error("run StartTCP failed", err)
	}
	defer listener1.Close()

	listener2 := NewListener()
	err = listener2.StartTCP("127.0.0.1:9998")
	if err == nil {
		t.Error("run StartTCP on same port")
	}
	defer listener2.Close()

}

func TestStartErr(t *testing.T) {
	defer time.Sleep(time.Millisecond)

	listener1 := NewListener()
	err := listener1.StartTCP("127.0.0.11:9998")
	if err == nil {
		fmt.Println("run StartTCP on bad address should failed")
		//t.Error("run StartTCP on bad address should failed")
		defer listener1.Close()
	}
	fmt.Println(err)

	err = listener1.StartSocket("/not/exist_file")
	if err == nil {
		t.Error("run StartTCP on bad address should failed")
		defer listener1.Close()
	}
	fmt.Println(err)
}
