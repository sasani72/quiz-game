package main

import (
	"quiz-game/entity"
	"quiz-game/repository/mysql"
)

func main() {
	mysqlRepo := mysql.New()

	mysqlRepo.Register(entity.User{
		PhoneNumber: "0912",
		Name:        "Test User",
	})
}
