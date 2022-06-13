package resource

import (
	"context"
	"sync"
)

import (
	"github.com/seata/seata-go/pkg/protocol/branch"
)

// Resource that can be managed by Resource Manager and involved into global transaction
type Resource interface {
	GetResourceGroupId() string
	GetResourceId() string
	GetBranchType() branch.BranchType
}

// Control a branch transaction commit or rollback
type ResourceManagerInbound interface {
	// Commit a branch transaction
	BranchCommit(ctx context.Context, branchType branch.BranchType, xid string, branchId int64, resourceId string, applicationData []byte) (branch.BranchStatus, error)
	// Rollback a branch transaction
	BranchRollback(ctx context.Context, ranchType branch.BranchType, xid string, branchId int64, resourceId string, applicationData []byte) (branch.BranchStatus, error)
}

// Resource Manager: send outbound request to TC
type ResourceManagerOutbound interface {
	// Branch register long
	BranchRegister(ctx context.Context, ranchType branch.BranchType, resourceId, clientId, xid, applicationData, lockKeys string) (int64, error)
	//  Branch report
	BranchReport(ctx context.Context, ranchType branch.BranchType, xid string, branchId int64, status branch.BranchStatus, applicationData string) error
	// Lock query boolean
	LockQuery(ctx context.Context, ranchType branch.BranchType, resourceId, xid, lockKeys string) (bool, error)
}

//  Resource Manager: common behaviors
type ResourceManager interface {
	ResourceManagerInbound
	ResourceManagerOutbound

	// Register a Resource to be managed by Resource Manager
	RegisterResource(resource Resource) error
	//  Unregister a Resource from the Resource Manager
	UnregisterResource(resource Resource) error
	// Get all resources managed by this manager
	GetManagedResources() sync.Map
	// Get the BranchType
	GetBranchType() branch.BranchType
}

type ResourceManagerGetter interface {
	GetResourceManager() ResourceManager
}
