package config

import (
	"quiz-game/repository/mysql"
	"quiz-game/service/authservice"
)

type HTTPServer struct {
	Port int
}
type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	MySQL      mysql.Config
}
