package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecodeTotalRetailPriceRequest(t *testing.T) {
	tests := []struct {
		request interface{}
		err     error
	}{
		{
			request: totalRetailPriceRequest{Code: "aaa111", Qty: 15},
			err:     nil,
		},
		{
			request: "test",
			err:     ErrInvalidRequest,
		},
	}

	for id, test := range tests {
		testJson, _ := json.Marshal(test.request)
		testBody := bytes.NewReader(testJson)
		r, _ := http.NewRequest(http.MethodPost, "/resource", testBody)

		decodeResult, err := DecodeTotalRetailPriceRequest(nil, r)
		assert.True(t, test.err == err, "~2|Test #%d expected error: %s, not error %s~", id, test.err, err)
		if err == nil {
			actualReq := decodeResult.(totalRetailPriceRequest)

			testReq := test.request.(totalRetailPriceRequest)

			assert.True(t, testReq.Code == actualReq.Code, "~2|Test #%d expected code: %s, not code %s~", id, testReq.Code, actualReq.Code)
			assert.True(t, testReq.Qty == actualReq.Qty, "~2|Test #%d expected qty: %.2f, not qty %.2f~", id, testReq.Qty, actualReq.Qty)
		}
	}
}

func Test_Encode(t *testing.T) {
	tests := []struct {
		response interface{}
		err      error
	}{
		{
			response: totalRetailPriceResponse{Total: 100.99, Err: ""},
			err:      nil,
		},
	}

	for id, test := range tests {
		w := httptest.NewRecorder()

		err := EncodeResponse(nil, w, test.response)
		assert.True(t, test.err == err, "~2|Test #%d expected error: %s, not error %s~", id, test.err, err)
		if err == nil {
			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			var actualRes totalRetailPriceResponse
			json.Unmarshal(data, &actualRes)

			testRes := test.response.(totalRetailPriceResponse)

			assert.True(t, testRes.Err == actualRes.Err, "~2|Test #%d expected error: %s, not error %s~", id, testRes.Err, actualRes.Err)
			assert.True(t, testRes.Total == actualRes.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testRes.Total, actualRes.Total)
		}
	}
}
