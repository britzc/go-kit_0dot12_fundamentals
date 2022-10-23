package transport

import "fmt"

type TotalRetailPriceRequest struct {
	Code string `json:"code"`
	Qty  int    `json:"qty"`
}

type TotalRetailPriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

type TotalWholesalePriceRequest struct {
	Partner string `json:"partner"`
	Code    string `json:"code"`
	Qty     int    `json:"qty"`
}

type TotalWholesalePriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

type ErrorResponse struct {
	Err string `json:"err, omitEmpty"`
}

func (e *ErrorResponse) Error() string {
	return e.Err
}

func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"err":"%s"}`, e.Err)), nil
}
