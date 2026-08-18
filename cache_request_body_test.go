package gin

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheRequestBody(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := CreateTestContext(w)

	r.Use(CacheRequestBody())

	r.POST("/test", func(c *Context) {
		// First read: direct read from c.Request.Body
		body, err := io.ReadAll(c.Request.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"name":"gin"}`, string(body))

		// Second read: ShouldBindJSON
		var obj struct {
			Name string `json:"name"`
		}
		err = c.ShouldBindJSON(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "gin", obj.Name)

		// Third read: ShouldBindJSON again
		var obj2 struct {
			Name string `json:"name"`
		}
		err = c.ShouldBindJSON(&obj2)
		assert.NoError(t, err)
		assert.Equal(t, "gin", obj2.Name)

		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"name":"gin"}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCacheRequestBodyEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := CreateTestContext(w)

	r.Use(CacheRequestBody())

	r.POST("/test", func(c *Context) {
		var obj struct {
			Name string `json:"name"`
		}
		err := c.ShouldBindJSON(&obj)
		assert.Error(t, err) // Should return EOF or validation error, but not panic
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodPost, "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCacheDecodedBody(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := CreateTestContext(w)

	r.Use(CacheDecodedBody())

	r.POST("/test", func(c *Context) {
		var obj struct {
			Name string `json:"name"`
		}
		err := c.ShouldBindJSON(&obj)
		assert.NoError(t, err)
		assert.Equal(t, "gin", obj.Name)
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"name":"gin"}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
