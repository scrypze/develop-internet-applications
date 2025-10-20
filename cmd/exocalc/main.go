package main

import (
	"develop-internet-applications/internal/app"
)

// @title ExoCalc API
// @version 1.0
// @description API для работы с экзопланетами и звездами

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @license.name AS IS (NO WARRANTY)

// @host localhost:8080
// @schemes http https
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer {token}

func main() {
	app := app.NewApp()
	app.RunApp()
}
