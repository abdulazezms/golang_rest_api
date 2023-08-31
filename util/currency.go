package util

//supported currencies.
const (
	USD = "USD"
	EUR = "EUR"
	SAR = "SAR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, SAR:
		{
			return true
		}
	}
	return false
}
