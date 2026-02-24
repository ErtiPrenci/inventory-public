package main

import (
	"context"
	"embed"
	"inventory-backend/internal/database"
	"inventory-backend/internal/repository"
	"inventory-backend/internal/service"
	"log"
	"net/http"
	"os"

	internalMiddleware "inventory-backend/internal/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

//go:embed assets/*
var assets embed.FS

func main() {
	//Load enviroment variables (only if local)
	if os.Getenv("LAMBDA_FUNCTION_NAME") == "" {
		_ = godotenv.Load()
		log.Println("Loaded enviroment variables")
	}

	companyName := os.Getenv("COMPANY_NAME")
	if companyName == "" {
		companyName = "Sin Nombre"
	}

	//Initialize Database
	dbURL := os.Getenv("DATABASE_URL")
	log.Println("DATABASE_URL: ", dbURL)
	if dbURL == "" {
		log.Fatal("URL from DB is required")
	}
	database.InitDB(dbURL)

	//Initialize Document Service
	docService := service.NewDocumentService(assets, "assets/logo.png", companyName)

	//Configure Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Log request body for debugging
	r.Use(internalMiddleware.LogRequestBody)

	var origins []string
	if os.Getenv("FRONTEND_URL") != "" {
		url_origin := os.Getenv("FRONTEND_URL")
		log.Println("FRONTEND_URL: ", url_origin)
		origins = []string{url_origin}
	} else {
		log.Println("FRONTEND_URL is not set, using default origins")
		origins = []string{"https://*", "http://*"}
	}

	// Basic CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link", "Content-Disposition"},

		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	categoryRepo := repository.NewCategoryRepository(database.DB)
	categoryService := service.NewCategoryService(categoryRepo)

	productRepo := repository.NewProductRepository(database.DB)
	productService := service.NewProductService(productRepo)

	orderRepo := repository.NewOrderRepository(database.DB)
	orderService := service.NewOrderService(orderRepo, productRepo, docService)

	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)

	r.Route("/v1", func(r chi.Router) {
		// Public routes
		log.Println("Registering public routes")
		authService.RegisterRoutes(r)

		//Protected routes
		r.Group(func(r chi.Router) {
			r.Use(internalMiddleware.AuthMiddleware)
			productService.RegisterRoutes(r)
			categoryService.RegisterRoutes(r)
			orderService.RegisterRoutes(r)

		})

	})

	//Decision of start: Lambda or Local?
	if os.Getenv("LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Starting in Lambda mode (HTTP API v2)")
		adapter := httpadapter.NewV2(r)

		lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			// Log detallado para ver la estructura real
			log.Printf("EVENTO RECIBIDO - RawPath: [%s], Method: [%s], Path: [%s]",
				req.RawPath,
				req.RequestContext.HTTP.Method,
				req.RequestContext.HTTP.Path)

			// Si RawPath está vacío, intentamos usar el Path del contexto
			if req.RawPath == "" && req.RequestContext.HTTP.Path != "" {
				req.RawPath = req.RequestContext.HTTP.Path
			}

			return adapter.ProxyWithContext(ctx, req)
		})
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Starting local server on http://localhost:%s", port)
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}
