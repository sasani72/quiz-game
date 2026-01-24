package main

import (
	"fmt"
	"quiz-game/config"
	"quiz-game/delivery/httpserver"
	"quiz-game/repository/mysql"
	"quiz-game/service/authservice"
	"quiz-game/service/userservice"
	"time"
)

const (
	JwtSigningKey              = "jwt_secret"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth: authservice.Config{
			SignKey:               JwtSigningKey,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
		},
		MySQL: mysql.Config{
			Username: "gameapp",
			Password: "root",
			Host:     "localhost",
			Port:     3308,
			DBName:   "gameapp_db",
		},
	}

	// add command for migrations
	//mgr := migrator.New(cfg.MySQL)
	//mgr.Up()
	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)
	fmt.Println("start echo server")
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.MySQL)
	userSvc := userservice.New(authSvc, mysqlRepo)

	return authSvc, userSvc
}
