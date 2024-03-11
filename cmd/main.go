package main

import (
	"os"

	"github.com/edwinbm5/go-product-web/internal/application"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	App := application.NewDefaultApp(application.ConfigDefaultApp{
		Title:    os.Getenv("APP_TITLE"),
		Color:    os.Getenv("APP_CLI_COLOR"),
		FilePath: os.Getenv("DB_PATH") + os.Getenv("DB_FILE_NAME"),
		Token:    os.Getenv("TOKEN"),
	})

	App.Run()
}
