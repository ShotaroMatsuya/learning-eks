// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package sagemaker

import (
	"time"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/aws"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/aws/request"
)

// WaitUntilEndpointDeleted uses the SageMaker API operation
// DescribeEndpoint to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilEndpointDeleted(input *DescribeEndpointInput) error {
	return c.WaitUntilEndpointDeletedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilEndpointDeletedWithContext is an extended version of WaitUntilEndpointDeleted.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilEndpointDeletedWithContext(ctx aws.Context, input *DescribeEndpointInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilEndpointDeleted",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(30 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:    request.SuccessWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "EndpointStatus",
				Expected: "Failed",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeEndpointInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeEndpointRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilEndpointInService uses the SageMaker API operation
// DescribeEndpoint to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilEndpointInService(input *DescribeEndpointInput) error {
	return c.WaitUntilEndpointInServiceWithContext(aws.BackgroundContext(), input)
}

// WaitUntilEndpointInServiceWithContext is an extended version of WaitUntilEndpointInService.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilEndpointInServiceWithContext(ctx aws.Context, input *DescribeEndpointInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilEndpointInService",
		MaxAttempts: 120,
		Delay:       request.ConstantWaiterDelay(30 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "EndpointStatus",
				Expected: "InService",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "EndpointStatus",
				Expected: "Failed",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeEndpointInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeEndpointRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilImageCreated uses the SageMaker API operation
// DescribeImage to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilImageCreated(input *DescribeImageInput) error {
	return c.WaitUntilImageCreatedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilImageCreatedWithContext is an extended version of WaitUntilImageCreated.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilImageCreatedWithContext(ctx aws.Context, input *DescribeImageInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilImageCreated",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageStatus",
				Expected: "CREATED",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageStatus",
				Expected: "CREATE_FAILED",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeImageInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeImageRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilImageDeleted uses the SageMaker API operation
// DescribeImage to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilImageDeleted(input *DescribeImageInput) error {
	return c.WaitUntilImageDeletedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilImageDeletedWithContext is an extended version of WaitUntilImageDeleted.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilImageDeletedWithContext(ctx aws.Context, input *DescribeImageInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilImageDeleted",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:    request.SuccessWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ResourceNotFoundException",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageStatus",
				Expected: "DELETE_FAILED",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeImageInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeImageRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilImageUpdated uses the SageMaker API operation
// DescribeImage to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilImageUpdated(input *DescribeImageInput) error {
	return c.WaitUntilImageUpdatedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilImageUpdatedWithContext is an extended version of WaitUntilImageUpdated.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilImageUpdatedWithContext(ctx aws.Context, input *DescribeImageInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilImageUpdated",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageStatus",
				Expected: "CREATED",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageStatus",
				Expected: "UPDATE_FAILED",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeImageInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeImageRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilImageVersionCreated uses the SageMaker API operation
// DescribeImageVersion to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilImageVersionCreated(input *DescribeImageVersionInput) error {
	return c.WaitUntilImageVersionCreatedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilImageVersionCreatedWithContext is an extended version of WaitUntilImageVersionCreated.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilImageVersionCreatedWithContext(ctx aws.Context, input *DescribeImageVersionInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilImageVersionCreated",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageVersionStatus",
				Expected: "CREATED",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageVersionStatus",
				Expected: "CREATE_FAILED",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeImageVersionInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeImageVersionRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilImageVersionDeleted uses the SageMaker API operation
// DescribeImageVersion to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilImageVersionDeleted(input *DescribeImageVersionInput) error {
	return c.WaitUntilImageVersionDeletedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilImageVersionDeletedWithContext is an extended version of WaitUntilImageVersionDeleted.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilImageVersionDeletedWithContext(ctx aws.Context, input *DescribeImageVersionInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilImageVersionDeleted",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:    request.SuccessWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ResourceNotFoundException",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ImageVersionStatus",
				Expected: "DELETE_FAILED",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeImageVersionInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeImageVersionRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilNotebookInstanceDeleted uses the SageMaker API operation
// DescribeNotebookInstance to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilNotebookInstanceDeleted(input *DescribeNotebookInstanceInput) error {
	return c.WaitUntilNotebookInstanceDeletedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilNotebookInstanceDeletedWithContext is an extended version of WaitUntilNotebookInstanceDeleted.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilNotebookInstanceDeletedWithContext(ctx aws.Context, input *DescribeNotebookInstanceInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilNotebookInstanceDeleted",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(30 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:    request.SuccessWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "NotebookInstanceStatus",
				Expected: "Failed",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeNotebookInstanceInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeNotebookInstanceRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilNotebookInstanceInService uses the SageMaker API operation
// DescribeNotebookInstance to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilNotebookInstanceInService(input *DescribeNotebookInstanceInput) error {
	return c.WaitUntilNotebookInstanceInServiceWithContext(aws.BackgroundContext(), input)
}

// WaitUntilNotebookInstanceInServiceWithContext is an extended version of WaitUntilNotebookInstanceInService.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilNotebookInstanceInServiceWithContext(ctx aws.Context, input *DescribeNotebookInstanceInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilNotebookInstanceInService",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(30 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "NotebookInstanceStatus",
				Expected: "InService",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "NotebookInstanceStatus",
				Expected: "Failed",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeNotebookInstanceInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeNotebookInstanceRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilNotebookInstanceStopped uses the SageMaker API operation
// DescribeNotebookInstance to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilNotebookInstanceStopped(input *DescribeNotebookInstanceInput) error {
	return c.WaitUntilNotebookInstanceStoppedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilNotebookInstanceStoppedWithContext is an extended version of WaitUntilNotebookInstanceStopped.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilNotebookInstanceStoppedWithContext(ctx aws.Context, input *DescribeNotebookInstanceInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilNotebookInstanceStopped",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(30 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "NotebookInstanceStatus",
				Expected: "Stopped",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "NotebookInstanceStatus",
				Expected: "Failed",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeNotebookInstanceInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeNotebookInstanceRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilProcessingJobCompletedOrStopped uses the SageMaker API operation
// DescribeProcessingJob to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilProcessingJobCompletedOrStopped(input *DescribeProcessingJobInput) error {
	return c.WaitUntilProcessingJobCompletedOrStoppedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilProcessingJobCompletedOrStoppedWithContext is an extended version of WaitUntilProcessingJobCompletedOrStopped.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilProcessingJobCompletedOrStoppedWithContext(ctx aws.Context, input *DescribeProcessingJobInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilProcessingJobCompletedOrStopped",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ProcessingJobStatus",
				Expected: "Completed",
			},
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ProcessingJobStatus",
				Expected: "Stopped",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "ProcessingJobStatus",
				Expected: "Failed",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeProcessingJobInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeProcessingJobRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilTrainingJobCompletedOrStopped uses the SageMaker API operation
// DescribeTrainingJob to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilTrainingJobCompletedOrStopped(input *DescribeTrainingJobInput) error {
	return c.WaitUntilTrainingJobCompletedOrStoppedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilTrainingJobCompletedOrStoppedWithContext is an extended version of WaitUntilTrainingJobCompletedOrStopped.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilTrainingJobCompletedOrStoppedWithContext(ctx aws.Context, input *DescribeTrainingJobInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilTrainingJobCompletedOrStopped",
		MaxAttempts: 180,
		Delay:       request.ConstantWaiterDelay(120 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "TrainingJobStatus",
				Expected: "Completed",
			},
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "TrainingJobStatus",
				Expected: "Stopped",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "TrainingJobStatus",
				Expected: "Failed",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeTrainingJobInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeTrainingJobRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilTransformJobCompletedOrStopped uses the SageMaker API operation
// DescribeTransformJob to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *SageMaker) WaitUntilTransformJobCompletedOrStopped(input *DescribeTransformJobInput) error {
	return c.WaitUntilTransformJobCompletedOrStoppedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilTransformJobCompletedOrStoppedWithContext is an extended version of WaitUntilTransformJobCompletedOrStopped.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *SageMaker) WaitUntilTransformJobCompletedOrStoppedWithContext(ctx aws.Context, input *DescribeTransformJobInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilTransformJobCompletedOrStopped",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(60 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "TransformJobStatus",
				Expected: "Completed",
			},
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "TransformJobStatus",
				Expected: "Stopped",
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "TransformJobStatus",
				Expected: "Failed",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationException",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *DescribeTransformJobInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := c.DescribeTransformJobRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
