package go_openapi_ui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// EchoSwaggerUIAndSingleSpec appends handles to echo.Echo to provide swagger-ui functionality.
// Can be handled only one openapi specification file.
//
//	/url/prefix/<openapi.yaml|openapi.json> - handle spec file;
//	/url/prefix/swagger-initializer.js      - handle generate initializer for UI;
//	/url/prefix/index.html                  - handle swagger-ui static files;
func EchoSwaggerUIAndSingleSpec(router *echo.Echo, urlPrefix string, specFile embed.FS) {
	specFSDir, _ := specFile.ReadDir(".")
	if len(specFSDir) != 1 {
		panic("required single openapi specification file")
	}
	specBaseFileName := specFSDir[0].Name()

	if !strings.HasPrefix(urlPrefix, "/") {
		urlPrefix = "/" + urlPrefix
	}

	if strings.HasSuffix(urlPrefix, "/") {
		urlPrefix = strings.TrimSuffix(urlPrefix, "/")
	}

	var (
		specPath = fmt.Sprintf("%s/%s", urlPrefix, specBaseFileName)

		initJSPath = fmt.Sprintf("%s/swagger-initializer.js", urlPrefix)

		swaggerUIPath = fmt.Sprintf("%s/*", urlPrefix)
		stripPrefix   = fmt.Sprintf("%s/", urlPrefix)
	)

	router.FileFS(specPath, specBaseFileName, specFile)

	router.GET(initJSPath, func(c echo.Context) error {
		const initJsContentTpl = `
window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "%s",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
`

		initJsContent := fmt.Sprintf(initJsContentTpl, specPath)

		c.Response().Header().Set("Content-Type", "application/javascript")
		c.Response().WriteHeader(http.StatusOK)
		_, writeErr := c.Response().Write([]byte(initJsContent))
		return writeErr
	})

	swaggerUIFileSystem, _ := fs.Sub(swaggerUIStatic, "web/swagger-ui/dist")
	fileServer := http.FileServer(http.FS(swaggerUIFileSystem))

	router.GET(swaggerUIPath, echo.WrapHandler(http.StripPrefix(stripPrefix, fileServer)))
}
