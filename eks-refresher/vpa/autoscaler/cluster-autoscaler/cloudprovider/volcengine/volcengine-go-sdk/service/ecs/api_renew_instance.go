// Code generated by volcengine with private/model/cli/gen-api/main.go. DO NOT EDIT.

package ecs

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/request"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/response"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/volcengineutil"
)

const opRenewInstanceCommon = "RenewInstance"

// RenewInstanceCommonRequest generates a "volcengine/request.Request" representing the
// client's request for the RenewInstanceCommon operation. The "output" return
// value will be populated with the RenewInstanceCommon request's response once the request completes
// successfully.
//
// Use "Send" method on the returned RenewInstanceCommon Request to send the API call to the service.
// the "output" return value is not valid until after RenewInstanceCommon Send returns without error.
//
// See RenewInstanceCommon for more information on using the RenewInstanceCommon
// API call, and error handling.
//
//	// Example sending a request using the RenewInstanceCommonRequest method.
//	req, resp := client.RenewInstanceCommonRequest(params)
//
//	err := req.Send()
//	if err == nil { // resp is now filled
//	    fmt.Println(resp)
//	}
func (c *ECS) RenewInstanceCommonRequest(input *map[string]interface{}) (req *request.Request, output *map[string]interface{}) {
	op := &request.Operation{
		Name:       opRenewInstanceCommon,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &map[string]interface{}{}
	}

	output = &map[string]interface{}{}
	req = c.newRequest(op, input, output)

	return
}

// RenewInstanceCommon API operation for ECS.
//
// Returns volcengineerr.Error for service API and SDK errors. Use runtime type assertions
// with volcengineerr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the VOLCENGINE API reference guide for ECS's
// API operation RenewInstanceCommon for usage and error information.
func (c *ECS) RenewInstanceCommon(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := c.RenewInstanceCommonRequest(input)
	return out, req.Send()
}

// RenewInstanceCommonWithContext is the same as RenewInstanceCommon with the addition of
// the ability to pass a context and additional request options.
//
// See RenewInstanceCommon for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If the context is nil a panic will occur.
// In the future the SDK may create sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ECS) RenewInstanceCommonWithContext(ctx volcengine.Context, input *map[string]interface{}, opts ...request.Option) (*map[string]interface{}, error) {
	req, out := c.RenewInstanceCommonRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opRenewInstance = "RenewInstance"

// RenewInstanceRequest generates a "volcengine/request.Request" representing the
// client's request for the RenewInstance operation. The "output" return
// value will be populated with the RenewInstanceCommon request's response once the request completes
// successfully.
//
// Use "Send" method on the returned RenewInstanceCommon Request to send the API call to the service.
// the "output" return value is not valid until after RenewInstanceCommon Send returns without error.
//
// See RenewInstance for more information on using the RenewInstance
// API call, and error handling.
//
//	// Example sending a request using the RenewInstanceRequest method.
//	req, resp := client.RenewInstanceRequest(params)
//
//	err := req.Send()
//	if err == nil { // resp is now filled
//	    fmt.Println(resp)
//	}
func (c *ECS) RenewInstanceRequest(input *RenewInstanceInput) (req *request.Request, output *RenewInstanceOutput) {
	op := &request.Operation{
		Name:       opRenewInstance,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &RenewInstanceInput{}
	}

	output = &RenewInstanceOutput{}
	req = c.newRequest(op, input, output)

	return
}

// RenewInstance API operation for ECS.
//
// Returns volcengineerr.Error for service API and SDK errors. Use runtime type assertions
// with volcengineerr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the VOLCENGINE API reference guide for ECS's
// API operation RenewInstance for usage and error information.
func (c *ECS) RenewInstance(input *RenewInstanceInput) (*RenewInstanceOutput, error) {
	req, out := c.RenewInstanceRequest(input)
	return out, req.Send()
}

// RenewInstanceWithContext is the same as RenewInstance with the addition of
// the ability to pass a context and additional request options.
//
// See RenewInstance for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. Ifthe context is nil a panic will occur.
// In the future the SDK may create sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ECS) RenewInstanceWithContext(ctx volcengine.Context, input *RenewInstanceInput, opts ...request.Option) (*RenewInstanceOutput, error) {
	req, out := c.RenewInstanceRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type RenewInstanceInput struct {
	_ struct{} `type:"structure"`

	ClientToken *string `type:"string"`

	InstanceId *string `type:"string"`

	Period *int32 `type:"int32"`

	PeriodUnit *string `type:"string"`
}

// String returns the string representation
func (s RenewInstanceInput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s RenewInstanceInput) GoString() string {
	return s.String()
}

// SetClientToken sets the ClientToken field's value.
func (s *RenewInstanceInput) SetClientToken(v string) *RenewInstanceInput {
	s.ClientToken = &v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *RenewInstanceInput) SetInstanceId(v string) *RenewInstanceInput {
	s.InstanceId = &v
	return s
}

// SetPeriod sets the Period field's value.
func (s *RenewInstanceInput) SetPeriod(v int32) *RenewInstanceInput {
	s.Period = &v
	return s
}

// SetPeriodUnit sets the PeriodUnit field's value.
func (s *RenewInstanceInput) SetPeriodUnit(v string) *RenewInstanceInput {
	s.PeriodUnit = &v
	return s
}

type RenewInstanceOutput struct {
	_ struct{} `type:"structure"`

	Metadata *response.ResponseMetadata

	OrderId *string `type:"string"`
}

// String returns the string representation
func (s RenewInstanceOutput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s RenewInstanceOutput) GoString() string {
	return s.String()
}

// SetOrderId sets the OrderId field's value.
func (s *RenewInstanceOutput) SetOrderId(v string) *RenewInstanceOutput {
	s.OrderId = &v
	return s
}
