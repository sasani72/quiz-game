package config

import (
	"quiz-game/repository/mysql"
	"quiz-game/service/authservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}
type Config struct {
	HTTPServer HTTPServer         `koanf:"http_server"`
	Auth       authservice.Config `koanf:"auth"`
	MySQL      mysql.Config       `koanf:"mysql"`
}
