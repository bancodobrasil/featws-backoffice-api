package v1_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/bancodobrasil/featws-api/dtos"
	mock_services "github.com/bancodobrasil/featws-api/mocks/services"
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

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
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
