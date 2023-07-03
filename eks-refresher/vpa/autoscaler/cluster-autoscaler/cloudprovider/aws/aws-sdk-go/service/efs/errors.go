// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package efs

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/private/protocol"
)

const (

	// ErrCodeAccessPointAlreadyExists for service response error code
	// "AccessPointAlreadyExists".
	//
	// Returned if the access point that you are trying to create already exists,
	// with the creation token you provided in the request.
	ErrCodeAccessPointAlreadyExists = "AccessPointAlreadyExists"

	// ErrCodeAccessPointLimitExceeded for service response error code
	// "AccessPointLimitExceeded".
	//
	// Returned if the Amazon Web Services account has already created the maximum
	// number of access points allowed per file system. For more informaton, see
	// https://docs.aws.amazon.com/efs/latest/ug/limits.html#limits-efs-resources-per-account-per-region
	// (https://docs.aws.amazon.com/efs/latest/ug/limits.html#limits-efs-resources-per-account-per-region).
	ErrCodeAccessPointLimitExceeded = "AccessPointLimitExceeded"

	// ErrCodeAccessPointNotFound for service response error code
	// "AccessPointNotFound".
	//
	// Returned if the specified AccessPointId value doesn't exist in the requester's
	// Amazon Web Services account.
	ErrCodeAccessPointNotFound = "AccessPointNotFound"

	// ErrCodeAvailabilityZonesMismatch for service response error code
	// "AvailabilityZonesMismatch".
	//
	// Returned if the Availability Zone that was specified for a mount target is
	// different from the Availability Zone that was specified for One Zone storage.
	// For more information, see Regional and One Zone storage redundancy (https://docs.aws.amazon.com/efs/latest/ug/availability-durability.html).
	ErrCodeAvailabilityZonesMismatch = "AvailabilityZonesMismatch"

	// ErrCodeBadRequest for service response error code
	// "BadRequest".
	//
	// Returned if the request is malformed or contains an error such as an invalid
	// parameter value or a missing required parameter.
	ErrCodeBadRequest = "BadRequest"

	// ErrCodeDependencyTimeout for service response error code
	// "DependencyTimeout".
	//
	// The service timed out trying to fulfill the request, and the client should
	// try the call again.
	ErrCodeDependencyTimeout = "DependencyTimeout"

	// ErrCodeFileSystemAlreadyExists for service response error code
	// "FileSystemAlreadyExists".
	//
	// Returned if the file system you are trying to create already exists, with
	// the creation token you provided.
	ErrCodeFileSystemAlreadyExists = "FileSystemAlreadyExists"

	// ErrCodeFileSystemInUse for service response error code
	// "FileSystemInUse".
	//
	// Returned if a file system has mount targets.
	ErrCodeFileSystemInUse = "FileSystemInUse"

	// ErrCodeFileSystemLimitExceeded for service response error code
	// "FileSystemLimitExceeded".
	//
	// Returned if the Amazon Web Services account has already created the maximum
	// number of file systems allowed per account.
	ErrCodeFileSystemLimitExceeded = "FileSystemLimitExceeded"

	// ErrCodeFileSystemNotFound for service response error code
	// "FileSystemNotFound".
	//
	// Returned if the specified FileSystemId value doesn't exist in the requester's
	// Amazon Web Services account.
	ErrCodeFileSystemNotFound = "FileSystemNotFound"

	// ErrCodeIncorrectFileSystemLifeCycleState for service response error code
	// "IncorrectFileSystemLifeCycleState".
	//
	// Returned if the file system's lifecycle state is not "available".
	ErrCodeIncorrectFileSystemLifeCycleState = "IncorrectFileSystemLifeCycleState"

	// ErrCodeIncorrectMountTargetState for service response error code
	// "IncorrectMountTargetState".
	//
	// Returned if the mount target is not in the correct state for the operation.
	ErrCodeIncorrectMountTargetState = "IncorrectMountTargetState"

	// ErrCodeInsufficientThroughputCapacity for service response error code
	// "InsufficientThroughputCapacity".
	//
	// Returned if there's not enough capacity to provision additional throughput.
	// This value might be returned when you try to create a file system in provisioned
	// throughput mode, when you attempt to increase the provisioned throughput
	// of an existing file system, or when you attempt to change an existing file
	// system from Bursting Throughput to Provisioned Throughput mode. Try again
	// later.
	ErrCodeInsufficientThroughputCapacity = "InsufficientThroughputCapacity"

	// ErrCodeInternalServerError for service response error code
	// "InternalServerError".
	//
	// Returned if an error occurred on the server side.
	ErrCodeInternalServerError = "InternalServerError"

	// ErrCodeInvalidPolicyException for service response error code
	// "InvalidPolicyException".
	//
	// Returned if the FileSystemPolicy is malformed or contains an error such as
	// a parameter value that is not valid or a missing required parameter. Returned
	// in the case of a policy lockout safety check error.
	ErrCodeInvalidPolicyException = "InvalidPolicyException"

	// ErrCodeIpAddressInUse for service response error code
	// "IpAddressInUse".
	//
	// Returned if the request specified an IpAddress that is already in use in
	// the subnet.
	ErrCodeIpAddressInUse = "IpAddressInUse"

	// ErrCodeMountTargetConflict for service response error code
	// "MountTargetConflict".
	//
	// Returned if the mount target would violate one of the specified restrictions
	// based on the file system's existing mount targets.
	ErrCodeMountTargetConflict = "MountTargetConflict"

	// ErrCodeMountTargetNotFound for service response error code
	// "MountTargetNotFound".
	//
	// Returned if there is no mount target with the specified ID found in the caller's
	// Amazon Web Services account.
	ErrCodeMountTargetNotFound = "MountTargetNotFound"

	// ErrCodeNetworkInterfaceLimitExceeded for service response error code
	// "NetworkInterfaceLimitExceeded".
	//
	// The calling account has reached the limit for elastic network interfaces
	// for the specific Amazon Web Services Region. Either delete some network interfaces
	// or request that the account quota be raised. For more information, see Amazon
	// VPC Quotas (https://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_Appendix_Limits.html)
	// in the Amazon VPC User Guide (see the Network interfaces per Region entry
	// in the Network interfaces table).
	ErrCodeNetworkInterfaceLimitExceeded = "NetworkInterfaceLimitExceeded"

	// ErrCodeNoFreeAddressesInSubnet for service response error code
	// "NoFreeAddressesInSubnet".
	//
	// Returned if IpAddress was not specified in the request and there are no free
	// IP addresses in the subnet.
	ErrCodeNoFreeAddressesInSubnet = "NoFreeAddressesInSubnet"

	// ErrCodePolicyNotFound for service response error code
	// "PolicyNotFound".
	//
	// Returned if the default file system policy is in effect for the EFS file
	// system specified.
	ErrCodePolicyNotFound = "PolicyNotFound"

	// ErrCodeReplicationNotFound for service response error code
	// "ReplicationNotFound".
	//
	// Returned if the specified file system does not have a replication configuration.
	ErrCodeReplicationNotFound = "ReplicationNotFound"

	// ErrCodeSecurityGroupLimitExceeded for service response error code
	// "SecurityGroupLimitExceeded".
	//
	// Returned if the size of SecurityGroups specified in the request is greater
	// than five.
	ErrCodeSecurityGroupLimitExceeded = "SecurityGroupLimitExceeded"

	// ErrCodeSecurityGroupNotFound for service response error code
	// "SecurityGroupNotFound".
	//
	// Returned if one of the specified security groups doesn't exist in the subnet's
	// virtual private cloud (VPC).
	ErrCodeSecurityGroupNotFound = "SecurityGroupNotFound"

	// ErrCodeSubnetNotFound for service response error code
	// "SubnetNotFound".
	//
	// Returned if there is no subnet with ID SubnetId provided in the request.
	ErrCodeSubnetNotFound = "SubnetNotFound"

	// ErrCodeThrottlingException for service response error code
	// "ThrottlingException".
	//
	// Returned when the CreateAccessPoint API action is called too quickly and
	// the number of Access Points in the account is nearing the limit of 120.
	ErrCodeThrottlingException = "ThrottlingException"

	// ErrCodeThroughputLimitExceeded for service response error code
	// "ThroughputLimitExceeded".
	//
	// Returned if the throughput mode or amount of provisioned throughput can't
	// be changed because the throughput limit of 1024 MiB/s has been reached.
	ErrCodeThroughputLimitExceeded = "ThroughputLimitExceeded"

	// ErrCodeTooManyRequests for service response error code
	// "TooManyRequests".
	//
	// Returned if you don’t wait at least 24 hours before either changing the
	// throughput mode, or decreasing the Provisioned Throughput value.
	ErrCodeTooManyRequests = "TooManyRequests"

	// ErrCodeUnsupportedAvailabilityZone for service response error code
	// "UnsupportedAvailabilityZone".
	//
	// Returned if the requested Amazon EFS functionality is not available in the
	// specified Availability Zone.
	ErrCodeUnsupportedAvailabilityZone = "UnsupportedAvailabilityZone"

	// ErrCodeValidationException for service response error code
	// "ValidationException".
	//
	// Returned if the Backup service is not available in the Amazon Web Services
	// Region in which the request was made.
	ErrCodeValidationException = "ValidationException"
)

var exceptionFromCode = map[string]func(protocol.ResponseMetadata) error{
	"AccessPointAlreadyExists":          newErrorAccessPointAlreadyExists,
	"AccessPointLimitExceeded":          newErrorAccessPointLimitExceeded,
	"AccessPointNotFound":               newErrorAccessPointNotFound,
	"AvailabilityZonesMismatch":         newErrorAvailabilityZonesMismatch,
	"BadRequest":                        newErrorBadRequest,
	"DependencyTimeout":                 newErrorDependencyTimeout,
	"FileSystemAlreadyExists":           newErrorFileSystemAlreadyExists,
	"FileSystemInUse":                   newErrorFileSystemInUse,
	"FileSystemLimitExceeded":           newErrorFileSystemLimitExceeded,
	"FileSystemNotFound":                newErrorFileSystemNotFound,
	"IncorrectFileSystemLifeCycleState": newErrorIncorrectFileSystemLifeCycleState,
	"IncorrectMountTargetState":         newErrorIncorrectMountTargetState,
	"InsufficientThroughputCapacity":    newErrorInsufficientThroughputCapacity,
	"InternalServerError":               newErrorInternalServerError,
	"InvalidPolicyException":            newErrorInvalidPolicyException,
	"IpAddressInUse":                    newErrorIpAddressInUse,
	"MountTargetConflict":               newErrorMountTargetConflict,
	"MountTargetNotFound":               newErrorMountTargetNotFound,
	"NetworkInterfaceLimitExceeded":     newErrorNetworkInterfaceLimitExceeded,
	"NoFreeAddressesInSubnet":           newErrorNoFreeAddressesInSubnet,
	"PolicyNotFound":                    newErrorPolicyNotFound,
	"ReplicationNotFound":               newErrorReplicationNotFound,
	"SecurityGroupLimitExceeded":        newErrorSecurityGroupLimitExceeded,
	"SecurityGroupNotFound":             newErrorSecurityGroupNotFound,
	"SubnetNotFound":                    newErrorSubnetNotFound,
	"ThrottlingException":               newErrorThrottlingException,
	"ThroughputLimitExceeded":           newErrorThroughputLimitExceeded,
	"TooManyRequests":                   newErrorTooManyRequests,
	"UnsupportedAvailabilityZone":       newErrorUnsupportedAvailabilityZone,
	"ValidationException":               newErrorValidationException,
}
