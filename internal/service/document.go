package service

import (
	"bytes"
	"embed"
	"fmt"
	"inventory-backend/internal/core"
	"io"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type DocumentService struct {
	assets      embed.FS
	logoPath    string
	companyName string
}

func NewDocumentService(assets embed.FS, logoPath, companyName string) *DocumentService {
	return &DocumentService{
		assets:      assets,
		logoPath:    logoPath,
		companyName: companyName,
	}
}

func (s *DocumentService) GenerateInvoicePDF(order core.OrderResponse, w io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	logoData, err := s.assets.ReadFile(s.logoPath)
	if err == nil {
		// Si el logo existe, lo ponemos en el PDF
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("logo", opts, bytes.NewReader(logoData))
		pdf.ImageOptions("logo", 10, 10, 30, 0, false, opts, 0, "")
	} else {
		fmt.Println("Logo not found")
		fmt.Println(err)
	}

	pdf.SetY(13)
	pdf.SetX(45)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(100, 10, s.companyName)
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(100, 10, "RUT: ...")
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.Cell(100, 10, "Direccion: ...")
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.Cell(100, 10, "Telefono: ...")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "BOLETA DE VENTA")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 10)
	if order.InvoiceNumber != nil && *order.InvoiceNumber != "" {
		pdf.Cell(40, 10, fmt.Sprintf("Numero de Boleta: %s", *order.InvoiceNumber))
	} else {
		pdf.Cell(40, 10, fmt.Sprintf("Numero de Boleta: %s", "PENDIENTE"))
	}
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Fecha: %s", order.CreatedAt.Format("02/01/2006")))
	pdf.Ln(15)

	// Items Table
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(100, 8, "Producto", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Cant.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Precio", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Subtotal", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, item := range order.Items {
		pdf.CellFormat(100, 8, item.ProductName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, formatCLP(item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, formatCLP(item.SubTotal), "1", 1, "R", false, 0, "")
	}

	// Total
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(160, 10, "Subtotal Neto:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, formatCLP(order.TotalAmount*0.81), "1", 1, "R", false, 0, "")
	pdf.CellFormat(160, 10, "IVA 19%:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, formatCLP(order.TotalAmount*0.19), "1", 1, "R", false, 0, "")
	pdf.CellFormat(160, 10, "TOTAL con IVA:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, formatCLP(order.TotalAmount), "1", 1, "R", false, 0, "")

	return pdf.Output(w)
}

func (s *DocumentService) GenerateQuotePDF(order core.OrderResponse, w io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	logoData, err := s.assets.ReadFile(s.logoPath)
	if err == nil {
		// Si el logo existe, lo ponemos en el PDF
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("logo", opts, bytes.NewReader(logoData))
		pdf.ImageOptions("logo", 10, 10, 30, 0, false, opts, 0, "")
	} else {
		fmt.Println("Logo not found")
		fmt.Println(err)
	}

	pdf.SetY(13)
	pdf.SetX(45)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(100, 10, s.companyName)
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(100, 10, "RUT: ...")
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.Cell(100, 10, "Direccion: ...")
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.Cell(100, 10, "Telefono: ...")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "COTIZACION")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Cotizacion valida hasta 7 dias desde la fecha de emision")
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Fecha de emision: %s", order.CreatedAt.Format("02/01/2006")))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Valido hasta: %s", order.CreatedAt.AddDate(0, 0, 7).Format("02/01/2006")))
	pdf.Ln(15)

	// Items Table
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(100, 8, "Producto", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Cant.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Precio", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Subtotal", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, item := range order.Items {
		pdf.CellFormat(100, 8, item.ProductName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, formatCLP(item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, formatCLP(item.SubTotal), "1", 1, "R", false, 0, "")
	}

	// Total
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(160, 10, "Subtotal Neto:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, formatCLP(order.TotalAmount*0.81), "1", 1, "R", false, 0, "")
	pdf.CellFormat(160, 10, "IVA 19%:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, formatCLP(order.TotalAmount*0.19), "1", 1, "R", false, 0, "")
	pdf.CellFormat(160, 10, "TOTAL con IVA:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, formatCLP(order.TotalAmount), "1", 1, "R", false, 0, "")

	return pdf.Output(w)
}

func formatCLP(amount float64) string {
	intAmount := int64(amount)
	s := fmt.Sprintf("%d", intAmount)

	//Insert points every 3 digits
	var result []string
	n := len(s)
	for i := n; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{s[start:i]}, result...)
	}

	return "$" + strings.Join(result, ".")
}
