// Code generated by MockGen. DO NOT EDIT.
// Source: techno-store/internal/domain/definition (interfaces: SupplierRepository)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/supplier.go techno-store/internal/domain/definition SupplierRepository
//
// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"
	bo "techno-store/internal/domain/bo"

	gomock "go.uber.org/mock/gomock"
)

// MockSupplierRepository is a mock of SupplierRepository interface.
type MockSupplierRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSupplierRepositoryMockRecorder
}

// MockSupplierRepositoryMockRecorder is the mock recorder for MockSupplierRepository.
type MockSupplierRepositoryMockRecorder struct {
	mock *MockSupplierRepository
}

// NewMockSupplierRepository creates a new mock instance.
func NewMockSupplierRepository(ctrl *gomock.Controller) *MockSupplierRepository {
	mock := &MockSupplierRepository{ctrl: ctrl}
	mock.recorder = &MockSupplierRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSupplierRepository) EXPECT() *MockSupplierRepositoryMockRecorder {
	return m.recorder
}

// CreateSupplier mocks base method.
func (m *MockSupplierRepository) CreateSupplier(arg0 context.Context, arg1 *bo.Supplier) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSupplier", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSupplier indicates an expected call of CreateSupplier.
func (mr *MockSupplierRepositoryMockRecorder) CreateSupplier(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSupplier", reflect.TypeOf((*MockSupplierRepository)(nil).CreateSupplier), arg0, arg1)
}

// DeleteSupplier mocks base method.
func (m *MockSupplierRepository) DeleteSupplier(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSupplier", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSupplier indicates an expected call of DeleteSupplier.
func (mr *MockSupplierRepositoryMockRecorder) DeleteSupplier(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSupplier", reflect.TypeOf((*MockSupplierRepository)(nil).DeleteSupplier), arg0, arg1)
}

// GetSupplierByID mocks base method.
func (m *MockSupplierRepository) GetSupplierByID(arg0 context.Context, arg1 int64) (bo.Supplier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupplierByID", arg0, arg1)
	ret0, _ := ret[0].(bo.Supplier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupplierByID indicates an expected call of GetSupplierByID.
func (mr *MockSupplierRepositoryMockRecorder) GetSupplierByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupplierByID", reflect.TypeOf((*MockSupplierRepository)(nil).GetSupplierByID), arg0, arg1)
}

// ListSuppliers mocks base method.
func (m *MockSupplierRepository) ListSuppliers(arg0 context.Context, arg1 bo.SupplierQuery) (bo.PaginatedSupplierCollection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSuppliers", arg0, arg1)
	ret0, _ := ret[0].(bo.PaginatedSupplierCollection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSuppliers indicates an expected call of ListSuppliers.
func (mr *MockSupplierRepositoryMockRecorder) ListSuppliers(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSuppliers", reflect.TypeOf((*MockSupplierRepository)(nil).ListSuppliers), arg0, arg1)
}

// UpdateSupplier mocks base method.
func (m *MockSupplierRepository) UpdateSupplier(arg0 context.Context, arg1 bo.SupplierUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSupplier", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSupplier indicates an expected call of UpdateSupplier.
func (mr *MockSupplierRepositoryMockRecorder) UpdateSupplier(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSupplier", reflect.TypeOf((*MockSupplierRepository)(nil).UpdateSupplier), arg0, arg1)
}
