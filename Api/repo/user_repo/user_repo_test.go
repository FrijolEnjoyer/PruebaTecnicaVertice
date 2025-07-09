package user_repo

import (
	"testing"

	"pruebaVertice/Api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sirupsen/logrus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)
	return db
}

func TestCreateUser_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	input := &models.User{Username: "u1", Email: "e1@e.com", Password: "pwd"}
	created, err := repo.CreateUser(input)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "u1", created.Username)
	assert.Equal(t, "e1@e.com", created.Email)
}

func TestGetUserByID_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	u := &models.User{Username: "u2", Email: "e2@e.com", Password: "pwd"}
	repo.CreateUser(u)

	fetched, err := repo.GetUserByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "u2", fetched.Username)
	assert.Equal(t, "e2@e.com", fetched.Email)
}

func TestUpdateUser_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	u := &models.User{Username: "u3", Email: "e3@e.com", Password: "pwd"}
	repo.CreateUser(u)

	u.Username = "u3-upd"
	updated, err := repo.UpdateUser(u)
	assert.NoError(t, err)
	assert.Equal(t, "u3-upd", updated.Username)

	// verify persistence
	fetched, _ := repo.GetUserByID("1")
	assert.Equal(t, "u3-upd", fetched.Username)
}

func TestDeleteUser_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	u := &models.User{Username: "u4", Email: "e4@e.com", Password: "pwd"}
	repo.CreateUser(u)
	err := repo.DeleteUser("1")
	assert.NoError(t, err)

	_, err = repo.GetUserByID("1")
	assert.Error(t, err)
}

func TestUpdateUserToken_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	u := &models.User{Username: "u5", Email: "e5@e.com", Password: "pwd"}
	repo.CreateUser(u)

	u.Token = "tok"
	u.RefreshToken = "ref"
	res, err := repo.UpdateUserToken(u)
	assert.NoError(t, err)
	assert.Equal(t, "tok", res.Token)
	assert.Equal(t, "ref", res.RefreshToken)
}

func TestGetUserByEmail_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	u := &models.User{Username: "u6", Email: "e6@e.com", Password: "pwd"}
	repo.CreateUser(u)

	res, err := repo.GetUserByEmail("e6@e.com")
	assert.NoError(t, err)
	assert.Equal(t, "u6", res.Username)
}

func TestGetUserByEmail_Error(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewUserRepository(db, logger)

	_, err := repo.GetUserByEmail("not@found")
	assert.Error(t, err)
}
