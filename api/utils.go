package main

import (
	"strconv"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
)

func LoadEmailConfig() (host, sender, password string, port int) {
	env, _ := utils.ReadEnvVariable("APP_ENV")

	if env == "DEVELOPMENT" {
		host, _ = utils.ReadEnvVariable("HOST")
		sender, _ = utils.ReadEnvVariable("TEST_EMAIL_SENDER")
		password, _ = utils.ReadEnvVariable("TEST_PASSWORD_SENDER")

		p, _ := utils.ReadEnvVariable("PORT")
		port, _ = strconv.Atoi(p)

	}
	return host, sender, password, port
}
