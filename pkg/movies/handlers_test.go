package movies

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emilyzhang/goyagi/pkg/application"
	"github.com/emilyzhang/goyagi/pkg/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("lists movies on success", func(tt *testing.T) {
		c, rr := newContext(tt, nil)

		err := h.listHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response []model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.True(tt, len(response) >= 23)
	})

	t.Run("optional parameters work as expected", func(tt *testing.T) {
		payload := []byte(`{"limit": 2}`)
		c, rr := newContext(tt, payload)

		err := h.listHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response []model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.True(tt, len(response) == 2)
	})
}

func TestRetrieveHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("retrieves movie on success", func(tt *testing.T) {
		c, rr := newContext(tt, nil)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.retrieveHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.Equal(tt, response.ID, 1)
		assert.Equal(tt, response.Title, "Iron Man")
	})

	t.Run("returns 404 if user isn't found", func(tt *testing.T) {
		c, _ := newContext(tt, nil)
		c.SetParamNames("id")
		c.SetParamValues("9999")

		err := h.retrieveHandler(c)
		assert.Contains(tt, err.Error(), "movie not found")
	})
}

func TestCreateHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("successfully creates new movie", func(tt *testing.T) {
		payload := []byte(`{"title": "Goyagi", "release_date": "2019-01-30T00:00:00.00Z"}`)
		c, rr := newContext(tt, payload)

		err := h.createHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.True(tt, response.Title == "Goyagi")

		// check that goyagi was added to list
		c, rr = newContext(tt, nil)
		c.SetParamNames("id")
		c.SetParamValues("24")

		err = h.retrieveHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.Equal(tt, response.ID, 24)
		assert.Equal(tt, response.Title, "Goyagi")
	})
}

// newHandler returns a new handler to be used for tests.
func newHandler(t *testing.T) handler {
	t.Helper()

	app, err := application.New()
	require.NoError(t, err)
	return handler{app}
}

// newContext returns a new echo.Context, and *httptest.ResponseRecorder to be
// used for tests.
func newContext(t *testing.T, payload []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)
	return c, rr
}
