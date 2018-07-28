/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package manager

// Physical NIC specification
type Pnic struct {

	// device name or key
	DeviceName string `json:"device_name"`

	// Uplink name for this Pnic. This name will be used to reference this Pnic in other configurations.
	UplinkName string `json:"uplink_name"`
}
