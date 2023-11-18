// Code generated by MockGen. DO NOT EDIT.
// Source: techno-store/internal/domain/definition (interfaces: BrandRepository)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/brand.go techno-store/internal/domain/definition BrandRepository
//
// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"
	bo "techno-store/internal/domain/bo"

	gomock "go.uber.org/mock/gomock"
)

// MockBrandRepository is a mock of BrandRepository interface.
type MockBrandRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBrandRepositoryMockRecorder
}

// MockBrandRepositoryMockRecorder is the mock recorder for MockBrandRepository.
type MockBrandRepositoryMockRecorder struct {
	mock *MockBrandRepository
}

// NewMockBrandRepository creates a new mock instance.
func NewMockBrandRepository(ctrl *gomock.Controller) *MockBrandRepository {
	mock := &MockBrandRepository{ctrl: ctrl}
	mock.recorder = &MockBrandRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBrandRepository) EXPECT() *MockBrandRepositoryMockRecorder {
	return m.recorder
}

// CreateBrand mocks base method.
func (m *MockBrandRepository) CreateBrand(arg0 context.Context, arg1 *bo.Brand) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBrand", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBrand indicates an expected call of CreateBrand.
func (mr *MockBrandRepositoryMockRecorder) CreateBrand(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBrand", reflect.TypeOf((*MockBrandRepository)(nil).CreateBrand), arg0, arg1)
}

// DeleteBrand mocks base method.
func (m *MockBrandRepository) DeleteBrand(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBrand", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBrand indicates an expected call of DeleteBrand.
func (mr *MockBrandRepositoryMockRecorder) DeleteBrand(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBrand", reflect.TypeOf((*MockBrandRepository)(nil).DeleteBrand), arg0, arg1)
}

// GetBrandByID mocks base method.
func (m *MockBrandRepository) GetBrandByID(arg0 context.Context, arg1 int64) (bo.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBrandByID", arg0, arg1)
	ret0, _ := ret[0].(bo.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBrandByID indicates an expected call of GetBrandByID.
func (mr *MockBrandRepositoryMockRecorder) GetBrandByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBrandByID", reflect.TypeOf((*MockBrandRepository)(nil).GetBrandByID), arg0, arg1)
}

// ListBrands mocks base method.
func (m *MockBrandRepository) ListBrands(arg0 context.Context, arg1 bo.BrandQuery) (bo.PaginatedBrandCollection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBrands", arg0, arg1)
	ret0, _ := ret[0].(bo.PaginatedBrandCollection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBrands indicates an expected call of ListBrands.
func (mr *MockBrandRepositoryMockRecorder) ListBrands(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBrands", reflect.TypeOf((*MockBrandRepository)(nil).ListBrands), arg0, arg1)
}

// UpdateBrand mocks base method.
func (m *MockBrandRepository) UpdateBrand(arg0 context.Context, arg1 bo.BrandUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBrand", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBrand indicates an expected call of UpdateBrand.
func (mr *MockBrandRepositoryMockRecorder) UpdateBrand(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBrand", reflect.TypeOf((*MockBrandRepository)(nil).UpdateBrand), arg0, arg1)
}