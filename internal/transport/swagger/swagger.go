package swagger

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v5"
	swaggerFiles "github.com/swaggo/files/v2"
	"github.com/swaggo/swag"
)

type Config struct {
	URL          string
	DeepLinking  bool
	DocExpansion string
	DomID        string
}

func NewConfig() *Config {
	return &Config{
		URL:          "doc.json",
		DeepLinking:  true,
		DocExpansion: "list",
		DomID:        "swagger-ui",
	}
}

func Handler(config *Config) echo.HandlerFunc {
	if config == nil {
		config = NewConfig()
	}

	t := template.Must(template.New("swagger").Parse(indexTemplate))

	return func(c *echo.Context) error {
		path := c.Param("*")

		switch path {
		case "", "/":
			return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		case "index.html":
			c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
			return t.Execute(c.Response(), config)
		case "doc.json":
			doc, err := swag.ReadDoc()
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.JSONBlob(http.StatusOK, []byte(doc))
		default:
			file, err := swaggerFiles.FS.Open(path)
			if err != nil {
				return c.String(http.StatusNotFound, "file not found")
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			contentType := getContentType(path)
			c.Response().Header().Set("Content-Type", contentType)
			return c.Stream(http.StatusOK, contentType, &fileReaderWrapper{file: file, size: stat.Size()})
		}
	}
}

type fileReaderWrapper struct {
	file fs.File
	size int64
}

func (f *fileReaderWrapper) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func getContentType(path string) string {
	switch {
	case len(path) > 3 && path[len(path)-3:] == ".js":
		return "application/javascript"
	case len(path) > 4 && path[len(path)-4:] == ".css":
		return "text/css"
	case len(path) > 4 && path[len(path)-4:] == ".png":
		return "image/png"
	case len(path) > 5 && path[len(path)-5:] == ".html":
		return "text/html"
	default:
		return "application/octet-stream"
	}
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Swagger UI</title>
    <link rel="stylesheet" href="swagger-ui.css">
    <link rel="icon" type="image/png" href="favicon-32x32.png" sizes="32x32">
</head>
<body>
    <div id="{{.DomID}}"></div>
    <script src="swagger-ui-bundle.js"></script>
    <script src="swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "{{.URL}}",
                dom_id: '#{{.DomID}}',
                deepLinking: {{.DeepLinking}},
                docExpansion: "{{.DocExpansion}}",
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
