package services

import (
	_ "github.com/lib/pq"
)

// func TestValidLogin(t *testing.T) {
// 	db, err := utilities.GetMainDatabaseConnection(utilities.DatabaseConfig{
// 		Host:     "localhost",
// 		Userame:  "postgres",
// 		Password: "root",
// 		Database: "apiauth",
// 		Port:     "5432",
// 	})
// 	if err != nil {
// 		t.Error("Connection to database failed: ", err)
// 		return
// 	}
// 	authService := NewAuthService(db)
// 	_, err = authService.Login("jackie", "password", "", "")

// 	if err != nil {
// 		t.Error("Failed to authenticate")
// 	}
// }
