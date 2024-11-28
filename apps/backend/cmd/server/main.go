package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/erentaskiran/project123123123/internal/api"
	db "github.com/erentaskiran/project123123123/pkg/database"
	"github.com/joho/godotenv"
)

type app struct {
	db *sql.DB
}

func main() {
	_ = godotenv.Load()

	port := ":8080"

	Db, err := db.SetupDb()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := app{
		db: Db,
	}

	router := api.NewRouter(app.db)

	r := router.NewRouter()

	fmt.Println("Server is running on port", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
