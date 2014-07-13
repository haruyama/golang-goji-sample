package helpers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/coopernurse/gorp"
	"github.com/haruyama/golang-goji-sample/models"
)

func Login(dbMap *gorp.DbMap, email string, password string) (*models.User, error) {
	var user models.User
	err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE Email = ?", email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, err
}
