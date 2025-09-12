package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// SwaggerIndex serves a custom Swagger UI page that enables persistAuthorization
func SwaggerIndex(c *fiber.Ctx) error {
	// Use the embedded swagger doc endpoint provided by swaggo/fiber-swagger: /swagger/doc.json
	html := fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>API Docs</title>
    <link href="https://unpkg.com/swagger-ui-dist@4.18.3/swagger-ui.css" rel="stylesheet" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.18.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.18.3/swagger-ui-standalone-preset.js"></script>
    <script>
      window.addEventListener('load', function() {
        const ui = SwaggerUIBundle({
          url: window.location.origin + '/swagger/doc.json',
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
          layout: 'StandaloneLayout',
          persistAuthorization: true
        })
        window.ui = ui
      })
    </script>
  </body>
</html>`)

	c.Type("html")
	return c.SendString(html)
}
