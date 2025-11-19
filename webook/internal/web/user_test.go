package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)
import "golang.org/x/crypto/bcrypt"

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)

}
