package models

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

// PermissionMock PermissionMock
type PermissionMock struct {
	PermissionID string
	mock.Mock
}

// Get Get
func (p PermissionMock) Get(ctx context.Context, perGroupID, formID string) (*Permission, error) {
	panic("implement me")
}

// Create Create
func (p PermissionMock) Create(ctx context.Context, permission *Permission, ttl time.Duration) error {
	panic("implement me")
}

// Delete Delete
func (p PermissionMock) Delete(ctx context.Context, permission *Permission) error {
	panic("implement me")
}

// CreatePerMatch CreatePerMatch
func (p PermissionMock) CreatePerMatch(ctx context.Context, match *PermissionMatch) error {
	panic("implement me")
}

// GetPerMatch GetPerMatch
func (p PermissionMock) GetPerMatch(ctx context.Context, userID, appID string) (*PermissionMatch, error) {
	panic("implement me")
}

// DeletePerMatch DeletePerMatch
func (p PermissionMock) DeletePerMatch(ctx context.Context, appID string) error {
	panic("implement me")
}

// Lock Lock
func (p PermissionMock) Lock(ctx context.Context, key string, val interface{}, ttl time.Duration) (bool, error) {
	panic("implement me")
}

// UnLock UnLock
func (p PermissionMock) UnLock(ctx context.Context, key string) error {
	panic("implement me")
}

// PerMatchExpire PerMatchExpire
func (p PermissionMock) PerMatchExpire(ctx context.Context, key string, ttl time.Duration) error {
	panic("implement me")
}

// PermissionMockOpt PermissionMockOpt
type PermissionMockOpt func(*PermissionMock)

// WithPermissionMockID WithPermissionMockID
func WithPermissionMockID(id string) PermissionMockOpt {
	return func(p *PermissionMock) {
		p.PermissionID = id
	}
}

//NewPerMissionMock NewPerMissionMock
func NewPerMissionMock(opts ...PermissionMockOpt) PermissionRepo {
	per := new(PermissionMock)
	for _, opt := range opts {
		opt(per)
	}
	per.On("Get", mock.Anything).Return(map[string]Permission{
		"1": {
			PerGroupID: "1",
			Name:       "测试",
			FormID:     "1",
			AppID:      "1",
		},
	})
	per.On("Create", mock.Anything).Return()
	per.On("Delete", mock.Anything).Return()
	per.On("CreatePerMatch", mock.Anything).Return()
	per.On("GetPerMatch", mock.Anything).Return()
	per.On("DeletePerMatch", mock.Anything).Return()
	per.On("Lock", mock.Anything).Return()
	per.On("UnLock", mock.Anything).Return()
	per.On("PerMatchExpire", mock.Anything).Return()
	return per
}
