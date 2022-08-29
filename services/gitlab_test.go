package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bancodobrasil/featws-api/dtos"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/services"
)

func TestSaveSucess(t *testing.T) {
	dto := &dtos.Rulesheet{
		ID:          1,
		Name:        "Teste",
		Description: "Teste",
	}
	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Save", &dto, "test").Return(nil)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := services.Gitlab.Save(nil, dto, "test")
		if err != nil {
			t.Error("unexpected error:", err)
		}

	}))
	defer s.Close()
}
