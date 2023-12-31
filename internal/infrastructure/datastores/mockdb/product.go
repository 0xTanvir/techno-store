// Code generated by MockGen. DO NOT EDIT.
// Source: techno-store/internal/domain/definition (interfaces: ProductRepository)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/product.go techno-store/internal/domain/definition ProductRepository
//
// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"
	bo "techno-store/internal/domain/bo"

	gomock "go.uber.org/mock/gomock"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// CreateProduct mocks base method.
func (m *MockProductRepository) CreateProduct(arg0 context.Context, arg1 *bo.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockProductRepositoryMockRecorder) CreateProduct(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockProductRepository)(nil).CreateProduct), arg0, arg1)
}

// DeleteProduct mocks base method.
func (m *MockProductRepository) DeleteProduct(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockProductRepositoryMockRecorder) DeleteProduct(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockProductRepository)(nil).DeleteProduct), arg0, arg1)
}

// GetProductByID mocks base method.
func (m *MockProductRepository) GetProductByID(arg0 context.Context, arg1 int64) (bo.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByID", arg0, arg1)
	ret0, _ := ret[0].(bo.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID.
func (mr *MockProductRepositoryMockRecorder) GetProductByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockProductRepository)(nil).GetProductByID), arg0, arg1)
}

// ListProducts mocks base method.
func (m *MockProductRepository) ListProducts(arg0 context.Context, arg1 bo.ProductSearchQuery) (bo.PaginatedProductCollection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProducts", arg0, arg1)
	ret0, _ := ret[0].(bo.PaginatedProductCollection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProducts indicates an expected call of ListProducts.
func (mr *MockProductRepositoryMockRecorder) ListProducts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProducts", reflect.TypeOf((*MockProductRepository)(nil).ListProducts), arg0, arg1)
}

// UpdateProduct mocks base method.
func (m *MockProductRepository) UpdateProduct(arg0 context.Context, arg1 bo.ProductUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockProductRepositoryMockRecorder) UpdateProduct(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockProductRepository)(nil).UpdateProduct), arg0, arg1)
}
