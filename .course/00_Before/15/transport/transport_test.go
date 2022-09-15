package transport

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
		input    TotalRetailPriceRequest
		expected TotalRetailPriceRequest
	}{
		{
			input:    TotalRetailPriceRequest{Code: "test", Qty: 0},
			expected: TotalRetailPriceRequest{Code: "test", Qty: 0},
		},
		{
			input:    TotalRetailPriceRequest{Code: "", Qty: 12},
			expected: TotalRetailPriceRequest{Code: "", Qty: 12},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual TotalRetailPriceRequest
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Code == actual.Code, "~2|Test #%d expected code: %s, not code %s~", id, test.expected.Code, actual.Code)
		assert.True(t, test.expected.Qty == actual.Qty, "~2|Test #%d expected qty: %d, not qty %d~", id, test.expected.Qty, actual.Qty)
	}
}

func Test_TotalRetailPriceResponse(t *testing.T) {
	tests := []struct {
		input    TotalRetailPriceResponse
		expected TotalRetailPriceResponse
	}{
		{
			input:    TotalRetailPriceResponse{Total: 100.99},
			expected: TotalRetailPriceResponse{Total: 100.99},
		},
		{
			input:    TotalRetailPriceResponse{Total: 0.0, Err: "test"},
			expected: TotalRetailPriceResponse{Total: 0.0, Err: "test"},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual TotalRetailPriceResponse
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
		{
			request:  "test",
			response: TotalRetailPriceResponse{Err: "Invalid Request"},
		},
	}

	mockPricingService := new(MockPricingService)

	logger := &MockLogger{}
	totalRetailPriceHandler := MakeTotalRetailPriceHttpHandler(logger, mockPricingService)

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

		testResponse := test.response.(TotalRetailPriceResponse)

		assert.True(t, testResponse.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, testResponse.Err, actualResponse.Err)
		assert.True(t, testResponse.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testResponse.Total, actualResponse.Total)
	}
}

func Test_TotalWholesalePriceRequest(t *testing.T) {
	tests := []struct {
		input    TotalWholesalePriceRequest
		expected TotalWholesalePriceRequest
	}{
		{
			input:    TotalWholesalePriceRequest{Partner: "test", Code: "", Qty: 0},
			expected: TotalWholesalePriceRequest{Partner: "test", Code: "", Qty: 0},
		},
		{
			input:    TotalWholesalePriceRequest{Partner: "", Code: "test", Qty: 0},
			expected: TotalWholesalePriceRequest{Partner: "", Code: "test", Qty: 0},
		},
		{
			input:    TotalWholesalePriceRequest{Partner: "", Code: "", Qty: 12},
			expected: TotalWholesalePriceRequest{Partner: "", Code: "", Qty: 12},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual TotalWholesalePriceRequest
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Partner == actual.Partner, "~2|Test #%d expected partner: %s, not partner %s~", id, test.expected.Partner, actual.Partner)
		assert.True(t, test.expected.Code == actual.Code, "~2|Test #%d expected code: %s, not code %s~", id, test.expected.Code, actual.Code)
		assert.True(t, test.expected.Qty == actual.Qty, "~2|Test #%d expected qty: %d, not qty %d~", id, test.expected.Qty, actual.Qty)
	}
}

func Test_TotalWholesalePriceResponse(t *testing.T) {
	tests := []struct {
		input    TotalWholesalePriceResponse
		expected TotalWholesalePriceResponse
	}{
		{
			input:    TotalWholesalePriceResponse{Total: 100.99},
			expected: TotalWholesalePriceResponse{Total: 100.99},
		},
		{
			input:    TotalWholesalePriceResponse{Total: 0.0, Err: "test"},
			expected: TotalWholesalePriceResponse{Total: 0.0, Err: "test"},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual TotalWholesalePriceResponse
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
		{
			request:  TotalWholesalePriceRequest{Partner: "test", Code: "aaa111", Qty: 10},
			response: TotalWholesalePriceResponse{Err: "Partner Not Found"},
		},
		{
			request:  TotalWholesalePriceRequest{Partner: "superstore", Code: "fff000", Qty: 10},
			response: TotalWholesalePriceResponse{Err: "Code Not Found"},
		},
		{
			request:  "test",
			response: TotalWholesalePriceResponse{Err: "Invalid Request"},
		},
	}

	mockPricingService := new(MockPricingService)

	logger := &MockLogger{}
	totalWholesalePriceHandler := MakeTotalWholesalePriceHttpHandler(logger, mockPricingService)

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

		testResponse := test.response.(TotalWholesalePriceResponse)

		assert.True(t, testResponse.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, testResponse.Err, actualResponse.Err)
		assert.True(t, testResponse.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testResponse.Total, actualResponse.Total)
	}
}
