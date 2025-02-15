package data

import (
	"errors"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserType interface {
	ValidateUserRequest(v *validator.ValidateData)
}

type Users struct {
	UserName    string   `json:"user_name" redis:"user_name"`
	Email       string   `json:"email" redis:"email"`
	Password    password `json:"-"`
	IsAdmin     bool     `json:"-"`
	IsActivated bool     `json:"-"`
}

type password struct {
	PasswordText   *string
	HashedPassword []byte
}

func validateEmail(input string, v *validator.ValidateData) {
	v.CheckIsError(input == "", "email", "email field cannot be empty")
	v.CheckIsError(!validator.MatchPattern(validator.EMAIL_REGEX, input), "email", "invalid email provided")
}

func validatePassword(password string, v *validator.ValidateData) {
	v.CheckIsError((len(password) >= 8) == false, "password", "password too short must be greater than or equal to 8")
	v.CheckIsError((len(password) <= 72) == false, "password", "password field too long must not be more than 72 bytes long")
}
func (u *Users) ValidateUserRequest(v *validator.ValidateData) {
	v.CheckIsError(u.UserName == "", "username", "username field cannot be empty")
	v.CheckIsError((len(u.UserName) <= 500) == false, "username", "username field cannot be 500 bytes long")

	// validate email
	validateEmail(u.Email, v)

	// set user password
	if u.Password.PasswordText != nil {
		validatePassword(*u.Password.PasswordText, v)
	}
	if u.Password.HashedPassword == nil {
		panic("Missing password hashing")
	}
}

// handle user password
func (p *password) SetPassword(plainPassword string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 15)

	if err != nil {
		return err
	}
	p.PasswordText = &plainPassword
	p.HashedPassword = hashPassword
	return nil
}

// password match...
func (p *password) PasswordMatch(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
