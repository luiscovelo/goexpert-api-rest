package database

import (
	"testing"

	"github.com/luiscovelo/goexpert-api-rest/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUserAndFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.User{})

	user, err := entity.NewUser("Luis", "luis@luis.com.br", "1234")
	assert.Nil(t, err)

	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	userCreated, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.NotNil(t, userCreated)

	assert.Equal(t, user.ID, userCreated.ID)
	assert.Equal(t, "Luis", userCreated.Name)
	assert.Equal(t, "luis@luis.com.br", userCreated.Email)
	assert.True(t, userCreated.ValidatePassword("1234"))
}
