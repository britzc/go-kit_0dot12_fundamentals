package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
)

var (
	ErrInvalidPartner  = errors.New("Invalid Partner Requested")
	ErrPartnerNotFound = errors.New("Partner Not Found")
	ErrInvalidCode     = errors.New("Invalid Code Requested")
	ErrCodeNotFound    = errors.New("Code Not Found")
	ErrInvalidQty      = errors.New("Invalid Quantity Requested")
)

type MockPricingService struct{}

func (MockPricingService) GetRetailTotal(code string, qty int) (total float64, err error) {
	if code == "" {
		return 0.0, ErrInvalidCode
	}
	if qty <= 0 {
		return 0.0, ErrInvalidQty
	}

	data := []string{
		"aaa111,12.99,10.99",
		"bbb222,2.90,2.50",
		"ccc333,22.50,21.00",
	}

	for _, line := range data {
		parts := strings.Split(line, ",")
		if parts[0] == code {
			price, _ := strconv.ParseFloat(parts[1], 64)

			return (price * float64(qty)), nil
		}
	}

	return 0.0, ErrCodeNotFound
}

func (MockPricingService) GetWholesaleTotal(partner string, code string, qty int) (total float64, err error) {
	if partner == "" {
		return 0.0, ErrInvalidPartner
	}
	if code == "" {
		return 0.0, ErrInvalidCode
	}
	if qty <= 0 {
		return 0.0, ErrInvalidQty
	}

	prices := []string{
		"aaa111,12.99",
		"bbb222,2.90",
		"ccc333,22.50",
	}

	price := 0.0
	priceFound := false
	for _, line := range prices {
		parts := strings.Split(line, ",")
		if parts[0] == code {
			priceFound = true
			price, _ = strconv.ParseFloat(parts[1], 64)
		}
	}

	if !priceFound {
		return 0.0, ErrCodeNotFound
	}

	partners := []string{
		"superstore,0.15",
		"joesdiscount,0.05",
	}

	discount := 0.0
	discountFound := false
	for _, line := range partners {
		parts := strings.Split(line, ",")
		if parts[0] == partner {
			discountFound = true
			discount, _ = strconv.ParseFloat(parts[1], 64)
		}
	}

	if !discountFound {
		return 0.0, ErrPartnerNotFound
	}

	saved := (price * discount)
	total = (price - saved) * float64(qty)

	return math.Round(total*100) / 100, nil
}

func Test_MakeTotalRetailPriceEndpoint(t *testing.T) {
	tests := []struct {
		request  TotalRetailPriceRequest
		response TotalRetailPriceResponse
	}{
		{
			request:  TotalRetailPriceRequest{Code: "", Qty: 0},
			response: TotalRetailPriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  TotalRetailPriceRequest{Code: "aaa111", Qty: 0},
			response: TotalRetailPriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  TotalRetailPriceRequest{Code: "aaa111", Qty: 15},
			response: TotalRetailPriceResponse{Total: 194.85},
		},
		{
			request:  TotalRetailPriceRequest{Code: "fff000", Qty: 10},
			response: TotalRetailPriceResponse{Err: "Code Not Found"},
		},
	}

	mockPricingService := new(MockPricingService)

	totalRetailPriceHandler := httptransport.NewServer(
		MakeTotalRetailPriceEndpoint(mockPricingService),
		decodeTotalRetailPriceRequest,
		encodeResponse,
	)

	server := httptest.NewServer(totalRetailPriceHandler)
	defer server.Close()

	for id, test := range tests {
		postBody, _ := json.Marshal(test.request)

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(server.URL, "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		var actualResponse TotalRetailPriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		assert.True(t, test.response.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, test.response.Err, actualResponse.Err)
		assert.True(t, test.response.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.response.Total, actualResponse.Total)
	}
}

func Test_MakeTotalWholesalePriceEndpoint(t *testing.T) {
	tests := []struct {
		request  TotalWholesalePriceRequest
		response TotalWholesalePriceResponse
	}{
		{
			request:  TotalWholesalePriceRequest{Partner: "", Code: "aaa111", Qty: 0},
			response: TotalWholesalePriceResponse{Err: "Invalid Partner Requested"},
		},
		{
			request:  TotalWholesalePriceRequest{Partner: "superstore", Code: "", Qty: 0},
			response: TotalWholesalePriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  TotalWholesalePriceRequest{Partner: "superstore", Code: "aaa111", Qty: 0},
			response: TotalWholesalePriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  TotalWholesalePriceRequest{Partner: "superstore", Code: "aaa111", Qty: 15},
			response: TotalWholesalePriceResponse{Total: 165.62},
		},
	}

	mockPricingService := new(MockPricingService)

	totalWholesalePriceHandler := httptransport.NewServer(
		MakeTotalWholesalePriceEndpoint(mockPricingService),
		decodeTotalWholesalePriceRequest,
		encodeResponse,
	)

	server := httptest.NewServer(totalWholesalePriceHandler)
	defer server.Close()

	for id, test := range tests {
		postBody, _ := json.Marshal(test.request)

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(server.URL, "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		var actualResponse TotalWholesalePriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		assert.True(t, test.response.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, test.response.Err, actualResponse.Err)
		assert.True(t, test.response.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.response.Total, actualResponse.Total)
	}
}
