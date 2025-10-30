// Package currencycodes provides constants for various currency codes
package currencycodes

var (
	// Currency codes as per ISO 4217 standard
	BRL = "BRL" // Brazilian Real
	EUR = "EUR" // Euro
	GBP = "GBP" // British Pound Sterling
	JPY = "JPY" // Japanese Yen
	USD = "USD" // United States Dollar
	ARS = "ARS" // Argentine Peso
	PYG = "PYG" // Paraguayan Guarani
	CLP = "CLP" // Chilean Peso
	UYU = "UYU" // Uruguayan Peso
)

func IsValid(code string) bool {
	switch code {
	case BRL, EUR, GBP, JPY, USD, ARS, PYG, CLP, UYU:
		return true
	default:
		return false
	}
}
