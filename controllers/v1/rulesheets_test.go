package v1_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/bancodobrasil/featws-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Rulesheet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func TestListRulesheets(t *testing.T) {

	gin.SetMode(gin.TestMode)
	// service, _ := services.Rulesheets.Find()

	r := gin.Default()
	r.GET("/v1/rulesheets", v1.Rulesheets.GetRulesheets())

	req, err := http.NewRequest("GET", "/api/v1/rulesheets", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	defer w.Result().Body.Close()

	var rulesheets []Rulesheet
	json.Unmarshal([]byte(w.Body.String()), &rulesheets)
	//Using testify
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 10, len(rulesheets))

}

func TestFindRulesheet(t *testing.T) {

	// iniciolizacao de um servidor Gin
	ts := httptest.NewServer(routes.SetupRoutesWithReturn())
	defer ts.Close()

	// criacao da folha de regras esperada
	expectedRulesheet := Rulesheet{
		ID:          "1",
		Name:        "Rulesheet 1",
		Description: "Rulesheet 1 description",
	}

	// chamada do servico de busca da folha de regras especifica
	resp, err := http.Get(fmt.Sprintf("%s/v1/rulesheets/1", ts.URL))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}
	defer resp.Body.Close()
	//Usando o testify

	//Testa se os erros sao nulos
	assert.Nil(t, err)

	// Testa se o codigo de resposta e 200
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Faz a leitura do corpo da resposta
	data, _ := ioutil.ReadAll(resp.Body)

	// instancia um objeto do tipo Rulesheet, esta errado nao sei qual utlizar
	var actualRulesheet Rulesheet

	// faz o parse do corpo da resposta para o objeto actualRulesheet
	json.Unmarshal(data, &actualRulesheet)

	//Using testify
	assert.Equal(t, expectedRulesheet.Name, actualRulesheet.Name)
}
