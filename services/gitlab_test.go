package services_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/dtos"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func SetupConfig(url *httptest.Server) *config.Config {
	cfg := config.Config{
		GitlabURL:       url.URL,
		GitlabNamespace: "test",
		GitlabToken:     "test",
		GitlabPrefix:    "prefix-",
		GitlabCIScript:  "test ci-script",
	}
	return &cfg
}

func SetupRulesheet() *dtos.Rulesheet {
	rulesheet := dtos.Rulesheet{
		ID:          1,
		Name:        "Test",
		Description: "Test",
	}
	return &rulesheet
}

// Functions for test Save function
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestSaveTestFilesCreationWithRuleInterface(t *testing.T) {
	dto := SetupRulesheet()

	rules := []interface{}{
		&dtos.Rule{
			Condition: "test",
			Value: map[string]string{
				"nomeAplicativo": "testAplicativo",
				"textoUrlPadrao": "testURLpadrao",
				"textoUrlDesvio": "testURLdesvio",
			},
			Type: "testType",
		},
	}

	mappedRules := map[string]interface{}{
		"tags": rules,
	}

	dto.Rules = &mappedRules

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
			assert.Equal(t, "[[tags]]\ncondition = test\nvalue = {\"nomeAplicativo\":\"testAplicativo\",\"textoUrlDesvio\":\"testURLdesvio\",\"textoUrlPadrao\":\"testURLpadrao\"}\ntype = object\n\n", rulesFeatws)

			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestSaveTestFilesCreationWithRuleString(t *testing.T) {
	dto := SetupRulesheet()

	mappedRules := map[string]interface{}{
		"rule1": "true",
		"rule2": "\"test\"",
	}
	dto.Rules = &mappedRules

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
			assert.Equal(t, "rule1 = true\nrule2 = \"test\"\n", rulesFeatws)

			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestSaveTestFilesCreationWithStringRule(t *testing.T) {
	dto := SetupRulesheet()

	rules := make(map[string]interface{})

	rules["test1"] = "test1"
	rules["test2"] = "test2"
	rules["test3"] = "test3"

	dto.Rules = &rules

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
			assert.Equal(t, "test1 = test1\ntest2 = test2\ntest3 = test3\n", rulesFeatws)

			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestSaveTestFilesCreationWithDefaultRule(t *testing.T) {
	dto := SetupRulesheet()

	rules := []interface{}{
		&dtos.Rule{
			Condition: "test",
			Value: map[string]string{
				"NomeAplicativo": "testAplicativo",
				"TextoURLPadrao": "testURLpadrao",
				"TextoURLDesvio": "testURLdesvio",
			},
			Type: "testType",
		},
	}

	mappedRules := map[string]interface{}{
		// "mystring": "teste",
		"tags": rules,
		// "mycomplexrule": &dtos.Rule{
		// 	Condition: "true",
		// 	Value: map[string]string{
		// 		"field1": "value1",
		// 		"field2": "value2",
		// 	},
		// 	Type: "testType",
		// },
	}

	dto.Rules = &mappedRules

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
			assert.Equal(t, "[[tags]]\ncondition = test\nvalue = {\"NomeAplicativo\":\"testAplicativo\",\"TextoURLDesvio\":\"testURLdesvio\",\"TextoURLPadrao\":\"testURLpadrao\"}\ntype = object\n\n", rulesFeatws)

			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}
}

// Functions to test fill function
func TestFill(t *testing.T) {

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste", "path":"testpath"}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/testpath/prefix-Test" {
			w.Write([]byte(`{"id":1,"description":"testeDesc","name":"teste"}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/1/repository/files/rules.featws" {
			content := base64.StdEncoding.EncodeToString([]byte("regra = $test"))

			file := gitlab.File{
				Content: content,
			}
			data, _ := json.Marshal(file)
			w.Write(data)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/1/repository/files/parameters.json" {
			content := base64.StdEncoding.EncodeToString([]byte(`[
				{
					"name": "param1",
					"type": "string"
				},
				{
					"name": "param2",
					"type": "string"
				  }				
			]`))

			file := gitlab.File{
				Content: content,
			}
			data, _ := json.Marshal(file)
			w.Write(data)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/1/repository/files/features.json" {
			content := base64.StdEncoding.EncodeToString([]byte(`[
				{
					"name": "feat1",
					"type": "string"
				},
				{
					"name": "feat2",
					"type": "string"
				  }				
			]`))

			file := gitlab.File{
				Content: content,
			}
			data, _ := json.Marshal(file)
			w.Write(data)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	dto := SetupRulesheet()
	err := ngl.Fill(dto)
	if err != nil {
		t.Errorf("unexpected error on fill: %s", err.Error())
		return
	}

	if (*dto.Rules)["regra"].(string) != "$test" {
		t.Error("error on unmarshalling rules")
		return
	}

	if (dto.Parameters) == nil || len(*dto.Parameters) != 2 {
		t.Error("error on unmarshalling parameters")
		return
	}

	param1 := (*dto.Parameters)[0]
	if param1["name"] != "param1" || param1["type"] != "string" {
		t.Error("error on unmarshalling parameter 1")
		return
	}

	param2 := (*dto.Parameters)[1]
	if param2["name"] != "param2" || param1["type"] != "string" {
		t.Error("error on unmarshalling parameter 2")
		return
	}

	if (dto.Features) == nil || len(*dto.Features) != 2 {
		t.Error("error on unmarshalling Features")
		return
	}

	feat1 := (*dto.Features)[0]
	if feat1["name"] != "feat1" || feat1["type"] != "string" {
		t.Error("error on unmarshalling Feature 1")
		return
	}

	feat2 := (*dto.Features)[1]
	if feat2["name"] != "feat2" || param1["type"] != "string" {
		t.Error("error on unmarshalling Feature 2")
		return
	}
}

// Functions to test fill function
func TestFillRulesSlices(t *testing.T) {
	dto := SetupRulesheet()

	namespace := "test"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/namespaces/"+namespace {
			w.Write([]byte(`{"id":1,"name":"teste", "path":"testpath"}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/testpath/prefix-Test" {
			w.Write([]byte(`{"id":1,"description":"testeDesc","name":"teste"}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/1/repository/files/rules.featws" {
			content := base64.StdEncoding.EncodeToString([]byte("[feat]\n condition = $test\n[[tags]]\n condition = $test"))

			file := gitlab.File{
				Content: content,
			}
			data, _ := json.Marshal(file)
			w.Write(data)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/v4/projects/1/repository/files/parameters.json" {
			content := base64.StdEncoding.EncodeToString([]byte("[]"))

			file := gitlab.File{
				Content: content,
			}
			data, _ := json.Marshal(file)
			w.Write(data)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Fill(dto)

	feat, ok := (*dto.Rules)["feat"]
	if !ok {
		t.Error("error on unmarshalling rules: not found feat")
		return
	}

	featMap, ok := feat.(map[string]interface{})
	if !ok {
		t.Error("error on unmarshalling rules: feat is no map")
		return
	}

	if featMap["condition"].(string) != "$test" {
		t.Error("error on unmarshalling rules")
		return
	}

	tags, ok := (*dto.Rules)["tags"]
	if !ok {
		t.Error("error on unmarshalling rules: not found tags array")
		return
	}

	arr, ok := tags.([]map[string]interface{})
	if !ok {
		t.Error("error on unmarshalling rules: tags is no array")
		return
	}

	if arr[0]["condition"].(string) != "$test" {
		t.Error("error on unmarshalling rules")
		return
	}

	if err != nil {
		t.Error("unexpected error")
		return
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
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

	cfg := SetupConfig(s)

	ngl := services.NewGitlab(cfg)
	ngl.Connect()
	err := ngl.Save(dto, "test")

	if err != nil {
		t.Error("unexpected error")
	}
}
