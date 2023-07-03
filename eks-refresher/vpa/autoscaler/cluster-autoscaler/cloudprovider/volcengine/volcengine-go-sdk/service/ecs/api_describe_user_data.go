// Code generated by volcengine with private/model/cli/gen-api/main.go. DO NOT EDIT.

package ecs

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/request"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/response"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/volcengineutil"
)

const opDescribeUserDataCommon = "DescribeUserData"

// DescribeUserDataCommonRequest generates a "volcengine/request.Request" representing the
// client's request for the DescribeUserDataCommon operation. The "output" return
// value will be populated with the DescribeUserDataCommon request's response once the request completes
// successfully.
//
// Use "Send" method on the returned DescribeUserDataCommon Request to send the API call to the service.
// the "output" return value is not valid until after DescribeUserDataCommon Send returns without error.
//
// See DescribeUserDataCommon for more information on using the DescribeUserDataCommon
// API call, and error handling.
//
//	// Example sending a request using the DescribeUserDataCommonRequest method.
//	req, resp := client.DescribeUserDataCommonRequest(params)
//
//	err := req.Send()
//	if err == nil { // resp is now filled
//	    fmt.Println(resp)
//	}
func (c *ECS) DescribeUserDataCommonRequest(input *map[string]interface{}) (req *request.Request, output *map[string]interface{}) {
	op := &request.Operation{
		Name:       opDescribeUserDataCommon,
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

// DescribeUserDataCommon API operation for ECS.
//
// Returns volcengineerr.Error for service API and SDK errors. Use runtime type assertions
// with volcengineerr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the VOLCENGINE API reference guide for ECS's
// API operation DescribeUserDataCommon for usage and error information.
func (c *ECS) DescribeUserDataCommon(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := c.DescribeUserDataCommonRequest(input)
	return out, req.Send()
}

// DescribeUserDataCommonWithContext is the same as DescribeUserDataCommon with the addition of
// the ability to pass a context and additional request options.
//
// See DescribeUserDataCommon for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If the context is nil a panic will occur.
// In the future the SDK may create sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ECS) DescribeUserDataCommonWithContext(ctx volcengine.Context, input *map[string]interface{}, opts ...request.Option) (*map[string]interface{}, error) {
	req, out := c.DescribeUserDataCommonRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opDescribeUserData = "DescribeUserData"

// DescribeUserDataRequest generates a "volcengine/request.Request" representing the
// client's request for the DescribeUserData operation. The "output" return
// value will be populated with the DescribeUserDataCommon request's response once the request completes
// successfully.
//
// Use "Send" method on the returned DescribeUserDataCommon Request to send the API call to the service.
// the "output" return value is not valid until after DescribeUserDataCommon Send returns without error.
//
// See DescribeUserData for more information on using the DescribeUserData
// API call, and error handling.
//
//	// Example sending a request using the DescribeUserDataRequest method.
//	req, resp := client.DescribeUserDataRequest(params)
//
//	err := req.Send()
//	if err == nil { // resp is now filled
//	    fmt.Println(resp)
//	}
func (c *ECS) DescribeUserDataRequest(input *DescribeUserDataInput) (req *request.Request, output *DescribeUserDataOutput) {
	op := &request.Operation{
		Name:       opDescribeUserData,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DescribeUserDataInput{}
	}

	output = &DescribeUserDataOutput{}
	req = c.newRequest(op, input, output)

	return
}

// DescribeUserData API operation for ECS.
//
// Returns volcengineerr.Error for service API and SDK errors. Use runtime type assertions
// with volcengineerr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the VOLCENGINE API reference guide for ECS's
// API operation DescribeUserData for usage and error information.
func (c *ECS) DescribeUserData(input *DescribeUserDataInput) (*DescribeUserDataOutput, error) {
	req, out := c.DescribeUserDataRequest(input)
	return out, req.Send()
}

// DescribeUserDataWithContext is the same as DescribeUserData with the addition of
// the ability to pass a context and additional request options.
//
// See DescribeUserData for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. Ifthe context is nil a panic will occur.
// In the future the SDK may create sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ECS) DescribeUserDataWithContext(ctx volcengine.Context, input *DescribeUserDataInput, opts ...request.Option) (*DescribeUserDataOutput, error) {
	req, out := c.DescribeUserDataRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type DescribeUserDataInput struct {
	_ struct{} `type:"structure"`

	InstanceId *string `type:"string"`
}

// String returns the string representation
func (s DescribeUserDataInput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeUserDataInput) GoString() string {
	return s.String()
}

// SetInstanceId sets the InstanceId field's value.
func (s *DescribeUserDataInput) SetInstanceId(v string) *DescribeUserDataInput {
	s.InstanceId = &v
	return s
}

type DescribeUserDataOutput struct {
	_ struct{} `type:"structure"`

	Metadata *response.ResponseMetadata

	InstanceId *string `type:"string"`

	UserData *string `type:"string"`
}

// String returns the string representation
func (s DescribeUserDataOutput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeUserDataOutput) GoString() string {
	return s.String()
}

// SetInstanceId sets the InstanceId field's value.
func (s *DescribeUserDataOutput) SetInstanceId(v string) *DescribeUserDataOutput {
	s.InstanceId = &v
	return s
}

// SetUserData sets the UserData field's value.
func (s *DescribeUserDataOutput) SetUserData(v string) *DescribeUserDataOutput {
	s.UserData = &v
	return s
}
