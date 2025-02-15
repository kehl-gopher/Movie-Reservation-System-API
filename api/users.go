package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/data"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/models"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/validator"
)

// handle users api routing

// handle user registration... routes
func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {

	var users struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}
	err := readFromJson(r, &users)

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	user := &data.Users{
		UserName: users.UserName,
		Email:    users.Email,
		IsAdmin:  users.IsAdmin,
	}

	err = user.Password.SetPassword(users.Password)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	v := validator.NewValidator()
	if user.ValidateUserRequest(v); !v.CheckErrorExists() {
		fmt.Println(v.CheckErrorExists())
		app.validationErrorResponse(w, toJson{"errors": v.Errors})
		return
	}

	err = app.model.Users.CreateUser(user)

	if err != nil {

		switch {
		case errors.Is(err, models.ErrDuplicationValue):
			v.AddError("username", "user with this username already exists")
			app.validationErrorResponse(w, toJson{"errors": v.Errors})
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.mailer.Send(user.Email, "user_welcome.html", user)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	// var user data.UserType = &data.Users{UserName: users.UserName, Email: users.Email}
	app.writeResponse(w, http.StatusCreated, toJson{"message": "User created successfully"})
}
