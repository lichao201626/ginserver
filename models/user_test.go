package models

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	u := NewUser()
	assert.Equal(t, "models.User", reflect.TypeOf(u).String())
	// assert.Equal(t, "", u.Email, "email should be blank")
}
