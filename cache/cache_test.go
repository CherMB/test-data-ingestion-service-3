package cache

import (
	"reflect"
	"testing"

	"github.com/calculi-corp/common/pkg/messaging"
	coredata "github.com/calculi-corp/core-data-cache"
	cmock "github.com/calculi-corp/core-data-cache/mock"
	client "github.com/calculi-corp/grpc-client"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSetMockCache(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResourceCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEndpointCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	SetMockCache(mockResourceCache, mockEndpointCache, mockGrpcClient)

	assert.Equal(t, mockResourceCache, CoreDataResourceCache, "CoreDataResourceCache should be set to mockResourceCache")
	assert.Equal(t, mockEndpointCache, epCache, "epCache should be set to mockEndpointCache")
	assert.Equal(t, mockGrpcClient, GrpcClient, "GrpcClient should be set to mockGrpcClient")
}

func TestInitializeCache(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	type args struct {
		clt    client.GrpcClient
		msgClt messaging.Messaging
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				clt:    mockGrpcClient,
				msgClt: nil,
			},
		},
	}
	mockGrpcClient.EXPECT().SendGrpc(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitializeCache(tt.args.clt, tt.args.msgClt); (err != nil) != tt.wantErr {
				assert.Equal(t, err.Error(), "config must contain a MessagingClient or MsgClient", "Validate message client nil")
			}
		})
	}
}

func TestGetGrpcClient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	GrpcClient = mockGrpcClient
	tests := []struct {
		name string
		want client.GrpcClient
	}{
		{
			name: "Validate grpc client",
			want: mockGrpcClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGrpcClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGrpcClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCoreDataCache(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	defer mockCtrl.Finish()
	mockResourceCache := cmock.NewMockResourceCacheI(mockCtrl)
	GrpcClient = mockGrpcClient
	tests := []struct {
		name string
		want coredata.ResourceCacheI
	}{
		{
			name: "Validate coredatacache",
			want: mockResourceCache,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCoreDataCache()
			assert.NotEqual(t, tt.want, got, "Unexpected cache instance")
		})
	}
}

func TestGetEndpointCache(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	defer mockCtrl.Finish()
	mockEndpointCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	GrpcClient = mockGrpcClient
	tests := []struct {
		name string
		want coredata.EndpointsCacheI
	}{
		{
			name: "Validate endpointcache",
			want: mockEndpointCache,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEndpointCache()
			assert.NotEqual(t, tt.want, got, "Unexpected cache instance")
		})
	}
}
