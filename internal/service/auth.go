package service

import (
	"encoding/json"
	"inventory-backend/internal/repository"
	"inventory-backend/internal/utils/auth"
	"inventory-backend/internal/utils/response"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterRoutes(r chi.Router) {
	r.Post("/login", s.Login)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON")
		return
	}
	log.Println("Username: ", req.Username)
	user, err := s.repo.GetByUsername(r.Context(), req.Username)
	if err != nil {
		response.Unauthorized(w, "Unauthorized", "Invalid credentials")
		return
	}
	log.Println("CheckPasswordHash")
	match := auth.CheckPasswordHash(req.Password, user.PasswordHash)
	if !match {
		response.Unauthorized(w, "Unauthorized", "Invalid credentials")
		return
	}
	log.Println("GenerateToken")
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		response.InternalServerError(w, "Error", "Could not generate token")
		return
	}
	log.Println("Response.SuccessData")
	response.SuccessData(w, "Login Successful", "Token generated", map[string]string{
		"token":    token,
		"username": user.Username,
	})
}
