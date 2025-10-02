package router

import (
	// initialize the Swagger documentation
	_ "app/src/docs"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func DocsRoutes(v1 fiber.Router) {
	docs := v1.Group("/docs")

	// Konfigurasi Swagger dengan persistence authorization
	config := swagger.Config{
		URL:                    "/v1/docs/doc.json",
		DeepLinking:            true,
		DocExpansion:           "none",
		Layout:                 "StandaloneLayout",
		PersistAuthorization:   true, // Ini yang paling penting! Mengaktifkan persistence authorization
		TryItOutEnabled:        true,
		DisplayRequestDuration: true,
		ShowExtensions:         false,
		ShowCommonExtensions:   false,
		// Custom CSS untuk improve tampilan
		CustomStyle: template.CSS(`
			/* Custom styles untuk improve UX */
			.swagger-ui .topbar { 
				background-color: #1f2937; 
			}
			.swagger-ui .info .title {
				color: #3b82f6;
				font-weight: bold;
			}
			.swagger-ui .info .description {
				color: #6b7280;
			}
			.swagger-ui .auth-wrapper .authorize:hover {
				background-color: #059669;
				border-color: #059669;
			}
			/* Success message styling */
			.swagger-ui .success {
				color: #10b981;
			}
		`),
		// Custom JS untuk memberikan feedback ke user bahwa token akan persist
		CustomScript: template.JS(`
			// Script untuk memberikan informasi ke user
			(function() {
				console.log('Swagger UI: Authorization persistence is ENABLED');
				
				// Tambahkan informasi di UI
				setTimeout(function() {
					const authWrapper = document.querySelector('.auth-wrapper');
					if (authWrapper) {
						const infoDiv = document.createElement('div');
						authWrapper.appendChild(infoDiv);
					}
				}, 1000);
			})();
		`),
	}

	docs.Get("/*", swagger.New(config))
}
