package service

import (
	"encoding/json"
	"fmt"
	"inventory-backend/internal/core"
	"inventory-backend/internal/repository"
	"inventory-backend/internal/utils/gen"
	"inventory-backend/internal/utils/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) RegisterRoutes(r chi.Router) {
	r.Get("/categories", s.ListCategories)
	r.Post("/categories", s.CreateCategory)
	r.Get("/categories/{id}", s.GetCategoryByID)
	r.Put("/categories/{id}", s.UpdateCategory)
	r.Delete("/categories/{id}", s.DeleteCategory)
}

func (s *CategoryService) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := s.repo.GetAll(r.Context())

	if err != nil {
		response.InternalServerError(w, "Internal Server Error", err.Error())
		return
	}

	response.SuccessData(w, "Categories List", "Categories successfully obtained", categories)
}

func (s *CategoryService) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}

	category, err := s.repo.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Category not found")
		return
	}

	response.SuccessData(w, "Category Found", "Category successfully obtained", category)
}

func (s *CategoryService) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var c core.Category

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON format")
		return
	}

	createdCategory, err := s.repo.Create(r.Context(), c)
	if err != nil {
		response.Conflict(w, "Conflict", "Possible category duplicate or DB error")
		return
	}

	response.SuccessData(w, "Category Created", "Saved in inventory", createdCategory)
}

func (s *CategoryService) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}

	var c core.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON format")
		return
	}

	c.ID = id
	updatedCategory, err := s.repo.Update(r.Context(), c)
	if err != nil {
		response.Conflict(w, "Conflict", "Possible category duplicate or DB error")
		return
	}

	response.SuccessData(w, "Category Updated", "Updated in inventory", updatedCategory)
}

func (s *CategoryService) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}

	err := s.repo.Delete(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Category not found")
		return
	}

	response.SuccessData(w, "Success", fmt.Sprintf("Category %s deleted successfully", id), nil)
}
