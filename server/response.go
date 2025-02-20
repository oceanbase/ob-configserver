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
	"fmt"
	"net/http"
)

type ApiResponse struct {
	Code       int         `json:"Code"`
	Message    string      `json:"Message"`
	Successful bool        `json:"Success"`
	Data       interface{} `json:"Data"`
	TraceId    string      `json:"Trace"`
	Server     string      `json:"Server"`
	Cost       int64       `json:"Cost"`
}

type IterableData struct {
	Contents interface{} `json:"Contents"`
}

func NewSuccessResponse(data interface{}) *ApiResponse {
	return &ApiResponse{
		Code:       http.StatusOK,
		Message:    "successful",
		Successful: true,
		Data:       data,
	}
}

func NewBadRequestResponse(err error) *ApiResponse {
	return &ApiResponse{
		Code:       http.StatusBadRequest,
		Message:    fmt.Sprintf("bad request: %v", err),
		Successful: false,
	}
}

func NewIllegalArgumentResponse(err error) *ApiResponse {
	return &ApiResponse{
		Code:       http.StatusBadRequest,
		Message:    fmt.Sprintf("illegal argument: %v", err),
		Successful: false,
	}
}

func NewNotFoundResponse(err error) *ApiResponse {
	return &ApiResponse{
		Code:       http.StatusNotFound,
		Message:    fmt.Sprintf("resource not found: %v", err),
		Successful: false,
	}
}

func NewNotImplementedResponse(err error) *ApiResponse {
	return &ApiResponse{
		Code:       http.StatusNotImplemented,
		Message:    fmt.Sprintf("request not implemented: %v", err),
		Successful: false,
	}
}

func NewErrorResponse(err error) *ApiResponse {
	return &ApiResponse{
		Code:       http.StatusInternalServerError,
		Message:    fmt.Sprintf("got internal error: %v", err),
		Successful: false,
	}
}
