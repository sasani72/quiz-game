package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	mux := http.DefaultServeMux
	mux.HandleFunc("/health-check", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func userProfileHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		fmt.Fprintf(writer, `"error": "invalid method"`)
	}

	authSvc := authservice.New(JwtSigningKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenExpireDuration, RefreshTokenExpireDuration)

	authToken := request.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(authToken)
	if err != nil {
		fmt.Fprintf(writer, `"error": "invalid authorization token"`)
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)

	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
	}
	writer.Write(data)
}

func userLoginHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		fmt.Fprintf(writer, `"error": "invalid method"`)
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	var lReq userservice.LoginRequest
	err = json.Unmarshal(data, &lReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}

	authSvc := authservice.New(JwtSigningKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenExpireDuration, RefreshTokenExpireDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)
	resp, err := userSvc.Login(lReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}

	data, err = json.Marshal(resp)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	writer.Write(data)
}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `"error": "invalid method"`)
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}

	authSvc := authservice.New(JwtSigningKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenExpireDuration, RefreshTokenExpireDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(authSvc, mysqlRepo)
	_, err = userSvc.Register(uReq)
	if err != nil {
		writer.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error())))

		return
	}

	writer.Write([]byte(`{"message": "user created successfully!"}`))
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "OK"`)
}
