package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TotalRetailPriceRequest(t *testing.T) {
	tests := []struct {
		input    totalRetailPriceRequest
		expected totalRetailPriceRequest
	}{
		{
			input:    totalRetailPriceRequest{Code: "test", Qty: 0},
			expected: totalRetailPriceRequest{Code: "test", Qty: 0},
		},
		{
			input:    totalRetailPriceRequest{Code: "", Qty: 12},
			expected: totalRetailPriceRequest{Code: "", Qty: 12},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual totalRetailPriceRequest
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Code == actual.Code, "~2|Test #%d expected code: %s, not code %s~", id, test.expected.Code, actual.Code)
		assert.True(t, test.expected.Qty == actual.Qty, "~2|Test #%d expected qty: %d, not qty %d~", id, test.expected.Qty, actual.Qty)
	}
}

func Test_TotalRetailPriceResponse(t *testing.T) {
	tests := []struct {
		input    totalRetailPriceResponse
		expected totalRetailPriceResponse
	}{
		{
			input:    totalRetailPriceResponse{Total: 100.99},
			expected: totalRetailPriceResponse{Total: 100.99},
		},
		{
			input:    totalRetailPriceResponse{Total: 0.0, Err: "test"},
			expected: totalRetailPriceResponse{Total: 0.0, Err: "test"},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual totalRetailPriceResponse
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
		{
			request:  "test",
			response: totalRetailPriceResponse{Err: "Invalid Request"},
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

		var actualResponse totalRetailPriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		testResponse := test.response.(totalRetailPriceResponse)

		assert.True(t, testResponse.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, testResponse.Err, actualResponse.Err)
		assert.True(t, testResponse.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testResponse.Total, actualResponse.Total)
	}
}
