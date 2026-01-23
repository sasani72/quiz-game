package main

import (
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
			Port:     3306,
			DBName:   "gameapp_db",
		},
	}

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)

	server.Serve()
}

//func userProfileHandler(writer http.ResponseWriter, request *http.Request) {
//	if request.Method != http.MethodGet {
//		fmt.Fprintf(writer, `"error": "invalid method"`)
//	}
//
//	authSvc := authservice.New(JwtSigningKey, AccessTokenSubject, RefreshTokenSubject,
//		AccessTokenExpireDuration, RefreshTokenExpireDuration)
//
//	authToken := request.Header.Get("Authorization")
//	claims, err := authSvc.ParseToken(authToken)
//	if err != nil {
//		fmt.Fprintf(writer, `"error": "invalid authorization token"`)
//	}
//
//	mysqlRepo := mysql.New()
//	userSvc := userservice.New(authSvc, mysqlRepo)
//
//	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
//	if err != nil {
//		writer.Write([]byte(
//			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	data, err := json.Marshal(resp)
//	if err != nil {
//		writer.Write([]byte(
//			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//	}
//	writer.Write(data)
//}
//
//func userLoginHandler(writer http.ResponseWriter, request *http.Request) {
//	if request.Method != http.MethodPost {
//		fmt.Fprintf(writer, `"error": "invalid method"`)
//	}
//
//	data, err := io.ReadAll(request.Body)
//	if err != nil {
//		writer.Write([]byte(
//			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
//		))
//
//		return
//	}
//
//	var lReq userservice.LoginRequest
//	err = json.Unmarshal(data, &lReq)
//	if err != nil {
//		writer.Write([]byte(
//			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	authSvc := authservice.New(JwtSigningKey, AccessTokenSubject, RefreshTokenSubject,
//		AccessTokenExpireDuration, RefreshTokenExpireDuration)
//
//	mysqlRepo := mysql.New()
//	userSvc := userservice.New(authSvc, mysqlRepo)
//	resp, err := userSvc.Login(lReq)
//	if err != nil {
//		writer.Write([]byte(
//			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	data, err = json.Marshal(resp)
//	if err != nil {
//		writer.Write([]byte(
//			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//		return
//	}
//
//	writer.Write(data)
//}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.MySQL)
	userSvc := userservice.New(authSvc, mysqlRepo)

	return authSvc, userSvc
}
