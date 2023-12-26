package v1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/bancodobrasil/featws-api/dtos"
	mock_services "github.com/bancodobrasil/featws-api/mocks/services"
	payloads "github.com/bancodobrasil/featws-api/payloads/v1"
	responses "github.com/bancodobrasil/featws-api/responses/v1"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRulesheet_GetRulesheetByID is a test function that tests the functionality of
// retrieving a rulesheet by ID using Gin framework and mock services.
func TestRulesheet_GetRulesheetByID(t *testing.T) {
	// It's testing the behavior of a function that retrieves a rulesheet by ID. The test case is checking if
	// the function returns a 404 status code when the rulesheet with the given ID is not found. It creates a new
	// HTTP request with a mock context and a mock service that returns nil for the rulesheet with ID "1".
	// The test case then calls the function and checks if the HTTP response code is equal to 404.
	t.Run("Not Found Rulesheet Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
		srv := new(mock_services.Rulesheets)
		srv.On("Get", mock.Anything, "1").Return(nil, nil)
		v1.NewRulesheets(srv).GetRulesheet()(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// It is testing the behavior of a function that retrieves a rulesheet without an ID.
	t.Run("Test Without ID Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		srv := new(mock_services.Rulesheets)
		v1.NewRulesheets(srv).GetRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// It tests the "GetRulesheet" function of the "Rulesheets" API endpoint. The test simulates a request
	// to the endpoint with a parameter "id" set to 1. The mock service is used to simulate the behavior
	// of the actual service. In this case, the mock service returns an error when the "Get" function is called
	// with the parameter "1". The test then checks if the HTTP response code is equal to 500 using the "assert" package.
	t.Run("Wrong ID Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
		srv := new(mock_services.Rulesheets)
		srv.On("Get", mock.Anything, "1").Return(nil, errors.New("error"))
		v1.NewRulesheets(srv).GetRulesheet()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// It is testing the "GetRulesheet" function of the "Rulesheets" struct in the "v1" package.
	t.Run("Entity Different Nil Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		reponseEntity := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
		srv := new(mock_services.Rulesheets)
		srv.On("Get", mock.Anything, "1").Return(reponseEntity, nil)
		v1.NewRulesheets(srv).GetRulesheet()(c)
		assert.Equal(t, http.StatusOK, w.Code)

	})

}

// TestRulesheet_GetAllRulesheets tests the GetAllRulesheets function in the Rulesheets API endpoint, covering various
// scenarios such as error handling and normal flow with and without count query.
func TestRulesheet_GetAllRulesheets(t *testing.T) {
	// It tests the error handling of a limit query parameter that cannot be parsed. It creates a new HTTP
	// request with a malformed limit query parameter and sends it to a mock Rulesheets service using the Gin.
	// Thetest then checks that the HTTP response code is equal to 500 using the assert package.
	t.Run("Error on parse limit Query Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=?^&page=1")

		srv := new(mock_services.Rulesheets)
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	// It tests the error handling of a function that retrieves rulesheets from a server. The test creates
	// a mock server and sends a request with an invalid page query parameter. The test expects the server to return a 500.
	t.Run("Error on parse page Query Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=1&page=?^")

		srv := new(mock_services.Rulesheets)
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	// It tests the normal flow of a function that retrieves rulesheets with pagination. It creates a test
	// context and sets the request URL with query parameters for pagination. It then creates a mock service
	// for rulesheets and sets the expected behavior for the Find method of the service. Finally, it calls the
	// GetRulesheets function of the Rulesheets API with the test context and asserts that the response status code is HTTP 200 OK.
	t.Run("Normal flow with count query", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=1&page=1")

		srv := new(mock_services.Rulesheets)
		findOpts := &services.FindOptions{
			Limit: 1,
			Page:  1,
		}
		filter := make(map[string]interface{})
		srv.On("Find", mock.Anything, filter, findOpts).Return(nil, nil)
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// This test is testing the normal flow of getting rulesheets without count query and array of rulesheets.
	// It creates a test context and sets the request URL with limit and page parameters. It then creates a mock
	// service for rulesheets and sets the expected response entities. The test then calls the GetRulesheets function
	// and asserts that the HTTP status code returned is 200 OK.
	t.Run("Normal flow without count query and array of rulesheets", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=1&page=1")

		srv := new(mock_services.Rulesheets)
		findOpts := &services.FindOptions{
			Limit: 1,
			Page:  1,
		}
		filter := make(map[string]interface{})
		reponseEntities := []*dtos.Rulesheet{
			{
				ID:   uint(1),
				Name: "Test",
			},
			{
				ID:   uint(2),
				Name: "Test2",
			},
		}

		srv.On("Find", mock.Anything, filter, findOpts).Return(reponseEntities, nil)
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// It tests the error flow of a function that retrieves rulesheets from a service. The test creates a
	// mock service and sets expectations for the service's Find method to return an error. It then sends
	// a GET request to the function with query parameters for limit and page. Finally, it asserts that the
	// response code is equal to http.StatusInternalServerError.
	t.Run("Error flow without count query", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=1&page=1")

		srv := new(mock_services.Rulesheets)
		findOpts := &services.FindOptions{
			Limit: 1,
			Page:  1,
		}
		filter := make(map[string]interface{})
		srv.On("Find", mock.Anything, filter, findOpts).Return(nil, errors.New("error"))
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// It tests the normal flow of a function that retrieves rulesheets with a count query. It creates a
	// test context and sets the request URL with query parameters for limit, page, and count. It then
	// creates a mock service for rulesheets and sets expectations for the Find and Count methods. Finally,
	// it calls the GetRulesheets function and asserts that the HTTP status code returned is 200.
	t.Run("Normal flow with count query", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=1&page=1&count=true")

		srv := new(mock_services.Rulesheets)
		findOpts := &services.FindOptions{}
		filter := make(map[string]interface{})
		srv.On("Find", mock.Anything, filter, findOpts).Return(nil, nil)
		srv.On("Count", mock.Anything, filter).Return(int64(0), nil)
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// It tests the error flow of a function that retrieves rulesheets with a count query parameter. The test
	// creates a mock service for rulesheets and sets expectations for the Find and Count methods. It then sends
	// a GET request to the function with a URL containing the count query parameter. The test asserts that the response
	// code is HTTP 500.
	t.Run("Error flow with count query", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Request.URL, _ = url.Parse("?limit=1&page=1&count=true")

		srv := new(mock_services.Rulesheets)
		findOpts := &services.FindOptions{}
		filter := make(map[string]interface{})
		srv.On("Find", mock.Anything, filter, findOpts).Return(nil, nil)
		srv.On("Count", mock.Anything, filter).Return(int64(0), errors.New("error"))
		v1.NewRulesheets(srv).GetRulesheets()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

// TestRulesheet_CreateRulesheet is a test function for creating a rulesheet in a Go application, with different scenarios for
// normal flow and error handling.
func TestRulesheet_CreateRulesheet(t *testing.T) {
	// It tests the normal flow of creating a rulesheet using a mock service. It creates a new HTTP request with a JSON
	// payload, sets up a mock service to handle the request, and asserts that the response code is HTTP status code 201.
	t.Run("Normal flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		createdRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// It tests the behavior of a function that creates a new rulesheet by mocking the services layer. Specifically, it tests that
	// the function returns a 400 Bad Request HTTP status code when the request body is invalid. The test
	// creates a new HTTP request context, sets up a mock service, and calls the function being tested.
	// Finally, it uses the assert package to check that the HTTP status code returned by the function
	// matches the expected value.
	t.Run("Error on validate request body", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		srv := new(mock_services.Rulesheets)

		createdRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// It's testing the behavior of a function that creates a new rulesheet. The test is checking if the function
	// returns a status code of 400 (Bad Request) when the required fields in the payload are not provided. It creates a
	// mock service for the rulesheet and sets an expectation that the service's Create method will be
	// called with a specific argument. The test then calls the function and checks if the response status
	// code matches the expected value.
	t.Run("Error on validate required fields", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		createdRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error on validate required Name field", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		createdRulesheet := &dtos.Rulesheet{}

		srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := &responses.Error{}
		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, "Name", response.ValidationErrors[0].Field)
		assert.Equal(t, "required", response.ValidationErrors[0].Tag)
	})

	t.Run("Error on validate Name field does not start with digit", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			Name: "123Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		createdRulesheet := &dtos.Rulesheet{}

		srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := &responses.Error{}
		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, "Name", response.ValidationErrors[0].Field)
		assert.Equal(t, "doesNotStartWithDigit", response.ValidationErrors[0].Tag)
	})

	t.Run("Error on validate Name field must contain letter", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			Name: "#$%&*()",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		createdRulesheet := &dtos.Rulesheet{}

		srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := &responses.Error{}
		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, "Name", response.ValidationErrors[0].Field)
		assert.Equal(t, "mustContainLetter", response.ValidationErrors[0].Tag)
	})

	// t.Run("Error on define rulesheet entity", func(t *testing.T) {
	// 	gin.SetMode(gin.TestMode)

	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	c.Request = &http.Request{
	// 		Header: make(http.Header),
	// 	}

	// 	rules := make(map[string]interface{})
	// 	rules["test"] = 12
	// 	payload := &payloads.Rulesheet{
	// 		ID:    uint(1),
	// 		Name:  "test",
	// 		Rules: &rules,
	// 	}

	// 	bytedPayload, _ := json.Marshal(payload)

	// 	c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

	// 	srv := new(mock_services.Rulesheets)

	// 	createdRulesheet := &dtos.Rulesheet{
	// 		ID:   uint(1),
	// 		Name: "test",
	// 	}

	// 	srv.On("Create", mock.Anything, createdRulesheet).Return(nil)
	// 	v1.NewRulesheets(srv).CreateRulesheet()(c)
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// })

	// t.Run("Error on define rulesheet entity", func(t *testing.T) {
	// 	gin.SetMode(gin.TestMode)

	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	c.Request = &http.Request{
	// 		Header: make(http.Header),
	// 	}

	// 	payload := &payloads.Rulesheet{
	// 		ID:   uint(1),
	// 		Name: "Test",
	// 	}

	// 	bytedPayload, _ := json.Marshal(payload)

	// 	c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

	// 	srv := new(mock_services.Rulesheets)

	// 	createdRulesheet := &dtos.Rulesheet{
	// 		ID:   uint(1),
	// 		Name: "Test",
	// 	}

	// 	srv.On("Create", mock.Anything, createdRulesheet).Return(errors.New("error"))
	// 	v1.NewRulesheets(srv).CreateRulesheet()(c)
	// 	assert.Equal(t, http.StatusOK, w.Code)
	// })
}

// TestRulesheet_UpdateRulesheet is a test function that tests the update functionality of a rulesheet API endpoint
// using the Gin framework.
func TestRulesheet_UpdateRulesheet(t *testing.T) {
	// It tests the normal flow of the UpdateRulesheet function in the v1 package of a rulesheet API.
	t.Run("Normal flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// It tests the error flow of the UpdateRulesheet function in the v1 package of a rulesheet API. It creates a
	// mock service for the rulesheet API and sets expectations for the Get and Update methods of the service. It then
	// creates a test context and sets up a request with a payload. Finally, it calls the UpdateRulesheet
	// function and asserts that the response code is http.StatusBadRequest.
	t.Run("Error on check if the rulesheet exists flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// It's testing the error flow of the UpdateRulesheet function in the v1 package of a rulesheet API.
	// The test creates a mock service for the rulesheet API and sets expectations for the Get and Update
	// methods of the service. It then creates a test context and sets up the necessary request parameters
	// and payload. Finally, it calls the UpdateRulesheet function and asserts that the response code is http.StatusNotFound.
	t.Run("Error on check if the rulesheet exists flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, errors.New("error"))

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// It's testing the error scenario when binding the payload flow in the UpdateRulesheet function of the v1
	// package's Rulesheets struct.
	t.Run("Error on binding the payload flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// It is testing the "UpdateRulesheet" function of the "Rulesheets" API endpoint. The test is checking if
	// the API returns a status code of 400 (Bad Request) when a required field is missing in the request payload.
	// The test creates a mock service for the "Rulesheets" endpoint and sets up the necessary request context and
	// payload. It then calls the "UpdateRulesheet" function and checks if the response status code is 400.
	t.Run("Error on validate required fields flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// This code tests the error flow of updating a rulesheet with a slug that is already defined. It creates
	// a mock service for rulesheets and sets expectations for the Get and Update methods. It then creates a
	// test context and sets up aequest with a payload containing a rulesheet with a slug that is already defined.
	// It calls the UpdateRulesheet function and asserts that the response code is http.StatusBadRequest.
	t.Run("Error on slug already defined flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
			Slug: "test-slug",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// It tests the error flow of updating a rulesheet entity in a RESTful API. It creates a mock HTTP request context with a payload
	// containing the ID and name of the rulesheet to be updated. It then creates a mock service for
	// rulesheets and sets expectations for the Get and Update methods to return a rulesheet and an error,
	// respectively. Finally, it calls the UpdateRulesheet function of the v1 package with the mock
	// context and asserts that the HTTP response code is equal to 500 (Internal Server Error).
	t.Run("Error on update entity flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		newRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, errors.New("error"))
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// It tests the error status not found flow of the UpdateRulesheet function in the v1 package of a Rulesheets service.
	t.Run("Error status not found flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		payload := &payloads.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		bytedPayload, _ := json.Marshal(payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		oldRulesheet := &dtos.Rulesheet{
			ID:   uint(1),
			Name: "Test",
		}

		srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

		srv.On("Update", mock.Anything, *oldRulesheet).Return(nil, nil)
		v1.NewRulesheets(srv).UpdateRulesheet()(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// t.Run("Error on define entity Flow", func(t *testing.T) {
	// 	gin.SetMode(gin.TestMode)

	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	c.Request = &http.Request{
	// 		Header: make(http.Header),
	// 	}

	// 	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// 	rules := make(map[string]interface{})
	// 	rules["rules"] = 0

	// 	payload := &payloads.Rulesheet{
	// 		ID:    uint(1),
	// 		Name:  "Test",
	// 		Rules: &rules,
	// 	}

	// 	bytedPayload, _ := json.Marshal(payload)

	// 	c.Request.Body = ioutil.NopCloser(bytes.NewReader(bytedPayload))

	// 	srv := new(mock_services.Rulesheets)

	// 	oldRulesheet := &dtos.Rulesheet{
	// 		ID:   uint(1),
	// 		Name: "Test",
	// 	}

	// 	newRulesheet := &dtos.Rulesheet{
	// 		ID:   uint(1),
	// 		Name: "Test",
	// 	}

	// 	srv.On("Get", mock.Anything, "1").Return(oldRulesheet, nil)

	// 	srv.On("Update", mock.Anything, *oldRulesheet).Return(newRulesheet, nil)
	// 	v1.NewRulesheets(srv).UpdateRulesheet()(c)
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// })
}

// estRulesheet_DeleteRulesheet is a test function that tests the DeleteRulesheet function in the Rulesheets API endpoint,
// covering different scenarios and expected HTTP status codes.
func TestRulesheet_DeleteRulesheet(t *testing.T) {
	// It tests is the normal flow of the DeleteRulesheet function in the v1 package of a Rulesheets service.
	t.Run("Normal flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		srv.On("Delete", mock.Anything, "1").Return(true, nil)
		v1.NewRulesheets(srv).DeleteRulesheet()(c)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	// It's testing the DeleteRulesheet() function of the v1 package's Rulesheets struct. The test is checking if the
	// function returns a 404 status code when the Rulesheets service's Delete() function returns false
	// and no error. The test creates a mock Rulesheets service and sets the expected behavior of the
	// Delete() function. It then creates a new HTTP request context with a parameter "id" set to "1" and
	// calls the DeleteRulesheet() function with this context. Finally, it checks if the HTTP response
	t.Run("Error on delete rulesheet flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		srv.On("Delete", mock.Anything, "1").Return(false, nil)
		v1.NewRulesheets(srv).DeleteRulesheet()(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// The test checks if the function returns a 400 Bad Request HTTP status code when an ID is not passed as a parameter.
	// It uses a mock service object to simulate the behavior of the actual service.
	t.Run("Error on id not passed flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		srv := new(mock_services.Rulesheets)

		srv.On("Delete", mock.Anything, "1").Return(false, nil)
		v1.NewRulesheets(srv).DeleteRulesheet()(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// This test creates a new HTTP request context with a parameter `id` set to `1`. It then creates a mock `Rulesheets`
	// service and sets an expectation that the `Delete` method of the service will be called with the argument `"1"`. The
	// method is expected to return `false` and an error. The test then calls the `DeleteRulesheet` method
	// with the created context and asserts that
	t.Run("Error on delete rulesheet", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		srv := new(mock_services.Rulesheets)

		srv.On("Delete", mock.Anything, "1").Return(false, errors.New("error"))
		v1.NewRulesheets(srv).DeleteRulesheet()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
