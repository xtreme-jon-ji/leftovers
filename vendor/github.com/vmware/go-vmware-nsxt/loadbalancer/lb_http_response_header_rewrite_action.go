/*
 * NSX API
 *
 * VMware NSX REST API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package loadbalancer

// This action is used to rewrite header fields of HTTP response messages to specified new values at HTTP_RESPONSE_REWRITE phase. One action can be used to rewrite one header field. To rewrite multiple header fields, multiple actions must be defined. Captured variables and built-in variables can be used in the header_value field, header_name field does not support variables.
type LbHttpResponseHeaderRewriteAction struct {

	// Type of load balancer rule action
	Type_ string `json:"type"`

	// Name of a header field of HTTP request message
	HeaderName string `json:"header_name"`

	// Value of header field
	HeaderValue string `json:"header_value"`
}
