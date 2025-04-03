package router

import (
	"bankapi/docs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// ServeSwaggerYAML serves the Swagger YAML file
func ServeSwaggerYAML(c *gin.Context) {
	yamlFile, err := os.ReadFile("docs/swagger.yaml")
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not read swagger.yaml file")
		return
	}
	c.Data(http.StatusOK, "application/yaml", yamlFile)
}

func Endpoints(app *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/"
	app.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	// Serve Swagger YAML file
	app.GET("/get-swagger-file", ServeSwaggerYAML)

	// app.GET("/", func(c *gin.Context) {
	// 	responses.StatusOk(
	// 		c,
	// 		"Welcome to Bank api",
	// 		"",
	// 		"",
	// 	)
	// })

	CheckHealth(app)

	ProjectModules(app)
}

// @Tags Health API
// @Router /health [get]
func CheckHealth(app *gin.Engine) {
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"message": "Health check passed",
		})
	})
}
