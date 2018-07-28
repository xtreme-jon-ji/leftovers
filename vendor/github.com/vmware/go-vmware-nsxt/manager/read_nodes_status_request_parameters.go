/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package manager

// Node Status list request parameters
type ReadNodesStatusRequestParameters struct {

	// A comma separated list of request Node Ids.
	NodeIds []string `json:"node_ids"`
}
