package middleware
import (
	token "ecommerce/tokens"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)
func generateTestToken(email, uid string) (string, error) {
	token, _, err := token.TokenGenerator(email, "", "", uid)
	return token, err
}
func TestAuthenticationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/protected", Authentication(), func(c *gin.Context) {
		email, _ := c.Get("email")
		uid, _ := c.Get("uid")
		c.JSON(http.StatusOK, gin.H{"email": email, "uid": uid})
	})
	t.Run("Valid Token", func(t *testing.T) {
		token, err := generateTestToken("test@example.com", "123456")
		assert.NoError(t, err)
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("token", token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"email":"test@example.com"`)
		assert.Contains(t, w.Body.String(), `"uid":"123456"`)
	})
	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("token", "invalid-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expected := map[string]string{
			"error": "token contains an invalid number of segments",
		}
		expectedJSON, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("Error marshaling expected JSON: %v", err)
		}
		assert.Equal(t, string(expectedJSON), w.Body.String())
	})
	t.Run("Missing Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"No Authorization Header Provided"`)
	})
}