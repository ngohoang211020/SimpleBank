package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CAD, USD, EUR:
		return true
	}
	return false
}
