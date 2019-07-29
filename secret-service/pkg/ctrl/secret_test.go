package ctrl_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/nurali/secret-server/secret-service/pkg/app"
	"github.com/nurali/secret-server/secret-service/pkg/config"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SecretCtrlSuite struct {
	suite.Suite
	Router *mux.Router
}

func TestSecretCtrl(t *testing.T) {
	cfg := config.New()

	db, err := app.OpenDB(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := app.SetupDB(db); err != nil {
		t.Fatal(err)
	}

	s := &SecretCtrlSuite{
		Router: app.Router(db),
	}

	suite.Run(t, s)
}

func (s *SecretCtrlSuite) TestSecretCreate() {
	s.T().Run("ok", func(t *testing.T) {
		reqBody := `{"secret":"mysec1", "expireAfterViews":10, "expireAfter":5}`
		r := httptest.NewRequest("POST", "/api/secret", bytes.NewBuffer([]byte(reqBody)))
		w := httptest.NewRecorder()

		// test
		s.Router.ServeHTTP(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(respBody), "mysec1")
	})

	s.T().Run("invalid_expire_after_views", func(t *testing.T) {
		reqBody := `{"secret":"mysec1", "expireAfterViews":0, "expireAfter":5}`
		r := httptest.NewRequest("POST", "/api/secret", bytes.NewBuffer([]byte(reqBody)))
		w := httptest.NewRecorder()

		// test
		s.Router.ServeHTTP(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		respBody, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(respBody), "Invalid expireAfterViews")
	})
}

func (s *SecretCtrlSuite) TestSecretShow() {
	reqBody := `{"secret":"mysec1", "expireAfterViews":10, "expireAfter":5}`
	r := httptest.NewRequest("POST", "/api/secret", bytes.NewBuffer([]byte(reqBody)))
	w := httptest.NewRecorder()

	s.Router.ServeHTTP(w, r)

	resp := w.Result()
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	hash := parseResp(resp.Body)["hash"]
	require.NotEmpty(s.T(), hash)

	s.T().Run("ok", func(t *testing.T) {
		url := fmt.Sprintf("/api/secret/%s", hash)
		t.Logf("url=%s", url)
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		// test
		s.Router.ServeHTTP(w, r)

		resp := w.Result()
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		secretAttrs := parseResp(resp.Body)
		assert.Equal(t, hash, secretAttrs["hash"])
	})

	s.T().Run("not_found", func(t *testing.T) {
		url := fmt.Sprintf("/api/secret/%s", uuid.NewV4().String())
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		// test
		s.Router.ServeHTTP(w, r)

		resp := w.Result()
		assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
	})

}

func parseResp(content io.Reader) map[string]string {
	secretAttrs := make(map[string]string)
	decode := json.NewDecoder(content)
	decode.Decode(&secretAttrs)
	return secretAttrs
}
