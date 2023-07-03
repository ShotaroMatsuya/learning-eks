// Code generated by volcengine with private/model/cli/gen-api/main.go. DO NOT EDIT.

package ecs

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/request"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/response"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/volcengine/volcengine-go-sdk/volcengine/volcengineutil"
)

const opDeleteInstancesCommon = "DeleteInstances"

// DeleteInstancesCommonRequest generates a "volcengine/request.Request" representing the
// client's request for the DeleteInstancesCommon operation. The "output" return
// value will be populated with the DeleteInstancesCommon request's response once the request completes
// successfully.
//
// Use "Send" method on the returned DeleteInstancesCommon Request to send the API call to the service.
// the "output" return value is not valid until after DeleteInstancesCommon Send returns without error.
//
// See DeleteInstancesCommon for more information on using the DeleteInstancesCommon
// API call, and error handling.
//
//	// Example sending a request using the DeleteInstancesCommonRequest method.
//	req, resp := client.DeleteInstancesCommonRequest(params)
//
//	err := req.Send()
//	if err == nil { // resp is now filled
//	    fmt.Println(resp)
//	}
func (c *ECS) DeleteInstancesCommonRequest(input *map[string]interface{}) (req *request.Request, output *map[string]interface{}) {
	op := &request.Operation{
		Name:       opDeleteInstancesCommon,
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

// DeleteInstancesCommon API operation for ECS.
//
// Returns volcengineerr.Error for service API and SDK errors. Use runtime type assertions
// with volcengineerr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the VOLCENGINE API reference guide for ECS's
// API operation DeleteInstancesCommon for usage and error information.
func (c *ECS) DeleteInstancesCommon(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := c.DeleteInstancesCommonRequest(input)
	return out, req.Send()
}

// DeleteInstancesCommonWithContext is the same as DeleteInstancesCommon with the addition of
// the ability to pass a context and additional request options.
//
// See DeleteInstancesCommon for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If the context is nil a panic will occur.
// In the future the SDK may create sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ECS) DeleteInstancesCommonWithContext(ctx volcengine.Context, input *map[string]interface{}, opts ...request.Option) (*map[string]interface{}, error) {
	req, out := c.DeleteInstancesCommonRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opDeleteInstances = "DeleteInstances"

// DeleteInstancesRequest generates a "volcengine/request.Request" representing the
// client's request for the DeleteInstances operation. The "output" return
// value will be populated with the DeleteInstancesCommon request's response once the request completes
// successfully.
//
// Use "Send" method on the returned DeleteInstancesCommon Request to send the API call to the service.
// the "output" return value is not valid until after DeleteInstancesCommon Send returns without error.
//
// See DeleteInstances for more information on using the DeleteInstances
// API call, and error handling.
//
//	// Example sending a request using the DeleteInstancesRequest method.
//	req, resp := client.DeleteInstancesRequest(params)
//
//	err := req.Send()
//	if err == nil { // resp is now filled
//	    fmt.Println(resp)
//	}
func (c *ECS) DeleteInstancesRequest(input *DeleteInstancesInput) (req *request.Request, output *DeleteInstancesOutput) {
	op := &request.Operation{
		Name:       opDeleteInstances,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DeleteInstancesInput{}
	}

	output = &DeleteInstancesOutput{}
	req = c.newRequest(op, input, output)

	return
}

// DeleteInstances API operation for ECS.
//
// Returns volcengineerr.Error for service API and SDK errors. Use runtime type assertions
// with volcengineerr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the VOLCENGINE API reference guide for ECS's
// API operation DeleteInstances for usage and error information.
func (c *ECS) DeleteInstances(input *DeleteInstancesInput) (*DeleteInstancesOutput, error) {
	req, out := c.DeleteInstancesRequest(input)
	return out, req.Send()
}

// DeleteInstancesWithContext is the same as DeleteInstances with the addition of
// the ability to pass a context and additional request options.
//
// See DeleteInstances for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. Ifthe context is nil a panic will occur.
// In the future the SDK may create sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ECS) DeleteInstancesWithContext(ctx volcengine.Context, input *DeleteInstancesInput, opts ...request.Option) (*DeleteInstancesOutput, error) {
	req, out := c.DeleteInstancesRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type DeleteInstancesInput struct {
	_ struct{} `type:"structure"`

	InstanceIds []*string `type:"list"`
}

// String returns the string representation
func (s DeleteInstancesInput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteInstancesInput) GoString() string {
	return s.String()
}

// SetInstanceIds sets the InstanceIds field's value.
func (s *DeleteInstancesInput) SetInstanceIds(v []*string) *DeleteInstancesInput {
	s.InstanceIds = v
	return s
}

type DeleteInstancesOutput struct {
	_ struct{} `type:"structure"`

	Metadata *response.ResponseMetadata

	OperationDetails []*OperationDetailForDeleteInstancesOutput `type:"list"`
}

// String returns the string representation
func (s DeleteInstancesOutput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteInstancesOutput) GoString() string {
	return s.String()
}

// SetOperationDetails sets the OperationDetails field's value.
func (s *DeleteInstancesOutput) SetOperationDetails(v []*OperationDetailForDeleteInstancesOutput) *DeleteInstancesOutput {
	s.OperationDetails = v
	return s
}

type ErrorForDeleteInstancesOutput struct {
	_ struct{} `type:"structure"`

	Code *string `type:"string"`

	Message *string `type:"string"`
}

// String returns the string representation
func (s ErrorForDeleteInstancesOutput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s ErrorForDeleteInstancesOutput) GoString() string {
	return s.String()
}

// SetCode sets the Code field's value.
func (s *ErrorForDeleteInstancesOutput) SetCode(v string) *ErrorForDeleteInstancesOutput {
	s.Code = &v
	return s
}

// SetMessage sets the Message field's value.
func (s *ErrorForDeleteInstancesOutput) SetMessage(v string) *ErrorForDeleteInstancesOutput {
	s.Message = &v
	return s
}

type OperationDetailForDeleteInstancesOutput struct {
	_ struct{} `type:"structure"`

	Error *ErrorForDeleteInstancesOutput `type:"structure"`

	InstanceId *string `type:"string"`
}

// String returns the string representation
func (s OperationDetailForDeleteInstancesOutput) String() string {
	return volcengineutil.Prettify(s)
}

// GoString returns the string representation
func (s OperationDetailForDeleteInstancesOutput) GoString() string {
	return s.String()
}

// SetError sets the Error field's value.
func (s *OperationDetailForDeleteInstancesOutput) SetError(v *ErrorForDeleteInstancesOutput) *OperationDetailForDeleteInstancesOutput {
	s.Error = v
	return s
}

// SetInstanceId sets the InstanceId field's value.
func (s *OperationDetailForDeleteInstancesOutput) SetInstanceId(v string) *OperationDetailForDeleteInstancesOutput {
	s.InstanceId = &v
	return s
}
