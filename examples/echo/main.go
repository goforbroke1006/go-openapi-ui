package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	go_openapi_ui "github.com/goforbroke1006/go-openapi-ui"
)

//go:embed openapi.yaml
var openapiSpec embed.FS

// http://localhost:8080/doc/_swagger/openapi.yaml
// http://localhost:8080/doc/_swagger/index.html
func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	go_openapi_ui.EchoSwaggerUIAndSingleSpec(e, "/doc/_swagger", openapiSpec)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
