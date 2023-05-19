package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bancodobrasil/featws-api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// func TestHomeHandler(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	r := gin.Default()
// 	r.GET("/", controllers.HomeHandler)

// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatalf("Couldn't create request: %v\n", err)
// 	}

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)
// 	defer w.Result().Body.Close()

// 	mockUserResp := `{"message":"FeatWS API Works!!!"}`
// 	//Using testify

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	responseData, _ := ioutil.ReadAll(w.Body)
// 	assert.Equal(t, mockUserResp, string(responseData))

// }
// TestUnitHomeHandler is a unit test for the HomeHandler function in a Go web application.
func TestUnitHomeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controllers.HomeHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
