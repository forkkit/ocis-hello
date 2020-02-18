package proto_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

var invalidEndpoints = []string{"api", "api/v0", "greet", ""}

func TestPOSTCorrectEndpointVariousData(t *testing.T) {
	type TestStruct struct {
		testDataName         string
		requestBody          string
		expectedStatusCode   int
		responseBody         string
		expectedResponseBody string
	}

	var tests = []TestStruct{
		{"no-data", `{}`, 400, "", "missing a name\n"},
		{"value missing", `{"name":""}`, 400, "", "missing a name\n"},
		{"ASCII name", `{"name":"Milan"}`, 201, "", `{"message":"Hello Milan"}` + "\n"},
		{"UTF name", `{"name":"मिलन"}`, 201, "", `{"message":"Hello मिलन"}` + "\n"},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {

			rr := sendRequest(t, "POST", "/api/v0/greet", testCase.requestBody)
			assertResult(t, rr, testCase.expectedStatusCode, testCase.expectedResponseBody)
		})
	}
}

func TestPOSTIncorrectEndpoints(t *testing.T) {
	for _, endpoint := range invalidEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			rr := sendRequest(t, "POST", endpoint, `{"name":"Milan"}`)
			assertResult(t, rr, 404, "404 page not found\n")
		})
	}
}

func TestGETIncorrectEndpoints(t *testing.T) {
	for _, endpoint := range invalidEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			rr := sendRequest(t, "GET", endpoint, "")
			assertResult(t, rr, 404, "404 page not found\n")
		})
	}
}

func sendRequest(t *testing.T, method, endpoint, data string) *httptest.ResponseRecorder {
	var reader = strings.NewReader(data)
	req, err := http.NewRequest(method, endpoint, reader)
	assert.Nil(t, err)

	r := chi.NewRouter()
	proto.RegisterHelloWeb(r, svc.NewService())

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func assertResult(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int, expectedBody string) {
	assert.Equal(t, rr.Body.String(), expectedBody, "response body not as expected")
	assert.Equal(t, rr.Code, expectedStatusCode, "response code not as expected")
}
