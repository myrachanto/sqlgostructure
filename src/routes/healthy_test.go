package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	endpoint := "/health"
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	e.GET(endpoint, HealthCheck)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var res map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.EqualValues(t, "server is up and running", res["data"])
}
