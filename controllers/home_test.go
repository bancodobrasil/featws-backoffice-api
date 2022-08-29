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

func TestHomeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/", controllers.HomeHandler)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	defer w.Result().Body.Close()

	mockUserResp := `{"message":"FeatWS API Works!!!"}`
	//Using testify

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockUserResp, string(responseData))

}

r.GET("/recipes", ListRecipesHandler)

   req, _ := http.NewRequest("GET", "/recipes", nil)

   w := httptest.NewRecorder()

   r.ServeHTTP(w, req)
