package controllers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bancodobrasil/featws-api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthLive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/health/live", controllers.NewHealthController().HealthLiveHandler())

	req, err := http.NewRequest("GET", "/health/live", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	defer w.Result().Body.Close()

	mockUserResp := `{"goroutine-threshold": "OK"}`
	//Using testify

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockUserResp, string(responseData))

}
