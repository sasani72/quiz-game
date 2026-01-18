package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"quiz-game/entity"
	"quiz-game/repository/mysql"
	"quiz-game/service/userservice"
)

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/health-check", healthCheckHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
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

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)
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

func testUserMysqlRepo() {
	mysqlRepo := mysql.New()

	mysqlRepo.Register(entity.User{
		PhoneNumber: "0912",
		Name:        "Test User",
	})
}
