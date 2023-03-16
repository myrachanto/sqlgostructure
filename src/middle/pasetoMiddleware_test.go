package middle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/sqlgostructure/src/pasetos"
	"github.com/stretchr/testify/require"
)

var data = &pasetos.Data{
	Username: "myrachanto",
	Email:    "myrachanto@gmail.com",
	Bizname:  "shop1",
	// Shops: [],
}

func addAuthorization(t *testing.T, request *http.Request, Authorizationtype string, data *pasetos.Data, duration time.Duration) {
	PasetoMaker, _ := pasetos.NewPasetoMaker()
	token, payload, err := PasetoMaker.CreateToken(data, duration)
	require.EqualValues(t, nil, err)
	require.NotEmpty(t, payload)
	authorizationHeader := fmt.Sprintf("%s %s", Authorizationtype, token)
	request.Header.Set(authorisationHeaderKey, authorizationHeader)
}
func TestPasetoAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorisationType, data, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
				// addAuthorization(t, request, authorisationType, data, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "unsupported", data, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "", data, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "", data, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// server := newTe
			e := echo.New()
			authpath := "/auth"
			e.GET(authpath, func(c echo.Context) error {
				return c.JSON(http.StatusOK, "Querried successifuly")
			}, PasetoAuthMiddleware)

			// e.Logger.Fatal(e.Start("2200"))
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authpath, nil)
			require.EqualValues(t, nil, err)
			tc.setupAuth(t, request)
			e.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
