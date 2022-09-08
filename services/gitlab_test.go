package services_test

import (
	"encoding/base64"
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
	"github.com/xanzy/go-gitlab"
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

func TestSaveAndUpdateProject(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Retrieving namespace
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste", "path":"test"}`))
			return
		}

		// Get project
		if r.URL.Path == "/api/v4/projects/test%2Fprefix-Test" {
			w.Write([]byte(`{"id":1,"description:null","name":"test"}`))
			return
		}

		// Create project
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects" {

			data, _ := io.ReadAll(r.Body)
			w.Write(data)
			return
		}

		//Create commits
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

func TestSaveTestFilesCreation(t *testing.T) {
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

			version := c["actions"].([]interface{})[0].(map[string]interface{})["content"].(string)
			assert.Equal(t, "1\n", version)
			gitlab_ci := c["actions"].([]interface{})[1].(map[string]interface{})["content"].(string)
			assert.Equal(t, "test ci-script", gitlab_ci)
			features := c["actions"].([]interface{})[2].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", features)
			parameters := c["actions"].([]interface{})[3].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", parameters)
			rulesFeatws := c["actions"].([]interface{})[4].(map[string]interface{})["content"].(string)
			assert.Equal(t, "", rulesFeatws)

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

func TestSaveTestFilesCreationWithFeatures(t *testing.T) {
	dto := SetupRulesheet()

	features := make([]map[string]interface{}, 0)

	features = append(features, map[string]interface{}{
		"name": "test1",
	})

	features = append(features, map[string]interface{}{
		"name": "test2",
	})

	features = append(features, map[string]interface{}{
		"name": "test3",
	})

	dto.Features = &features

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

			version := c["actions"].([]interface{})[0].(map[string]interface{})["content"].(string)
			assert.Equal(t, "1\n", version)
			gitlab_ci := c["actions"].([]interface{})[1].(map[string]interface{})["content"].(string)
			assert.Equal(t, "test ci-script", gitlab_ci)
			features := c["actions"].([]interface{})[2].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[\n  {\n    \"name\": \"test1\"\n  },\n  {\n    \"name\": \"test2\"\n  },\n  {\n    \"name\": \"test3\"\n  }\n]", features)
			parameters := c["actions"].([]interface{})[3].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", parameters)
			rulesFeatws := c["actions"].([]interface{})[4].(map[string]interface{})["content"].(string)
			assert.Equal(t, "", rulesFeatws)

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

func TestSaveTestFilesCreationWithParameters(t *testing.T) {
	dto := SetupRulesheet()

	parameters := make([]map[string]interface{}, 0)

	parameters = append(parameters, map[string]interface{}{
		"name": "test1",
	})

	parameters = append(parameters, map[string]interface{}{
		"name": "test2",
	})

	parameters = append(parameters, map[string]interface{}{
		"name": "test3",
	})

	dto.Parameters = &parameters

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

			version := c["actions"].([]interface{})[0].(map[string]interface{})["content"].(string)
			assert.Equal(t, "1\n", version)
			gitlab_ci := c["actions"].([]interface{})[1].(map[string]interface{})["content"].(string)
			assert.Equal(t, "test ci-script", gitlab_ci)
			features := c["actions"].([]interface{})[2].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", features)
			parameters := c["actions"].([]interface{})[3].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[\n  {\n    \"name\": \"test1\"\n  },\n  {\n    \"name\": \"test2\"\n  },\n  {\n    \"name\": \"test3\"\n  }\n]", parameters)
			rulesFeatws := c["actions"].([]interface{})[4].(map[string]interface{})["content"].(string)
			assert.Equal(t, "", rulesFeatws)

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

func TestSaveTestFilesCreationWithRules(t *testing.T) {
	dto := SetupRulesheet()

	parameters := make([]map[string]interface{}, 0)

	parameters = append(parameters, map[string]interface{}{
		"name": "test1",
	})

	parameters = append(parameters, map[string]interface{}{
		"name": "test2",
	})

	parameters = append(parameters, map[string]interface{}{
		"name": "test3",
	})

	dto.Parameters = &parameters

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

			version := c["actions"].([]interface{})[0].(map[string]interface{})["content"].(string)
			assert.Equal(t, "1\n", version)
			gitlab_ci := c["actions"].([]interface{})[1].(map[string]interface{})["content"].(string)
			assert.Equal(t, "test ci-script", gitlab_ci)
			features := c["actions"].([]interface{})[2].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", features)
			parameters := c["actions"].([]interface{})[3].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[\n  {\n    \"name\": \"test1\"\n  },\n  {\n    \"name\": \"test2\"\n  },\n  {\n    \"name\": \"test3\"\n  }\n]", parameters)
			rulesFeatws := c["actions"].([]interface{})[4].(map[string]interface{})["content"].(string)
			assert.Equal(t, "", rulesFeatws)

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

func TestSaveTestFilesUpdate(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste", "path":"testpath"}`))
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects" {

			data, _ := io.ReadAll(r.Body)
			w.Write(data)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/0/repository/files/VERSION" {
			content := base64.StdEncoding.EncodeToString([]byte("1\n"))

			file := gitlab.File{
				Content: content,
			}
			data, _ := json.Marshal(file)
			w.Write(data)
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects/0/repository/commits" {
			data, _ := io.ReadAll(r.Body)
			c := make(map[string]interface{})
			json.Unmarshal(data, &c)

			version := c["actions"].([]interface{})[0].(map[string]interface{})["content"].(string)
			assert.Equal(t, "2\n", version)
			gitlab_ci := c["actions"].([]interface{})[1].(map[string]interface{})["content"].(string)
			assert.Equal(t, "test ci-script", gitlab_ci)
			features := c["actions"].([]interface{})[2].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", features)
			parameters := c["actions"].([]interface{})[3].(map[string]interface{})["content"].(string)
			assert.Equal(t, "[]", parameters)
			rulesFeatws := c["actions"].([]interface{})[4].(map[string]interface{})["content"].(string)
			assert.Equal(t, "", rulesFeatws)

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

// Testing Errors
func TestSaveGitlabTokenNil(t *testing.T) {
	dto := SetupRulesheet()

	cfg := config.Config{
		GitlabURL:       "test",
		GitlabNamespace: "test",
	}

	ngl := services.NewGitlab(&cfg)
	err := ngl.Save(dto, "test")
	if err != nil {
		t.Error("expected nil return if gitlab token is nil")
	}

}

func TestSaveErrorOnFetchNameSpace(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}))
	defer s.Close()

	cfg := config.Config{
		GitlabURL:       s.URL,
		GitlabNamespace: namespace,
		GitlabToken:     "test",
	}

	ngl := services.NewGitlab(&cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")
	if err == nil {
		t.Error("expected error on fetch namespace")
	}

}

func TestSaveErrorOnFetchProject(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste"}`))
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}))
	defer s.Close()

	cfg := config.Config{
		GitlabURL:       s.URL,
		GitlabNamespace: namespace,
		GitlabToken:     "test",
	}

	ngl := services.NewGitlab(&cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")
	if err == nil {
		t.Error("expected error on fetch project")
	}

}

func TestSaveErrorOnCreateProject(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste"}`))
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v4/projects" {

			// _, _ := io.ReadAll(r.Body)
			w.Write(nil)
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

	if err == nil {
		t.Error("expected error on create project")
	}

}

func TestSaveErrorOnResolveVersion(t *testing.T) {
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

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/0/repository/files/VERSION" {
			w.WriteHeader(http.StatusBadRequest)
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

	if err == nil {
		t.Error("expected error on resolve version")
	}

}

func TestSaveErrorOnParseVersion(t *testing.T) {
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

// TODO: test save error on fetch version,
// func TestSaveErrorOnFetchVersion(t *testing.T) {
// 	dto := SetupRulesheet()

// 	namespace := "test"

// 	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == "/api/v4/namespaces/"+namespace {
// 			w.Write([]byte(`{"id":1,"name":"teste"}`))
// 			return
// 		}
// 		if r.Method == "POST" && r.URL.Path == "/api/v4/projects" {
// 			data, _ := io.ReadAll(r.Body)
// 			w.Write(data)
// 			return
// 		}

// 		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/0/repository/files/VERSION" {

// 		}

// 		if r.Method == "POST" && r.URL.Path == "/api/v4/projects/0/repository/commits" {
// 			data, _ := io.ReadAll(r.Body)
// 			c := make(map[string]interface{})
// 			json.Unmarshal(data, &c)
// 			w.Write(data)
// 			return
// 		}
// 		w.WriteHeader(http.StatusNotFound)
// 	}))
// 	defer s.Close()

// 	cfg := config.Config{
// 		GitlabURL:       s.URL,
// 		GitlabNamespace: namespace,
// 		GitlabToken:     "test",
// 		GitlabPrefix:    "prefix-",
// 		GitlabCIScript:  "test ci-script",
// 	}

// 	ngl := services.NewGitlab(&cfg)
// 	ngl.Connect()
// 	err := ngl.Save(dto, "test")

// 	if err != nil {
// 		t.Error("expected error on resolve version")
// 	}

// }
