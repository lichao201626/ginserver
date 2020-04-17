package resources

import (
	"ginserver/resources"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Users(t *testing.T) {

	r := gin.Default()
	resources.NewUserResource(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/get/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	t.Error("case fail", w.Body.String())
	// assert.Equal(t, "pong", w.Body.String())
}
