package services_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/dtos"
	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/stretchr/testify/assert"
)

func SetupRulesheet() *dtos.Rulesheet {
	rulesheet := dtos.Rulesheet{
		ID:          1,
		Name:        "Test",
		Description: "Test",
	}
	return &rulesheet
}

func TestSaveAndCreateProject(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste"}`))
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects" {

			data, _ := io.ReadAll(r.Body)
			w.Write(data)
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects/0/repository/commits" {
			data, _ := io.ReadAll(r.Body)
			c := make(map[string]interface{})
			json.Unmarshal(data, &c)
			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := config.Config{
		GitlabURL:       s.URL,
		GitlabNamespace: namespace,
		GitlabToken:     "test",
		GitlabPrefix:    "prefix-",
		GitlabCIScript:  "test ci-script",
	}

	ngl := services.NewGitlab(&cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}

}

// func TestSaveCheckIfActionsAreAddedToCommit(t *testing.T) {
// 	dto := SetupRulesheet()

// }

// func TestSaveSucess(t *testing.T) {
// 	dto := &dtos.Rulesheet{
// 		ID:          1,
// 		Name:        "Teste",
// 		Description: "Teste",
// 	}

// 	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
// 	}))

// 	cfg := config.Config{
// 		GitlabURL: s.URL,
// 	}

// 	gitlabService := new(mocks_services.Gitlab)
// 	gitlabService.On("Save", &dto, "test").Return(nil)

// 	ngl := services.NewGitlab(&cfg)
// 	err := ngl.Save(dto, "test")
// 	if !assert.Nil(t, err) {
// 		t.Fail()
// 	}
// }
// func TestSaveGetNamespaceError(t *testing.T) {
// 	dto := &dtos.Rulesheet{
// 		ID:          1,
// 		Name:        "Teste",
// 		Description: "Teste",
// 	}

// 	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
// 	}))
// 	defer s.Close()

// 	cfg := config.Config{
// 		GitlabURL:       s.URL,
// 		GitlabNamespace: "",
// 		GitlabToken:     "glpat-kxhJPq3iYRr15dpaqAPr",
// 	}

// 	ngl := services.NewGitlab(&cfg)
// 	ngl.Connect()
// 	err := ngl.Save(dto, "test")

// 	if err == nil {
// 		t.Error("unexpected error")
// 	}
// }

// func TestSaveGetProject(t *testing.T) {
// 	dto := &dtos.Rulesheet{
// 		ID:          1,
// 		Name:        "Teste",
// 		Description: "Teste",
// 	}

// 	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

// 	}))
// 	defer s.Close()

// 	cfg := config.Config{
// 		GitlabURL:       s.URL,
// 		GitlabNamespace: "",
// 		GitlabToken:     "glpat-kxhJPq3iYRr15dpaqAPr",
// 		GitlabPrefix:    "",
// 	}

// 	ngl := services.NewGitlab(&cfg)
// 	ngl.Connect()
// 	err := ngl.Save(dto, "test")

// 	if err == nil {
// 		t.Error("unexpected error")
// 	}
// }

func TestFill(t *testing.T) {
	dto := &dtos.Rulesheet{
		ID:          1,
		Name:        "Teste",
		Description: "Teste",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	}))
	defer s.Close()

	cfg := config.Config{
		GitlabURL: s.URL,
	}

	gitlabService := new(mocks_services.Gitlab)
	gitlabService.On("Fill", &dto, "test").Return(nil)

	ngl := services.NewGitlab(&cfg)
	ngl.Connect()
	err := ngl.Fill(dto)
	if !assert.Nil(t, err) {
		t.Fail()
	}
}
