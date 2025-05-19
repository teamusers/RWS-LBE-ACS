package router

import (
	"fmt"
	"log"
	"net/http"

	general "rlp-email-service/api/http"
	"rlp-email-service/api/http/middleware"
	"rlp-email-service/config"
	"rlp-email-service/model"
	"rlp-email-service/system"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option type and global slice for router modifications.
type Option func(*gin.RouterGroup)

var options = []Option{}

var endpointList []map[string]string

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	// Include additional routers
	Include(general.Routers)

	db := system.GetDb()
	if err := model.MigrateAuditLog(db); err != nil {
		log.Fatalf("audit log migration: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.AuditLogger(db))

	apiGroup := r.Group("/api")
	for _, opt := range options {
		opt(apiGroup)
	}

	// Capture routes but exclude the HandlerFunc to avoid JSON marshalling errors.
	routes := r.Routes()
	for _, route := range routes {
		endpointList = append(endpointList, map[string]string{
			"method": route.Method,
			"path":   route.Path,
		})
	}
	r.Static("/docs", "./api/docs")

	// wire up the swagger UI, telling it to fetch /docs/swagger.json
	url := ginSwagger.URL("/docs/swagger.json")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		url,
		// ⟵ hides the Models section
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// redirect root to swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
	// also catch bare /swagger
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})

	r.Run(fmt.Sprintf(":%d", config.GetConfig().Http.Port))
	return r
}
