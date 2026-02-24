package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"inventory-backend/internal/core"
	"inventory-backend/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll(ctx context.Context, filter core.ProductFilter) ([]core.Product, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]core.Product), args.Error(1)
}

func (m *MockProductRepository) Count(ctx context.Context, filter core.ProductFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (core.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(core.Product), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, p core.Product) (core.Product, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(core.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, p core.Product) (core.Product, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(core.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductService_CreateProduct(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	svc := service.NewProductService(mockRepo)

	// Test Data
	newProduct := core.Product{
		SKU:          "TEST-SKU-001",
		Name:         "Test Product",
		Description:  nil,
		Price:        100.0,
		Cost:         50.0,
		CurrentStock: 10,
		MinStock:     5,
		ImageURL:     nil,
	}

	createdProduct := newProduct
	createdProduct.ID = "generated-uuid"
	createdProduct.CreatedAt = time.Now()

	// Mock Expectation
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p core.Product) bool {
		return p.SKU == newProduct.SKU && p.Name == newProduct.Name
	})).Return(createdProduct, nil)

	// Request
	body, _ := json.Marshal(newProduct)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Handler
	handler := http.HandlerFunc(svc.CreateProduct)
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Product Created", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "generated-uuid", data["id"])
	assert.Equal(t, "TEST-SKU-001", data["sku"])

	mockRepo.AssertExpectations(t)
}

func TestProductService_ListProducts(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := service.NewProductService(mockRepo)

	expectedProducts := []core.Product{
		{ID: "p1", Name: "Product 1", Price: 10.0},
		{ID: "p2", Name: "Product 2", Price: 20.0},
	}

	mockRepo.On("GetAll", mock.Anything, core.ProductFilter{}).Return(expectedProducts, nil)

	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(svc.ListProducts)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "Products List", response["message"])

	// Check data type assertion carefully
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestProductService_GetProductByID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := service.NewProductService(mockRepo)

	// Valid UUID for checking
	id := "550e8400-e29b-41d4-a716-446655440000"
	expectedProduct := core.Product{ID: id, Name: "Found Me"}

	mockRepo.On("GetByID", mock.Anything, id).Return(expectedProduct, nil)

	req, _ := http.NewRequest("GET", "/products/"+id, nil)

	// Need to simulate chi context for URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(svc.GetProductByID)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, id, data["id"])
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := service.NewProductService(mockRepo)

	id := "550e8400-e29b-41d4-a716-446655440000"
	updateData := core.Product{
		SKU: "TEST-UPD", Name: "Updated Name", Price: 150.0, Cost: 10.0, CurrentStock: 5, MinStock: 1,
	}

	// Return the same product but with ID set
	returnedProduct := updateData
	returnedProduct.ID = id
	returnedProduct.UpdatedAt = time.Now()

	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(p core.Product) bool {
		return p.ID == id && p.Name == "Updated Name"
	})).Return(returnedProduct, nil)

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/products/"+id, bytes.NewBuffer(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(svc.UpdateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "Product Updated", response["message"])
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := service.NewProductService(mockRepo)

	id := "550e8400-e29b-41d4-a716-446655440000"

	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	req, _ := http.NewRequest("DELETE", "/products/"+id, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(svc.DeleteProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, fmt.Sprintf("Product %s deleted successfully", id), response["description"])
}
