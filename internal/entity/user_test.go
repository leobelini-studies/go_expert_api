package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "1y3t3@example.com", "password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "1y3t3@example.com", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("John Doe", "1y3t3@example.com", "password")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("password"))
	assert.False(t, user.ValidatePassword("wrong-password"))
	assert.NotEqual(t, user.Password, "wrong-password")
}
