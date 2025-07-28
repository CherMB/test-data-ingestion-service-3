package mock

//go:generate mockgen -destination=./rbac_svc_mocks.go -package=mock github.com/calculi-corp/api/go/auth RBACServiceClient
//go:generate mockgen -destination=./org_svc_mocks.go -package=mock github.com/calculi-corp/api/go/auth OrganizationsServiceClient
