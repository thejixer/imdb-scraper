package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	godotenv.Load()
}

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "IMDB SCRAPER")
	})
	e.GET("/m/:slug", GetMovieData)
	listenAddr := fmt.Sprintf(":%v", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(listenAddr))

}
