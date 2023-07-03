// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package kinesisvideosignalingchannels

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/private/protocol"
)

const (

	// ErrCodeClientLimitExceededException for service response error code
	// "ClientLimitExceededException".
	//
	// Your request was throttled because you have exceeded the limit of allowed
	// client calls. Try making the call later.
	ErrCodeClientLimitExceededException = "ClientLimitExceededException"

	// ErrCodeInvalidArgumentException for service response error code
	// "InvalidArgumentException".
	//
	// The value for this input parameter is invalid.
	ErrCodeInvalidArgumentException = "InvalidArgumentException"

	// ErrCodeInvalidClientException for service response error code
	// "InvalidClientException".
	//
	// The specified client is invalid.
	ErrCodeInvalidClientException = "InvalidClientException"

	// ErrCodeNotAuthorizedException for service response error code
	// "NotAuthorizedException".
	//
	// The caller is not authorized to perform this operation.
	ErrCodeNotAuthorizedException = "NotAuthorizedException"

	// ErrCodeResourceNotFoundException for service response error code
	// "ResourceNotFoundException".
	//
	// The specified resource is not found.
	ErrCodeResourceNotFoundException = "ResourceNotFoundException"

	// ErrCodeSessionExpiredException for service response error code
	// "SessionExpiredException".
	//
	// If the client session is expired. Once the client is connected, the session
	// is valid for 45 minutes. Client should reconnect to the channel to continue
	// sending/receiving messages.
	ErrCodeSessionExpiredException = "SessionExpiredException"
)

var exceptionFromCode = map[string]func(protocol.ResponseMetadata) error{
	"ClientLimitExceededException": newErrorClientLimitExceededException,
	"InvalidArgumentException":     newErrorInvalidArgumentException,
	"InvalidClientException":       newErrorInvalidClientException,
	"NotAuthorizedException":       newErrorNotAuthorizedException,
	"ResourceNotFoundException":    newErrorResourceNotFoundException,
	"SessionExpiredException":      newErrorSessionExpiredException,
}
