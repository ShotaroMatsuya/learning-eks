// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package memorydbiface provides an interface to enable mocking the Amazon MemoryDB service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package memorydbiface

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/aws"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/aws/request"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws/aws-sdk-go/service/memorydb"
)

// MemoryDBAPI provides an interface to enable mocking the
// memorydb.MemoryDB service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//	// myFunc uses an SDK service client to make a request to
//	// Amazon MemoryDB.
//	func myFunc(svc memorydbiface.MemoryDBAPI) bool {
//	    // Make svc.BatchUpdateCluster request
//	}
//
//	func main() {
//	    sess := session.New()
//	    svc := memorydb.New(sess)
//
//	    myFunc(svc)
//	}
//
// In your _test.go file:
//
//	// Define a mock struct to be used in your unit tests of myFunc.
//	type mockMemoryDBClient struct {
//	    memorydbiface.MemoryDBAPI
//	}
//	func (m *mockMemoryDBClient) BatchUpdateCluster(input *memorydb.BatchUpdateClusterInput) (*memorydb.BatchUpdateClusterOutput, error) {
//	    // mock response/functionality
//	}
//
//	func TestMyFunc(t *testing.T) {
//	    // Setup Test
//	    mockSvc := &mockMemoryDBClient{}
//
//	    myfunc(mockSvc)
//
//	    // Verify myFunc's functionality
//	}
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type MemoryDBAPI interface {
	BatchUpdateCluster(*memorydb.BatchUpdateClusterInput) (*memorydb.BatchUpdateClusterOutput, error)
	BatchUpdateClusterWithContext(aws.Context, *memorydb.BatchUpdateClusterInput, ...request.Option) (*memorydb.BatchUpdateClusterOutput, error)
	BatchUpdateClusterRequest(*memorydb.BatchUpdateClusterInput) (*request.Request, *memorydb.BatchUpdateClusterOutput)

	CopySnapshot(*memorydb.CopySnapshotInput) (*memorydb.CopySnapshotOutput, error)
	CopySnapshotWithContext(aws.Context, *memorydb.CopySnapshotInput, ...request.Option) (*memorydb.CopySnapshotOutput, error)
	CopySnapshotRequest(*memorydb.CopySnapshotInput) (*request.Request, *memorydb.CopySnapshotOutput)

	CreateACL(*memorydb.CreateACLInput) (*memorydb.CreateACLOutput, error)
	CreateACLWithContext(aws.Context, *memorydb.CreateACLInput, ...request.Option) (*memorydb.CreateACLOutput, error)
	CreateACLRequest(*memorydb.CreateACLInput) (*request.Request, *memorydb.CreateACLOutput)

	CreateCluster(*memorydb.CreateClusterInput) (*memorydb.CreateClusterOutput, error)
	CreateClusterWithContext(aws.Context, *memorydb.CreateClusterInput, ...request.Option) (*memorydb.CreateClusterOutput, error)
	CreateClusterRequest(*memorydb.CreateClusterInput) (*request.Request, *memorydb.CreateClusterOutput)

	CreateParameterGroup(*memorydb.CreateParameterGroupInput) (*memorydb.CreateParameterGroupOutput, error)
	CreateParameterGroupWithContext(aws.Context, *memorydb.CreateParameterGroupInput, ...request.Option) (*memorydb.CreateParameterGroupOutput, error)
	CreateParameterGroupRequest(*memorydb.CreateParameterGroupInput) (*request.Request, *memorydb.CreateParameterGroupOutput)

	CreateSnapshot(*memorydb.CreateSnapshotInput) (*memorydb.CreateSnapshotOutput, error)
	CreateSnapshotWithContext(aws.Context, *memorydb.CreateSnapshotInput, ...request.Option) (*memorydb.CreateSnapshotOutput, error)
	CreateSnapshotRequest(*memorydb.CreateSnapshotInput) (*request.Request, *memorydb.CreateSnapshotOutput)

	CreateSubnetGroup(*memorydb.CreateSubnetGroupInput) (*memorydb.CreateSubnetGroupOutput, error)
	CreateSubnetGroupWithContext(aws.Context, *memorydb.CreateSubnetGroupInput, ...request.Option) (*memorydb.CreateSubnetGroupOutput, error)
	CreateSubnetGroupRequest(*memorydb.CreateSubnetGroupInput) (*request.Request, *memorydb.CreateSubnetGroupOutput)

	CreateUser(*memorydb.CreateUserInput) (*memorydb.CreateUserOutput, error)
	CreateUserWithContext(aws.Context, *memorydb.CreateUserInput, ...request.Option) (*memorydb.CreateUserOutput, error)
	CreateUserRequest(*memorydb.CreateUserInput) (*request.Request, *memorydb.CreateUserOutput)

	DeleteACL(*memorydb.DeleteACLInput) (*memorydb.DeleteACLOutput, error)
	DeleteACLWithContext(aws.Context, *memorydb.DeleteACLInput, ...request.Option) (*memorydb.DeleteACLOutput, error)
	DeleteACLRequest(*memorydb.DeleteACLInput) (*request.Request, *memorydb.DeleteACLOutput)

	DeleteCluster(*memorydb.DeleteClusterInput) (*memorydb.DeleteClusterOutput, error)
	DeleteClusterWithContext(aws.Context, *memorydb.DeleteClusterInput, ...request.Option) (*memorydb.DeleteClusterOutput, error)
	DeleteClusterRequest(*memorydb.DeleteClusterInput) (*request.Request, *memorydb.DeleteClusterOutput)

	DeleteParameterGroup(*memorydb.DeleteParameterGroupInput) (*memorydb.DeleteParameterGroupOutput, error)
	DeleteParameterGroupWithContext(aws.Context, *memorydb.DeleteParameterGroupInput, ...request.Option) (*memorydb.DeleteParameterGroupOutput, error)
	DeleteParameterGroupRequest(*memorydb.DeleteParameterGroupInput) (*request.Request, *memorydb.DeleteParameterGroupOutput)

	DeleteSnapshot(*memorydb.DeleteSnapshotInput) (*memorydb.DeleteSnapshotOutput, error)
	DeleteSnapshotWithContext(aws.Context, *memorydb.DeleteSnapshotInput, ...request.Option) (*memorydb.DeleteSnapshotOutput, error)
	DeleteSnapshotRequest(*memorydb.DeleteSnapshotInput) (*request.Request, *memorydb.DeleteSnapshotOutput)

	DeleteSubnetGroup(*memorydb.DeleteSubnetGroupInput) (*memorydb.DeleteSubnetGroupOutput, error)
	DeleteSubnetGroupWithContext(aws.Context, *memorydb.DeleteSubnetGroupInput, ...request.Option) (*memorydb.DeleteSubnetGroupOutput, error)
	DeleteSubnetGroupRequest(*memorydb.DeleteSubnetGroupInput) (*request.Request, *memorydb.DeleteSubnetGroupOutput)

	DeleteUser(*memorydb.DeleteUserInput) (*memorydb.DeleteUserOutput, error)
	DeleteUserWithContext(aws.Context, *memorydb.DeleteUserInput, ...request.Option) (*memorydb.DeleteUserOutput, error)
	DeleteUserRequest(*memorydb.DeleteUserInput) (*request.Request, *memorydb.DeleteUserOutput)

	DescribeACLs(*memorydb.DescribeACLsInput) (*memorydb.DescribeACLsOutput, error)
	DescribeACLsWithContext(aws.Context, *memorydb.DescribeACLsInput, ...request.Option) (*memorydb.DescribeACLsOutput, error)
	DescribeACLsRequest(*memorydb.DescribeACLsInput) (*request.Request, *memorydb.DescribeACLsOutput)

	DescribeClusters(*memorydb.DescribeClustersInput) (*memorydb.DescribeClustersOutput, error)
	DescribeClustersWithContext(aws.Context, *memorydb.DescribeClustersInput, ...request.Option) (*memorydb.DescribeClustersOutput, error)
	DescribeClustersRequest(*memorydb.DescribeClustersInput) (*request.Request, *memorydb.DescribeClustersOutput)

	DescribeEngineVersions(*memorydb.DescribeEngineVersionsInput) (*memorydb.DescribeEngineVersionsOutput, error)
	DescribeEngineVersionsWithContext(aws.Context, *memorydb.DescribeEngineVersionsInput, ...request.Option) (*memorydb.DescribeEngineVersionsOutput, error)
	DescribeEngineVersionsRequest(*memorydb.DescribeEngineVersionsInput) (*request.Request, *memorydb.DescribeEngineVersionsOutput)

	DescribeEvents(*memorydb.DescribeEventsInput) (*memorydb.DescribeEventsOutput, error)
	DescribeEventsWithContext(aws.Context, *memorydb.DescribeEventsInput, ...request.Option) (*memorydb.DescribeEventsOutput, error)
	DescribeEventsRequest(*memorydb.DescribeEventsInput) (*request.Request, *memorydb.DescribeEventsOutput)

	DescribeParameterGroups(*memorydb.DescribeParameterGroupsInput) (*memorydb.DescribeParameterGroupsOutput, error)
	DescribeParameterGroupsWithContext(aws.Context, *memorydb.DescribeParameterGroupsInput, ...request.Option) (*memorydb.DescribeParameterGroupsOutput, error)
	DescribeParameterGroupsRequest(*memorydb.DescribeParameterGroupsInput) (*request.Request, *memorydb.DescribeParameterGroupsOutput)

	DescribeParameters(*memorydb.DescribeParametersInput) (*memorydb.DescribeParametersOutput, error)
	DescribeParametersWithContext(aws.Context, *memorydb.DescribeParametersInput, ...request.Option) (*memorydb.DescribeParametersOutput, error)
	DescribeParametersRequest(*memorydb.DescribeParametersInput) (*request.Request, *memorydb.DescribeParametersOutput)

	DescribeServiceUpdates(*memorydb.DescribeServiceUpdatesInput) (*memorydb.DescribeServiceUpdatesOutput, error)
	DescribeServiceUpdatesWithContext(aws.Context, *memorydb.DescribeServiceUpdatesInput, ...request.Option) (*memorydb.DescribeServiceUpdatesOutput, error)
	DescribeServiceUpdatesRequest(*memorydb.DescribeServiceUpdatesInput) (*request.Request, *memorydb.DescribeServiceUpdatesOutput)

	DescribeSnapshots(*memorydb.DescribeSnapshotsInput) (*memorydb.DescribeSnapshotsOutput, error)
	DescribeSnapshotsWithContext(aws.Context, *memorydb.DescribeSnapshotsInput, ...request.Option) (*memorydb.DescribeSnapshotsOutput, error)
	DescribeSnapshotsRequest(*memorydb.DescribeSnapshotsInput) (*request.Request, *memorydb.DescribeSnapshotsOutput)

	DescribeSubnetGroups(*memorydb.DescribeSubnetGroupsInput) (*memorydb.DescribeSubnetGroupsOutput, error)
	DescribeSubnetGroupsWithContext(aws.Context, *memorydb.DescribeSubnetGroupsInput, ...request.Option) (*memorydb.DescribeSubnetGroupsOutput, error)
	DescribeSubnetGroupsRequest(*memorydb.DescribeSubnetGroupsInput) (*request.Request, *memorydb.DescribeSubnetGroupsOutput)

	DescribeUsers(*memorydb.DescribeUsersInput) (*memorydb.DescribeUsersOutput, error)
	DescribeUsersWithContext(aws.Context, *memorydb.DescribeUsersInput, ...request.Option) (*memorydb.DescribeUsersOutput, error)
	DescribeUsersRequest(*memorydb.DescribeUsersInput) (*request.Request, *memorydb.DescribeUsersOutput)

	FailoverShard(*memorydb.FailoverShardInput) (*memorydb.FailoverShardOutput, error)
	FailoverShardWithContext(aws.Context, *memorydb.FailoverShardInput, ...request.Option) (*memorydb.FailoverShardOutput, error)
	FailoverShardRequest(*memorydb.FailoverShardInput) (*request.Request, *memorydb.FailoverShardOutput)

	ListAllowedNodeTypeUpdates(*memorydb.ListAllowedNodeTypeUpdatesInput) (*memorydb.ListAllowedNodeTypeUpdatesOutput, error)
	ListAllowedNodeTypeUpdatesWithContext(aws.Context, *memorydb.ListAllowedNodeTypeUpdatesInput, ...request.Option) (*memorydb.ListAllowedNodeTypeUpdatesOutput, error)
	ListAllowedNodeTypeUpdatesRequest(*memorydb.ListAllowedNodeTypeUpdatesInput) (*request.Request, *memorydb.ListAllowedNodeTypeUpdatesOutput)

	ListTags(*memorydb.ListTagsInput) (*memorydb.ListTagsOutput, error)
	ListTagsWithContext(aws.Context, *memorydb.ListTagsInput, ...request.Option) (*memorydb.ListTagsOutput, error)
	ListTagsRequest(*memorydb.ListTagsInput) (*request.Request, *memorydb.ListTagsOutput)

	ResetParameterGroup(*memorydb.ResetParameterGroupInput) (*memorydb.ResetParameterGroupOutput, error)
	ResetParameterGroupWithContext(aws.Context, *memorydb.ResetParameterGroupInput, ...request.Option) (*memorydb.ResetParameterGroupOutput, error)
	ResetParameterGroupRequest(*memorydb.ResetParameterGroupInput) (*request.Request, *memorydb.ResetParameterGroupOutput)

	TagResource(*memorydb.TagResourceInput) (*memorydb.TagResourceOutput, error)
	TagResourceWithContext(aws.Context, *memorydb.TagResourceInput, ...request.Option) (*memorydb.TagResourceOutput, error)
	TagResourceRequest(*memorydb.TagResourceInput) (*request.Request, *memorydb.TagResourceOutput)

	UntagResource(*memorydb.UntagResourceInput) (*memorydb.UntagResourceOutput, error)
	UntagResourceWithContext(aws.Context, *memorydb.UntagResourceInput, ...request.Option) (*memorydb.UntagResourceOutput, error)
	UntagResourceRequest(*memorydb.UntagResourceInput) (*request.Request, *memorydb.UntagResourceOutput)

	UpdateACL(*memorydb.UpdateACLInput) (*memorydb.UpdateACLOutput, error)
	UpdateACLWithContext(aws.Context, *memorydb.UpdateACLInput, ...request.Option) (*memorydb.UpdateACLOutput, error)
	UpdateACLRequest(*memorydb.UpdateACLInput) (*request.Request, *memorydb.UpdateACLOutput)

	UpdateCluster(*memorydb.UpdateClusterInput) (*memorydb.UpdateClusterOutput, error)
	UpdateClusterWithContext(aws.Context, *memorydb.UpdateClusterInput, ...request.Option) (*memorydb.UpdateClusterOutput, error)
	UpdateClusterRequest(*memorydb.UpdateClusterInput) (*request.Request, *memorydb.UpdateClusterOutput)

	UpdateParameterGroup(*memorydb.UpdateParameterGroupInput) (*memorydb.UpdateParameterGroupOutput, error)
	UpdateParameterGroupWithContext(aws.Context, *memorydb.UpdateParameterGroupInput, ...request.Option) (*memorydb.UpdateParameterGroupOutput, error)
	UpdateParameterGroupRequest(*memorydb.UpdateParameterGroupInput) (*request.Request, *memorydb.UpdateParameterGroupOutput)

	UpdateSubnetGroup(*memorydb.UpdateSubnetGroupInput) (*memorydb.UpdateSubnetGroupOutput, error)
	UpdateSubnetGroupWithContext(aws.Context, *memorydb.UpdateSubnetGroupInput, ...request.Option) (*memorydb.UpdateSubnetGroupOutput, error)
	UpdateSubnetGroupRequest(*memorydb.UpdateSubnetGroupInput) (*request.Request, *memorydb.UpdateSubnetGroupOutput)

	UpdateUser(*memorydb.UpdateUserInput) (*memorydb.UpdateUserOutput, error)
	UpdateUserWithContext(aws.Context, *memorydb.UpdateUserInput, ...request.Option) (*memorydb.UpdateUserOutput, error)
	UpdateUserRequest(*memorydb.UpdateUserInput) (*request.Request, *memorydb.UpdateUserOutput)
}

var _ MemoryDBAPI = (*memorydb.MemoryDB)(nil)
