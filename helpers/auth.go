package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/go-gorp/gorp"
	"github.com/haruyama/golang-goji-sample/models"
)

func Login(dbMap *gorp.DbMap, email string, password string) (*models.User, error) {
	var user models.User
	err := dbMap.SelectOne(&user, "SELECT * FROM Users WHERE Email = ?", email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, err
}
