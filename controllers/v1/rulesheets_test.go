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
	"github.com/bancodobrasil/featws-api/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRulesheet_GetRulesheetByID(t *testing.T) {
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

func TestRulesheet_GetAllRulesheets(t *testing.T) {
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

func TestRulesheet_CreateRulesheet(t *testing.T) {
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

	t.Run("Error on define rulesheet entity", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		rules := make(map[string]interface{})
		rules["test"] = 12
		payload := &payloads.Rulesheet{
			ID:    uint(1),
			Name:  "test",
			Rules: &rules,
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
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Error on define rulesheet entity", func(t *testing.T) {
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

		srv.On("Create", mock.Anything, createdRulesheet).Return(errors.New("error"))
		v1.NewRulesheets(srv).CreateRulesheet()(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRulesheet_UpdateRulesheet(t *testing.T) {
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

	t.Run("Error on define entity Flow", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		rules := make(map[string]interface{})
		rules["rules"] = 0

		payload := &payloads.Rulesheet{
			ID:    uint(1),
			Name:  "Test",
			Rules: &rules,
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
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRulesheet_DeleteRulesheet(t *testing.T) {
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
