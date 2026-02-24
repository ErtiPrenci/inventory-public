package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inventory-backend/internal/core"
	"inventory-backend/internal/repository"
	"inventory-backend/internal/utils/gen"
	"inventory-backend/internal/utils/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type OrderService struct {
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
	docService  *DocumentService
}

func NewOrderService(repo repository.OrderRepository, productRepo repository.ProductRepository, docService *DocumentService) *OrderService {
	return &OrderService{repo: repo, productRepo: productRepo, docService: docService}
}

func (s *OrderService) RegisterRoutes(r chi.Router) {
	r.Get("/orders", s.ListOrders)
	r.Post("/orders", s.CreateOrder)
	r.Get("/orders/{id}", s.GetOrderByID)
	r.Put("/orders/{id}", s.UpdateOrder)
	r.Get("/orders/{id}/invoice", s.GetInvoiceURL)
	r.Get("/orders/{id}/quote", s.GetQuoteURL)
}

func (s *OrderService) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.repo.GetAll(r.Context())
	if err != nil {
		response.InternalServerError(w, "Internal Server Error", err.Error())
		return
	}

	response.SuccessData(w, "Orders List", "Orders successfully obtained", orders)
}

func (s *OrderService) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}

	order, err := s.repo.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Order not found")
		return
	}

	response.SuccessData(w, "Order Found", "Order successfully obtained", order)
}

func (s *OrderService) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req core.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON")
		return
	}

	finalOrder := core.Order{
		CustomerID:    req.CustomerID,
		Status:        req.Status,
		Notes:         req.Notes,
		CustomerQuote: req.CustomerQuote,
		InvoiceNumber: req.InvoiceNumber,
		Items:         make([]core.OrderItem, 0),
	}

	var calculatedTotal float64

	//Business logic
	for _, itemReq := range req.Items {
		// Search price in DB
		if !gen.IsValidUUID(itemReq.ProductID) {
			response.BadRequest(w, "Bad Request", "Invalid product uuid")
			return
		}
		product, err := s.productRepo.GetByID(r.Context(), itemReq.ProductID)
		if err != nil {
			response.BadRequest(w, "Bad Request", fmt.Sprintf("Product %s not found or has no price", itemReq.ProductID))
			return
		}
		finalPrice := product.Price

		// Override price logic (DTO vs DB)
		if itemReq.UnitPrice != nil {
			finalPrice = *itemReq.UnitPrice
		}

		calculatedTotal += finalPrice * float64(itemReq.Quantity)

		// Add to entity items array
		finalOrder.Items = append(finalOrder.Items, core.OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: finalPrice,
		})
	}

	finalOrder.TotalAmount = calculatedTotal

	createdOrder, err := s.repo.Create(r.Context(), finalOrder)
	if err != nil {
		if err.Error() == "insufficient stock for product ID" || fmt.Sprintf("%v", err) != "" {
			// A simpler way is just pass the error msg to a 409
			response.Conflict(w, "Conflict", err.Error())
			return
		}
		response.Conflict(w, "Conflict", "Possible order duplicate or DB error")
		return
	}

	response.SuccessData(w, "Order Created", "Saved in inventory", createdOrder)
}

func (s *OrderService) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID")
		return
	}

	var req core.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Bad Request", "Invalid JSON")
		return
	}

	orderToUpdate := core.Order{
		ID:            id,
		Status:        req.Status,
		Notes:         req.Notes,
		InvoiceNumber: req.InvoiceNumber,
		CustomerID:    req.CustomerID,
		CustomerQuote: req.CustomerQuote,
	}

	updatedOrder, err := s.repo.Update(r.Context(), orderToUpdate)
	if err != nil {
		response.Conflict(w, "Conflict", err.Error())
		return
	}

	response.SuccessData(w, "Order updated", "Fields modified successfully", updatedOrder)
}

func (s *OrderService) GetInvoiceURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Validación de formato por consistencia
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}

	// 1. Obtener la orden (ya trae items y nombres de productos por tu repo)
	order, err := s.repo.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Order not found")
		return
	}

	var buf bytes.Buffer

	if err := s.docService.GenerateInvoicePDF(order, &buf); err != nil {
		response.InternalServerError(w, "Error", "Could not generate PDF")
		return
	}
	var invoice string
	if order.InvoiceNumber == nil || *order.InvoiceNumber == "" {
		invoice = "Boleta"
	} else {
		invoice = *order.InvoiceNumber
	}
	// Force to open in the browser (inline)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=Boleta_%s.pdf", invoice))
	// Define the exact content type
	w.Header().Set("Content-Type", "application/pdf")
	// Prevent the browser from trying to "guess" the file type
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.Write(buf.Bytes())
}

func (s *OrderService) GetQuoteURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Validación de formato por consistencia
	if !gen.IsValidUUID(id) {
		response.BadRequest(w, "Bad Request", "Invalid ID format")
		return
	}

	order, err := s.repo.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Not Found", "Order not found")
		return
	}

	var buf bytes.Buffer

	if err := s.docService.GenerateQuotePDF(order, &buf); err != nil {
		response.InternalServerError(w, "Error", "Could not generate PDF")
		return
	}

	// Force to open in the browser (inline)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=Cotizacion_%s.pdf", order.ID))
	// Define the exact content type
	w.Header().Set("Content-Type", "application/pdf")
	// Prevent the browser from trying to "guess" the file type
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.Write(buf.Bytes())
}
