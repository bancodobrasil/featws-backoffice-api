package services

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/bancodobrasil/featws-api/config"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/require"
// 	gl "github.com/xanzy/go-gitlab"
// )

// // func setup(t *testing.T) (*gin.Engine, *httptest.Server, *gl.Client) {
// // 	gin := gin.Default()

// // 	server := httptest.NewServer(gin)

// // 	client, err := gl.NewClient("", gl.WithBaseURL(server.URL))
// // 	if err != nil {
// // 		server.Close()
// // 		t.Fatalf("Failed to create client: %v", err)
// // 	}
// // 	return gin, server, client
// // }

// func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *gl.Client) {
// 	// mux is the HTTP request multiplexer used with the test server.
// 	mux := http.NewServeMux()

// 	// server is a test HTTP server used to provide mock API responses.
// 	server := httptest.NewServer(mux)

// 	// client is the Gitlab client being tested.
// 	client, err := gl.NewClient("", gl.WithBaseURL(server.URL))
// 	if err != nil {
// 		server.Close()
// 		t.Fatalf("Failed to create client: %v", err)
// 	}

// 	return mux, server, client
// }

// func teardown(server *httptest.Server) {
// 	server.Close()
// }

// func testMethod(t *testing.T, r *http.Request, want string) {
// 	if got := r.Method; got != want {
// 		t.Errorf("Request method: %s, want %s", got, want)
// 	}
// }

// // func GetCommit(c *gin.Context) {
// // 	testMethod(t, r, http.MethodGet)
// // 	fmt.Fprintf(c.Writer, `
// // 		{
// // 			"file_name": "key.rb",
// // 			"file_path": "app/models/key.rb",
// // 			"size": 1476,
// // 			"encoding": "base64",
// // 			"content": "IyA9PSBTY2hlbWEgSW5mb3...",
// // 			"content_sha256": "4c294617b60715c1d218e61164a3abd4808a4284cbc30e6728a01ad9aada4481",
// // 			"ref": "master",
// // 			"blob_id": "79f7bbd25901e8334750839545a9bd021f0e4c83",
// // 			"commit_id": "d5a3ff139356ce33e37e73add446f16869741b50",
// // 			"last_commit_id": "570e7b2abdd848b95f2f578043fc23bd6f6fd24d"
// // 		}
// // 	`)
// // }

// func TestSaveInGitlab(t *testing.T) {
// 	mux, server, client := setup(t)
// 	defer teardown(server)

// 	mux.HandleFunc("/api/v4/projects/13083/repository/files/app%2Fmodels%2Fkey%2Erb?ref=master", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, http.MethodGet)
// 		fmt.Fprintf(w, `
// 			{
// 			"file_name": "key.rb",
// 			"file_path": "app/models/key.rb",
// 			"size": 1476,
// 			"encoding": "base64",
// 			"content": "IyA9PSBTY2hlbWEgSW5mb3...",
// 			"content_sha256": "4c294617b60715c1d218e61164a3abd4808a4284cbc30e6728a01ad9aada4481",
// 			"ref": "master",
// 			"blob_id": "79f7bbd25901e8334750839545a9bd021f0e4c83",
// 			"commit_id": "d5a3ff139356ce33e37e73add446f16869741b50",
// 			"last_commit_id": "570e7b2abdd848b95f2f578043fc23bd6f6fd24d"
// 			}
// 		`)
// 	})

// 	want := &gl.File{
// 		FileName:     "key.rb",
// 		FilePath:     "app/models/key.rb",
// 		Size:         1476,
// 		Encoding:     "base64",
// 		Content:      "IyA9PSBTY2hlbWEgSW5mb3...",
// 		Ref:          "master",
// 		BlobID:       "79f7bbd25901e8334750839545a9bd021f0e4c83",
// 		CommitID:     "d5a3ff139356ce33e37e73add446f16869741b50",
// 		SHA256:       "4c294617b60715c1d218e61164a3abd4808a4284cbc30e6728a01ad9aada4481",
// 		LastCommitID: "570e7b2abdd848b95f2f578043fc23bd6f6fd24d",
// 	}

// 	f, resp, err := client.RepositoryFiles.GetFile(13083, "app%2Fmodels%2Fkey%2Erb?ref=master", nil)
// 	require.NoError(t, err)
// 	require.NotNil(t, resp)
// 	require.Equal(t, want, f)

// 	f, resp, err = client.RepositoryFiles.GetFile(13083.01, "app%2Fmodels%2Fkey%2Erb?ref=master", nil)
// 	require.EqualError(t, err, "invalid ID type 13083.01, the ID must be an int or a string")
// 	require.Nil(t, resp)
// 	require.Nil(t, f)

// 	f, resp, err = client.RepositoryFiles.GetFile(13084, "app%2Fmodels%2Fkey%2Erb?ref=master", nil)
// 	require.Error(t, err)
// 	require.Nil(t, f)
// 	require.Equal(t, http.StatusNotFound, resp.StatusCode)
// }

// func TestCreateOrUpdateGitlabFileCommitAction(t *testing.T) {

// }

// func TestDefineCreateOrUpdateGitlabFileAction(t *testing.T) {

// }

// func TestFillWithGitlab(t *testing.T) {

// }

// func TestConnectGitlab(t *testing.T) {
// 	err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("Não foi possível carregar as configurações: %s\n", err)
// 	}

// 	cfg := config.Config{}

// 	c, err := ConnectGitlab(&cfg)
// 	if err != nil {
// 		t.Fatalf("Failed to create client: %v", err)

// 	}

// 	expectedBaseURL := "/api/v4/"

// 	if c.BaseURL().String() != expectedBaseURL {
// 		t.Errorf("NewClient BaseUrl is %s, want %s", c.BaseURL().String(), expectedBaseURL)
// 	}

// }

// type Config struct {
// 	Test string `mapstructure:"TEST"`
// }

// var config1 = &Config{}

// func GetMockConfig() *Config {
// 	return config1
// }
// func TestConnectGitlabFail(t *testing.T) {

// 	mockCfg := config1

// 	c, err := ConnectGitlab(mockCfg)
// 	got := err
// 	expected := ""

// 	if got.Error() != expected {

// 	}

// 	expectedBaseURL := "/api/v4/"

// 	if c.BaseURL().String() != expectedBaseURL {
// 		t.Errorf("NewClient BaseUrl is %s, want %s", c.BaseURL().String(), expectedBaseURL)
// 	}

// }

// func TestGitlabLoadJSON(t *testing.T) {

// }

// func TestGitlabLoadString(t *testing.T) {

// }
