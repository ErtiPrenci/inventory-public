package gen

import (
	"inventory-backend/internal/core"
	"regexp"
	"strings"
	"unicode"
)

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

func IsValidUUID(u string) bool {
	if u == "" {
		return false
	}
	return uuidRegex.MatchString(u)
}

func IsValidRUT(rut string) bool {
	if rut == "" {
		return false
	}

	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")
	rut = strings.ToUpper(strings.TrimSpace(rut))

	if len(rut) < 2 {
		return false
	}

	body := rut[:len(rut)-1]
	dv := rut[len(rut)-1]

	sum := 0
	multiplier := 2

	for i := len(body) - 1; i >= 0; i-- {
		char := body[i]

		if !unicode.IsDigit(rune(char)) {
			return false
		}

		digit := int(char - '0')

		sum += digit * multiplier
		multiplier++
		if multiplier > 7 {
			multiplier = 2
		}
	}

	remainder := sum % 11
	result := 11 - remainder

	var expectedDV uint8
	switch result {
	case 11:
		expectedDV = '0'
	case 10:
		expectedDV = 'K'
	default:
		expectedDV = uint8(result + '0')
	}

	return dv == expectedDV
}

func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}

func IsValidPhone(phone string) bool {
	if phone == "" {
		return false
	}

	return regexp.MustCompile(`^\+?\d{10,}$`).MatchString(phone)
}

func IsValidOrderTransition(previousState core.OrderStatus, newState core.OrderStatus) bool {
	if previousState == core.OrderStatusSold && newState == core.OrderStatusQuote {
		return false
	}
	if previousState == core.OrderStatusCancelled && newState == core.OrderStatusQuote {
		return false
	}
	return true
}
