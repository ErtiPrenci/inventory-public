package service

import (
	"encoding/json"
	"fmt"
	"inventory-backend/internal/core"
	"inventory-backend/internal/repository"
	"inventory-backend/internal/utils/gen"
	"inventory-backend/internal/utils/response"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) RegisterRoutes(r chi.Router) {
	r.Get("/products", s.ListProducts)
	r.Post("/products", s.CreateProduct)
	r.Get("/products/{id}", s.GetProductByID)
	r.Put("/products/{id}", s.UpdateProduct)
	r.Delete("/products/{id}", s.DeleteProduct)
}

func (s *ProductService) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	filter := core.ProductFilter{
		Page:       page,
		Limit:      limit,
		Search:     search,
		CategoryID: categoryID,
	}

	// If pagination is requested (page > 0)
	if page > 0 || limit > 0 {
		if limit <= 0 {
			limit = 10 // Default limit
			filter.Limit = 10
		}
		if page <= 0 {
			page = 1
			filter.Page = 1
		}

		total, err := s.repo.Count(r.Context(), filter)
		if err != nil {
			response.InternalServerError(w, "Internal Server Error", err.Error())
			return
		}

		products, err := s.repo.GetAll(r.Context(), filter)
		if err != nil {
			response.InternalServerError(w, "Internal Server Error", err.Error())
			return
		}

		totalPages := 0
		if limit > 0 {
			totalPages = (total + limit - 1) / limit
		}

		meta := core.MetaData{
			Total:      total,
			Page:       page,
			PerPage:    limit,
			TotalPages: totalPages,
		}

		res := core.PaginatedResponse{
			Data: products,
			Meta: meta,
		}

		response.SuccessData(w, "Products List", "Products successfully obtained", res)
		return
	}

	// Backward compatibility: Return all if no pagination params
	// Pass empty filter which defaults to "no limit" in GetAll implementation if Limit is 0
	products, err := s.repo.GetAll(r.Context(), filter)
	if err != nil {
		response.InternalServerError(w, "Internal Server Error", err.Error())
		return
	}
	response.SuccessData(w, "Products List", "Products successfully obtained", products)
}

func (s *ProductService) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}
	product, err := s.repo.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Product not found")
		return
	}

	response.SuccessData(w, "Product Found", "Product successfully obtained", product)
}

func (s *ProductService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p core.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON format")
		return
	}

	// Business validation
	if p.SKU == "" || p.Name == "" || p.Price < 0 || p.CurrentStock < 0 || p.MinStock < 0 {
		response.BadRequest(w, "Bad Request", "SKU and Name are required")
		return
	}

	createdProduct, err := s.repo.Create(r.Context(), p)
	if err != nil {
		response.Conflict(w, "Conflict", "Possible SKU duplicate or DB error")
		return
	}

	response.SuccessData(w, "Product Created", "Saved in inventory", createdProduct)
}

func (s *ProductService) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}
	var p core.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON format")
		return
	}

	// Business validation
	if p.SKU == "" || p.Name == "" || p.Price < 0 || p.CurrentStock < 0 || p.MinStock < 0 {
		response.BadRequest(w, "Bad Request", "SKU and Name are required")
		return
	}

	p.ID = id
	updatedProduct, err := s.repo.Update(r.Context(), p)
	if err != nil {
		response.Conflict(w, "Conflict", "Possible SKU duplicate or DB error")
		return
	}
	response.SuccessData(w, "Product Updated", "Updated in inventory", updatedProduct)
}

func (s *ProductService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}
	err := s.repo.Delete(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Product not found")
		return
	}
	response.SuccessData(w, "Success", fmt.Sprintf("Product %s deleted successfully", id), nil)
}
