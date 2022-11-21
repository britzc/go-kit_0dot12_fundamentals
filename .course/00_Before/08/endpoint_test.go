package main

import (
	"bytes"
	"context"
	"encoding/json"
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

func testDecodeTotalWholesalePriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request totalWholesalePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &errorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func Test_MakeTotalRetailPriceEndpoint(t *testing.T) {
	tests := []struct {
		request  totalRetailPriceRequest
		response totalRetailPriceResponse
	}{
		{
			request:  totalRetailPriceRequest{Code: "", Qty: 0},
			response: totalRetailPriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  totalRetailPriceRequest{Code: "aaa111", Qty: 0},
			response: totalRetailPriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  totalRetailPriceRequest{Code: "aaa111", Qty: 15},
			response: totalRetailPriceResponse{Total: 194.85},
		},
		{
			request:  totalRetailPriceRequest{Code: "fff000", Qty: 10},
			response: totalRetailPriceResponse{Err: "Code Not Found"},
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

		var actualResponse totalRetailPriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		assert.True(t, test.response.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, test.response.Err, actualResponse.Err)
		assert.True(t, test.response.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.response.Total, actualResponse.Total)
	}
}

func Test_MakeTotalWholesalePriceEndpoint(t *testing.T) {
	tests := []struct {
		request  totalWholesalePriceRequest
		response totalWholesalePriceResponse
	}{
		{
			request:  totalWholesalePriceRequest{Partner: "", Code: "aaa111", Qty: 0},
			response: totalWholesalePriceResponse{Err: "Invalid Partner Requested"},
		},
		{
			request:  totalWholesalePriceRequest{Partner: "superstore", Code: "", Qty: 0},
			response: totalWholesalePriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  totalWholesalePriceRequest{Partner: "superstore", Code: "aaa111", Qty: 0},
			response: totalWholesalePriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  totalWholesalePriceRequest{Partner: "superstore", Code: "aaa111", Qty: 15},
			response: totalWholesalePriceResponse{Total: 165.62},
		},
		/*
			{
				request:  totalWholesalePriceRequest{Partner: "blakes", Code: "aaa111", Qty: 10},
				response: totalWholesalePriceResponse{Err: "Partner Not Found"},
			},
			{
				request:  totalWholesalePriceRequest{Partner: "", Code: "fff000", Qty: 10},
				response: totalWholesalePriceResponse{Err: "Code Not Found"},
			},
		*/
	}

	mockPricingService := new(MockPricingService)

	totalWholesalePriceHandler := httptransport.NewServer(
		MakeTotalWholesalePriceEndpoint(mockPricingService),
		testDecodeTotalWholesalePriceRequest,
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

		var actualResponse totalWholesalePriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		assert.True(t, test.response.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, test.response.Err, actualResponse.Err)
		assert.True(t, test.response.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.response.Total, actualResponse.Total)
	}
}
