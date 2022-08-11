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

	"github.com/britzc/go-kit_0dot12_fundamentals/current/payload"
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

func (MockPricingService) GetWholesaleTotal(partner, code string, qty int) (total float64, err error) {
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

func Test_TotalRetailPriceRequest(t *testing.T) {
	tests := []struct {
		input    payload.TotalRetailPriceRequest
		expected payload.TotalRetailPriceRequest
	}{
		{
			input:    payload.TotalRetailPriceRequest{Code: "test", Qty: 0},
			expected: payload.TotalRetailPriceRequest{Code: "test", Qty: 0},
		},
		{
			input:    payload.TotalRetailPriceRequest{Code: "", Qty: 12},
			expected: payload.TotalRetailPriceRequest{Code: "", Qty: 12},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual payload.TotalRetailPriceRequest
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Code == actual.Code, "~2|Test #%d expected code: %s, not code %s~", id, test.expected.Code, actual.Code)
		assert.True(t, test.expected.Qty == actual.Qty, "~2|Test #%d expected qty: %d, not qty %d~", id, test.expected.Qty, actual.Qty)
	}
}

func Test_TotalRetailPriceResponse(t *testing.T) {
	tests := []struct {
		input    payload.TotalRetailPriceResponse
		expected payload.TotalRetailPriceResponse
	}{
		{
			input:    payload.TotalRetailPriceResponse{Total: 100.99},
			expected: payload.TotalRetailPriceResponse{Total: 100.99},
		},
		{
			input:    payload.TotalRetailPriceResponse{Total: 0.0, Err: "test"},
			expected: payload.TotalRetailPriceResponse{Total: 0.0, Err: "test"},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual payload.TotalRetailPriceResponse
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Total == actual.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.expected.Total, actual.Total)
		assert.True(t, test.expected.Err == actual.Err, "~2|Test #%d expected err: %s, not err %s~", id, test.expected.Err, actual.Err)
	}
}

func Test_MakeTotalRetailPriceHttpHandler(t *testing.T) {
	tests := []struct {
		request  interface{}
		response interface{}
	}{
		{
			request:  payload.TotalRetailPriceRequest{Code: "", Qty: 0},
			response: payload.TotalRetailPriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  payload.TotalRetailPriceRequest{Code: "aaa111", Qty: 0},
			response: payload.TotalRetailPriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  payload.TotalRetailPriceRequest{Code: "aaa111", Qty: 15},
			response: payload.TotalRetailPriceResponse{Total: 194.85},
		},
		{
			request:  payload.TotalRetailPriceRequest{Code: "fff000", Qty: 10},
			response: payload.TotalRetailPriceResponse{Err: "Code Not Found"},
		},
		{
			request:  "test",
			response: payload.TotalRetailPriceResponse{Err: "Invalid Request"},
		},
	}

	mockPricingService := new(MockPricingService)

	totalRetailPriceHandler := MakeTotalRetailPriceHttpHandler(mockPricingService)

	server := httptest.NewServer(totalRetailPriceHandler)
	defer server.Close()

	for id, test := range tests {
		postBody, _ := json.Marshal(test.request)

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(server.URL, "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		var actualResponse payload.TotalRetailPriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		testResponse := test.response.(payload.TotalRetailPriceResponse)

		assert.True(t, testResponse.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, testResponse.Err, actualResponse.Err)
		assert.True(t, testResponse.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testResponse.Total, actualResponse.Total)
	}
}

func Test_TotalWholesalePriceRequest(t *testing.T) {
	tests := []struct {
		input    payload.TotalWholesalePriceRequest
		expected payload.TotalWholesalePriceRequest
	}{
		{
			input:    payload.TotalWholesalePriceRequest{Partner: "test", Code: "", Qty: 0},
			expected: payload.TotalWholesalePriceRequest{Partner: "test", Code: "", Qty: 0},
		},
		{
			input:    payload.TotalWholesalePriceRequest{Partner: "", Code: "test", Qty: 0},
			expected: payload.TotalWholesalePriceRequest{Partner: "", Code: "test", Qty: 0},
		},
		{
			input:    payload.TotalWholesalePriceRequest{Partner: "", Code: "", Qty: 12},
			expected: payload.TotalWholesalePriceRequest{Partner: "", Code: "", Qty: 12},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual payload.TotalWholesalePriceRequest
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Partner == actual.Partner, "~2|Test #%d expected partner: %s, not partner %s~", id, test.expected.Partner, actual.Partner)
		assert.True(t, test.expected.Code == actual.Code, "~2|Test #%d expected code: %s, not code %s~", id, test.expected.Code, actual.Code)
		assert.True(t, test.expected.Qty == actual.Qty, "~2|Test #%d expected qty: %d, not qty %d~", id, test.expected.Qty, actual.Qty)
	}
}

func Test_TotalWholesalePriceResponse(t *testing.T) {
	tests := []struct {
		input    payload.TotalWholesalePriceResponse
		expected payload.TotalWholesalePriceResponse
	}{
		{
			input:    payload.TotalWholesalePriceResponse{Total: 100.99},
			expected: payload.TotalWholesalePriceResponse{Total: 100.99},
		},
		{
			input:    payload.TotalWholesalePriceResponse{Total: 0.0, Err: "test"},
			expected: payload.TotalWholesalePriceResponse{Total: 0.0, Err: "test"},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual payload.TotalWholesalePriceResponse
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Total == actual.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.expected.Total, actual.Total)
		assert.True(t, test.expected.Err == actual.Err, "~2|Test #%d expected err: %s, not err %s~", id, test.expected.Err, actual.Err)
	}
}

func Test_MakeTotalWholesalePriceHttpHandler(t *testing.T) {
	tests := []struct {
		request  interface{}
		response interface{}
	}{
		{
			request:  payload.TotalWholesalePriceRequest{Partner: "", Code: "aaa111", Qty: 0},
			response: payload.TotalWholesalePriceResponse{Err: "Invalid Partner Requested"},
		},
		{
			request:  payload.TotalWholesalePriceRequest{Partner: "superstore", Code: "", Qty: 0},
			response: payload.TotalWholesalePriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  payload.TotalWholesalePriceRequest{Partner: "superstore", Code: "aaa111", Qty: 0},
			response: payload.TotalWholesalePriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  payload.TotalWholesalePriceRequest{Partner: "superstore", Code: "aaa111", Qty: 15},
			response: payload.TotalWholesalePriceResponse{Total: 165.62},
		},
		{
			request:  payload.TotalWholesalePriceRequest{Partner: "test", Code: "aaa111", Qty: 10},
			response: payload.TotalWholesalePriceResponse{Err: "Partner Not Found"},
		},
		{
			request:  payload.TotalWholesalePriceRequest{Partner: "superstore", Code: "fff000", Qty: 10},
			response: payload.TotalWholesalePriceResponse{Err: "Code Not Found"},
		},
		{
			request:  "test",
			response: payload.TotalWholesalePriceResponse{Err: "Invalid Request"},
		},
	}

	mockPricingService := new(MockPricingService)

	totalWholesalePriceHandler := MakeTotalWholesalePriceHttpHandler(mockPricingService)

	server := httptest.NewServer(totalWholesalePriceHandler)
	defer server.Close()

	for id, test := range tests {
		postBody, _ := json.Marshal(test.request)

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(server.URL, "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		var actualResponse payload.TotalWholesalePriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		testResponse := test.response.(payload.TotalWholesalePriceResponse)

		assert.True(t, testResponse.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, testResponse.Err, actualResponse.Err)
		assert.True(t, testResponse.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testResponse.Total, actualResponse.Total)
	}
}
