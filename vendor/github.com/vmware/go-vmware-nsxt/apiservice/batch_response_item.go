/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package apiservice

type BatchResponseItem struct {

	// object returned by api
	Body *interface{} `json:"body,omitempty"`

	// http status code
	Code int64 `json:"code"`

	// The headers returned by the API call
	Headers *interface{} `json:"headers,omitempty"`
}
